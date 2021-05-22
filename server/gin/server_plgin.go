package gin

import (
	"context"
	"time"
	"fmt"
)

type ServerPlugin interface {
	BeforeRun(ctx context.Context) error
	WhereStop(ctx context.Context) error
	AfterStop(ctx context.Context) error
}


var _ ServerPlugin = (*DefaultServerPlugin)(nil)


type DefaultServerPlugin struct {
	engine *Engine
}

func newDefaultServerPlugin(engine *Engine) *DefaultServerPlugin{
	return &DefaultServerPlugin{
		engine: engine,
	}
}


func (d *DefaultServerPlugin) BeforeRun(ctx context.Context) error {

	return nil
}

func (d *DefaultServerPlugin) WhereStop(ctx context.Context) error {
	return nil
}

func (d *DefaultServerPlugin) AfterStop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := d.engine.server.Shutdown(ctx)
	if err != nil {
		debugPrint("server shutdown error: %v", err.Error())
	}
	fmt.Print("DefaultServerPlugin AfterStop\n")
	return nil
}

func(engine *Engine) SetServerPlugin(serverPlugin ServerPlugin){
	engine.serverPlugin = serverPlugin
}