package observable

import (
	"context"
	"time"
)

// Emitter emits events asynchronously.
type Emitter struct {
	ctx              context.Context
	subscriptionTime time.Time
	Error            error
}

// NewEmitter constructs a new emitter
func NewEmitter(ctx context.Context) *Emitter {
	return &Emitter{
		ctx: ctx,
	}
}

//Deadline return when the emitter is completed, and whether or not a deadline is set.
func (e *Emitter) Deadline() (deadline time.Time, ok bool) {
	return e.ctx.Deadline()
}

// Done returns a channel that signals when the emitter has completed.
func (e *Emitter) Done() <-chan struct{} {
	return e.ctx.Done()
}

// Err returns the error that has occured for the emitter.
func (e *Emitter) Err() error {
	if e.Error != nil {
		return e.Error
	}
	return e.ctx.Err()
}

// Value assigns a key to the emitter.
func (e *Emitter) Value(key interface{}) interface{} {
	return e.ctx.Value(key)
}

// Subscribe records the time of subscription.
func (e *Emitter) Subscribe() {
	e.subscriptionTime = time.Now()
}

// SubscribeAt returns the time which the subscription occurred.
func (e Emitter) SubscribeAt() time.Time {
	return e.subscriptionTime
}
