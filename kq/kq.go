package kq

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/colinrs/pkgx/concurrent"
	goSafe "github.com/colinrs/pkgx/fx"
	"github.com/colinrs/pkgx/kq/core"

	"go.uber.org/atomic"
)

var (
	kqStatusInit    int32 = 1
	kqStatusRunning int32 = 2
	kqStatusClosed  int32 = 3
)

type KQ struct {
	kqStatus *atomic.Int32

	input       core.Input
	extractor   core.Extractor
	middleware  core.Middleware
	transformer core.Transformer
	output      core.Output

	inputMessageChannel  chan *core.InputMessage
	commitMessageChanel  chan *core.InputMessage
	extractorMessageChan chan *internalMessage
	outPutMessageChanel  chan *internalMessage
	limitGoroutines      *concurrent.Limit
	maxGoroutines        int
}

var (
	internalMessagePool = sync.Pool{
		New: func() interface{} {
			return &internalMessage{}
		},
	}
)

func NewKQ(ctx context.Context, opts ...Option) *KQ {
	o := &options{
		inputMessageChannelSize:     defaultInputMessageChannelSize,
		maxGoroutines:               defaultMaxGoroutines,
		commitMessageChannelSize:    defaultCommitMessageChannelSize,
		extractorMessageChannelSize: defaultExtractorMessageChannelSize,
		outPutMessageChannelSize:    defaultOutPutMessageChannelSize,
	}
	for _, opt := range opts {
		opt(o)
	}

	return &KQ{
		kqStatus:             atomic.NewInt32(kqStatusInit),
		inputMessageChannel:  make(chan *core.InputMessage, o.inputMessageChannelSize),
		commitMessageChanel:  make(chan *core.InputMessage, o.commitMessageChannelSize),
		extractorMessageChan: make(chan *internalMessage, o.extractorMessageChannelSize),
		outPutMessageChanel:  make(chan *internalMessage, o.extractorMessageChannelSize),
		limitGoroutines:      concurrent.NewLimit(o.maxGoroutines),
	}
}

func (k *KQ) SetInput(ctx context.Context, input core.Input) *KQ {
	k.input = input
	return k
}

func (k *KQ) SetExtractor(ctx context.Context, extractor core.Extractor) *KQ {
	k.extractor = extractor
	return k
}

func (k *KQ) SetMiddleware(ctx context.Context, middleware core.Middleware) *KQ {
	k.middleware = middleware
	return k
}

func (k *KQ) SetTransformer(ctx context.Context, transformer core.Transformer) *KQ {
	k.transformer = transformer
	return k
}

func (k *KQ) SetOutput(ctx context.Context, output core.Output) *KQ {
	k.output = output
	return k
}

func (k *KQ) Run(ctx context.Context) error {
	k.kqStatus.Store(kqStatusRunning)
	goSafe.GoSafeWithRecover(func() {
		err := k.inputMessage(ctx)
		if err != nil {
			fmt.Println(err.Error())
		}
	}, kqRecover())
	goSafe.GoSafeWithRecover(func() {
		err := k.commitMessage(ctx)
		if err != nil {
			fmt.Println(err.Error())
		}
	}, kqRecover())
	goSafe.GoSafeWithRecover(func() {
		err := k.inputMessageExtractor(ctx)
		if err != nil {
			fmt.Println(err.Error())
		}
	}, kqRecover())
	goSafe.GoSafeWithRecover(func() {
		err := k.extractorMessageProcess(ctx)
		if err != nil {
			fmt.Println(err.Error())
		}
	}, kqRecover())
	return nil
}

func (k *KQ) inputMessage(ctx context.Context) error {
	ticket := time.Tick(time.Millisecond * 500)
	for k.kqStatus.Load() == kqStatusRunning {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticket:
			runtime.Gosched()
		default:
			m, err := k.input.Consumer(ctx)
			if err != nil || m == nil {
				continue
			}
			m.Lock()
			k.inputMessageChannel <- m
			k.commitMessageChanel <- m
		}
	}
	return nil
}

