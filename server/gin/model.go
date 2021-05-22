package gin

type Model interface {
	Init(engine *Engine) error
}

func (engine *Engine) registerModelToRouter()  error{
	var err error
	for _, model := range engine.models{
		if err = model.Init(engine);err!=nil{
			return err
		}
	}
	return nil
}

func (engine *Engine) RegisterModel(models []Model)  {
	engine.lock.Lock()
	engine.models = append(engine.models, models...)
	engine.lock.Unlock()
}
