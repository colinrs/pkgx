package gin

import (
	"context"
	"github.com/gin-contrib/pprof"
	rgin "github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"sync"

	"github.com/colinrs/pkgx/shutdown"
)
type Engine struct {
	*rgin.Engine
	pprof bool
	healthCheck bool
	models []Model
	lock sync.RWMutex
	resourceCleanup []func()
	server *http.Server
	serverPlugin ServerPlugin
	responses Responses
	hook shutdown.Hook
	errChan chan error
	errDecode Decode
}

// Default returns an Engine instance with the Logger and Recovery middleware already attached.
func Default() *Engine {
	debugPrintWARNINGDefault()
	reneging := rgin.Default()
	engine := &Engine{
		Engine:reneging,
		pprof:true,
		healthCheck:true,
		models: []Model{},
		hook: shutdown.NewHook(),
		errChan: make(chan error),
	}
	engine.serverPlugin = newDefaultServerPlugin(engine)
	engine.responses = newDefaultResponses()
	engine.errDecode = newDefaultDecode()
	return engine
}


func New() *Engine {
	debugPrintWARNINGNew()
	reneging := rgin.New()
	engine := &Engine{
		Engine:reneging,
		pprof:true,
		healthCheck:true,
		models: []Model{},
		hook: shutdown.NewHook(),
	}
	engine.serverPlugin = newDefaultServerPlugin(engine)
	engine.responses = newDefaultResponses()
	return engine
}


func (engine *Engine) Run(addr ...string) (err error) {
	defer func() { debugPrintError(err) }()
	engine.registerAPI()
	address := resolveAddress(addr)
	server := &http.Server{Addr: address, Handler: engine.Engine}
	engine.server = server
	engine.serverPlugin.BeforeRun(context.Background())
	go func() {
		debugPrint("Listening and serving HTTP on %s\n", address)
		debugPrint("Listening and serving HTTP err:%s\n", server.ListenAndServe().Error())
	}()
	resourceCleanups := []func(){
		func() {
			engine.serverPlugin.AfterStop(context.Background())
		},
	}
	resourceCleanups = append(resourceCleanups, engine.resourceCleanup...)
	engine.hook.Close(resourceCleanups...) // 资源清理, 关闭服务器


	return
}

func (engine *Engine) registerAPI(){
	if engine.healthCheck{
		engine.Engine.GET("/health", Health)
	}
	if engine.pprof{
		pprof.Register(engine.Engine)
	}
	if err := engine.registerModelToRouter();err!=nil{
		panic(err)
	}

}

type HandlerFunc interface {}

func (engine *Engine) Group(relativePath string, handlers ...HandlerFunc) *rgin.RouterGroup {
	return engine.Engine.Group(relativePath, wrapHandlers(engine,handlers...)...)
}

func (engine *Engine) Any(relativePath string, handlers ...HandlerFunc) {
	engine.Engine.Any(relativePath, wrapHandlers(engine,handlers...)...)
}

func (engine *Engine) GET(relativePath string, handlers ...HandlerFunc) {
	engine.Engine.GET(relativePath, wrapHandlers(engine,handlers...)...)
}

func (engine *Engine) POST(relativePath string, handlers ...HandlerFunc) {
	engine.Engine.POST(relativePath, wrapHandlers(engine,handlers...)...)
}

func (engine *Engine) DELETE(relativePath string, handlers ...HandlerFunc) {
	engine.Engine.DELETE(relativePath, wrapHandlers(engine,handlers...)...)
}

func (engine *Engine) PATCH(relativePath string, handlers ...HandlerFunc) {
	engine.Engine.PATCH(relativePath, wrapHandlers(engine,handlers...)...)
}

func (engine *Engine) PUT(relativePath string, handlers ...HandlerFunc) {
	engine.Engine.PUT(relativePath, wrapHandlers(engine,handlers...)...)
}

func (engine *Engine) OPTIONS(relativePath string, handlers ...HandlerFunc) {
	engine.Engine.OPTIONS(relativePath, wrapHandlers(engine,handlers...)...)
}

func (engine *Engine) HEAD(relativePath string, handlers ...HandlerFunc) {
	engine.Engine.HEAD(relativePath, wrapHandlers(engine,handlers...)...)
}

// RegisterResourceCleanupFunc
func (engine *Engine) RegisterResourceCleanupFunc(funcs ...func()) {
	engine.lock.Lock()
	engine.resourceCleanup = append(engine.resourceCleanup, funcs...)
	engine.lock.Unlock()
}


func wrapHandlers(engine *Engine, handlers ...HandlerFunc) []rgin.HandlerFunc {
	resultHandlers := make([]rgin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		resultHandlers[i] = wrapHandler(engine, handler)
	}
	return resultHandlers
}

func wrapHandler(engine *Engine, impl interface{}) rgin.HandlerFunc{
	handlerFuncType := reflect.TypeOf(impl)
	f := reflect.ValueOf(impl)

	if f.Type().NumIn() != 2 && f.Type().NumIn() != 1 {
		panic("handler must be func(ctx context.Context, req *ReqType) or func(c *gin.Context), ReqType is alternative name")
	}
	if f.Type().NumOut() != 2 {
		panic("response must be func() (res *ReqResponse, err error)")
	}

	var reqType reflect.Type
	if f.Type().NumIn() == 2 {
		reqType = handlerFuncType.In(1)
		if reqType.Kind() == reflect.Ptr {
			reqType = reqType.Elem()
		}
	}

	inner := func(in []reflect.Value) []reflect.Value {
		c := in[0].Interface().(*rgin.Context)
		var realParam []reflect.Value
		var err error
		ctx := c.Request.Context()
		realParam = append(realParam, reflect.ValueOf(ctx))
		if reqType != nil {
			reqParam := reflect.New(reqType)
			req := reqParam.Interface()
			err = c.ShouldBind(req)
			if err != nil {
				engine.responses.Response(c, "ShouldBind", err, engine.errDecode)
				c.Abort()
				return nil
			}
			if request, ok := req.(Request);ok{
				err = request.Validator(ctx)
			}
			if err != nil {
				engine.responses.Response(c, "Validator", err, engine.errDecode)
				c.Abort()
				return nil
			}
			realParam = append(realParam, reqParam)
		}
		ret := f.Call(realParam)
		if !ret[1].IsNil() {
			err = ret[1].Interface().(error)
		}
		engine.responses.Response(c, ret[0].Interface(), err, engine.errDecode)
		return nil

	}
	v := reflect.MakeFunc(reflect.ValueOf(func(c *rgin.Context) {}).Type(), inner)
	return v.Interface().(func(c *rgin.Context))
}