func (k *KQ) inputMessageExtractor(ctx context.Context) error {
	for k.kqStatus.Load() == kqStatusRunning {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case inputMessage, ok := <-k.inputMessageChannel:
			if !ok {
				return nil
			}
			k.limitGoroutines.Acquire()
			goSafe.GoSafeWithRecover(func() {
				defer k.limitGoroutines.Release()
				extractorMessage, err := k.extractor.Unmarshal(ctx, inputMessage)
				iMessage := getInternalMessage()
				iMessage.inputMessage = inputMessage
				iMessage.extractorMessage = extractorMessage
				if err != nil {
					inputMessage.Ack()
					k.extractor.OnError(ctx, inputMessage, err)
				} else {
					k.extractor.OnDone(ctx, inputMessage)
					k.extractorMessageChan <- iMessage
				}
			}, kqRecover())
		}
	}
	return nil
}

func (k *KQ) extractorMessageProcess(ctx context.Context) error {
	for k.kqStatus.Load() == kqStatusRunning {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case iMessage, ok := <-k.extractorMessageChan:
			if !ok {
				return nil
			}
			k.limitGoroutines.Acquire()
			goSafe.GoSafeWithRecover(func() {
				defer k.limitGoroutines.Release()
				outPutMessage, err := k.transformer.Process(ctx, iMessage.extractorMessage)
				if err != nil {
					iMessage.inputMessage.Ack()
					k.transformer.OnError(ctx, iMessage.extractorMessage, err)
				} else {
					k.transformer.OnDone(ctx, iMessage.extractorMessage)
					iMessage.outPutMessage = outPutMessage
					iMessage.extractorMessage = nil
					k.outPutMessageChanel <- iMessage
				}
			}, kqRecover())
		}
	}
	return nil
}

func (k *KQ) transformerMessageOutPut(ctx context.Context) error {
	for k.kqStatus.Load() == kqStatusRunning {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case iMessage, ok := <-k.outPutMessageChanel:
			if !ok {
				return nil
			}
			k.limitGoroutines.Acquire()
			goSafe.GoSafeWithRecover(func() {
				defer k.limitGoroutines.Release()
				err := k.output.SendOutput(ctx, iMessage.outPutMessage)
				if err != nil {
					k.output.OnError(ctx, iMessage.outPutMessage, err)
				} else {
					k.output.OnDone(ctx, iMessage.outPutMessage)
				}
				iMessage.inputMessage.Ack()
				iMessage.inputMessage = nil
				iMessage.outPutMessage = nil
				putInternalMessage(iMessage)
			}, kqRecover())
		}
	}
	return nil
}

func (k *KQ) commitMessage(ctx context.Context) error {
	ticket := time.Tick(time.Millisecond * 500)
	for k.kqStatus.Load() == kqStatusRunning {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticket:
			runtime.Gosched()
		case commitMessage := <-k.commitMessageChanel:
			commitMessage.Lock()
			err := k.input.CommitMessage(ctx, commitMessage)
			if err != nil {
				commitMessage.Release()
				continue
			}
			commitMessage.Release()
		}
	}
	return nil
}

func (k *KQ) Close(ctx context.Context) error {
	k.kqStatus.Store(kqStatusClosed)
	return nil
}

type internalMessage struct {
	inputMessage     *core.InputMessage
	extractorMessage core.Message
	outPutMessage    *core.OutputMessage
}

func getInternalMessage() *internalMessage {
	iMessageInterface := internalMessagePool.Get()
	iMessage, ok := iMessageInterface.(*internalMessage)
	if !ok {
		return &internalMessage{}
	}
	return iMessage
}

func putInternalMessage(iMessage *internalMessage) *internalMessage {
	internalMessagePool.Put(iMessage)
	return iMessage
}
