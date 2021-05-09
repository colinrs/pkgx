package fx

// Parallel runs fns parallel and waits for done.
func Parallel(fns ...func()) {
	group := NewRoutineGroup()
	for _, fn := range fns {
		group.RunGoSafe(fn)
	}
	group.Wait()
}
