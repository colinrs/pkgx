package concurrent

// Limit ...
type Limit struct {
	bufSize int
	channel chan struct{}
}

// NewLimit ...
func NewLimit(concurrencyNum int) *Limit {
	return &Limit{channel: make(chan struct{}, concurrencyNum), bufSize: concurrencyNum}
}

// TryAcquire ...
func (limit *Limit) TryAcquire() bool {
	select {
	case limit.channel <- struct{}{}:
		return true
	default:
		return false
	}
}

// Acquire ...
func (limit *Limit) Acquire() {
	limit.channel <- struct{}{}
}

// Release ...
func (limit *Limit) Release() {
	<-limit.channel
}

// AvailablePermits ...
func (limit *Limit) AvailablePermits() int {
	return limit.bufSize - len(limit.channel)
}
