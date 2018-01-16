package fakes

type Throttler struct {
	TooManyReceived bool
	FinishCalled    bool
}

func (t *Throttler) Throttle() bool {
	return t.TooManyReceived
}

func (t *Throttler) Finish() {
	t.FinishCalled = true
}
