package fx

// Parallel runs fns parallel and waits for done.
func Parallel(fns ...func()) {
	group := NewRoutineGroup()
	for _, fn := range fns {
		group.RunSafe(fn)
	}
	group.Wait()
}

// ParallelWithReturn runs fns parallel and waits for done. and return result
func ParallelWithReturn(fns ...func()) ([]interface{}, error){
	group := NewRoutineGroup()
	for _, fn := range fns {
		group.RunSafe(fn)
	}
	group.Wait()
	return nil, nil
}
