package gin

import (
	"github.com/gin-contrib/pprof"
	rgin "github.com/gin-gonic/gin"

	"net/http"
)
type Engine struct {
	*rgin.Engine
	pprof bool
	healthCheck bool
}

// Default returns an Engine instance with the Logger and Recovery middleware already attached.
func Default() *Engine {
	debugPrintWARNINGDefault()
	reneging := rgin.Default()
	engine := &Engine{
		reneging,
		true,
		true,
	}
	return engine
}


func New() *Engine {
	debugPrintWARNINGNew()
	reneging := rgin.New()
	engine := &Engine{
		reneging,
		true,
		true,
	}
	return engine
}


func (engine *Engine) Run(addr ...string) (err error) {
	defer func() { debugPrintError(err) }()

	address := resolveAddress(addr)
	debugPrint("Listening and serving HTTP on %s\n", address)
	if engine.healthCheck{
		engine.GET("/health", Health)
	}
	if engine.pprof{
		pprof.Register(engine.Engine)
	}
	err = http.ListenAndServe(address, engine.Engine)
	return
}


func(engine *Engine) AfterRun(){

}