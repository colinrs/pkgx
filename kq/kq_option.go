package kq

const (
	defaultInputMessageChannelSize     = 1000
	defaultMaxGoroutines               = 10000
	defaultCommitMessageChannelSize    = 1000
	defaultExtractorMessageChannelSize = 1000
	defaultOutPutMessageChannelSize    = 1000
)

type options struct {
	inputMessageChannelSize     int
	maxGoroutines               int
	commitMessageChannelSize    int
	extractorMessageChannelSize int
	outPutMessageChannelSize    int
}

type Option func(*options)

func WithOptions(opts ...Option) Option {
	return func(o *options) {
		for _, opt := range opts {
			opt(o)
		}
	}
}

func WithInputMessageChannelSize(inputMessageChannelSize int) Option {
	return func(o *options) {
		o.inputMessageChannelSize = inputMessageChannelSize
	}
}

func WithCommitMessageChannelSize(commitMessageChannelSize int) Option {
	return func(o *options) {
		o.commitMessageChannelSize = commitMessageChannelSize
	}
}

func WithMaxGoroutines(maxGoroutines int) Option {
	return func(o *options) {
		o.maxGoroutines = maxGoroutines
	}
}

func WithExtractorMessageChannelSize(extractorMessageChannelSize int) Option {
	return func(o *options) {
		o.extractorMessageChannelSize = extractorMessageChannelSize
	}
}

func WithOutPutMessageChannelSize(outPutMessageChannelSize int) Option {
	return func(o *options) {
		o.outPutMessageChannelSize = outPutMessageChannelSize
	}
}
