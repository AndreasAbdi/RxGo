package observable

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/reactivex/rxgo/rx"
	"github.com/stretchr/testify/assert"
)

func TestImplementsEmitter(t *testing.T) {
	assert.Implements(t, (*rx.Emitter)(nil), &Emitter{})
}

func TestConstructEmitterFromContext(t *testing.T) {
	NewEmitter(context.Background())
}

func TestEmitterFromCancellableContext(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	emitter := NewEmitter(ctx)
	done := emitter.Done()
	cancelFunc()
	select {
	case <-done:
		assert.EqualError(t, emitter.Err(), context.Canceled.Error())
	default:
		t.Error("cancelling the context did not signal the emitter")
	}
}

func TestEmitterFromTimeoutContext(t *testing.T) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(0)) //nolint
	emitter := NewEmitter(ctx)
	done := emitter.Done()
	select {
	case <-done:
		assert.EqualError(t, emitter.Err(), context.DeadlineExceeded.Error())
	default:
		t.Error("cancelling the context did not signal the emitter")
	}
	cancelFunc()
}

func TestEmitterSubscribeAt(t *testing.T) {
	ctx := context.Background()
	emitter := NewEmitter(ctx)
	beforeSubscription := time.Now()
	emitter.Subscribe()
	afterSubscription := time.Now()

	emitterSubscriptionTime := emitter.SubscribeAt()
	assert.True(t, emitterSubscriptionTime.After(beforeSubscription), "time less than a time evaluation before subscription")
	assert.True(t, emitterSubscriptionTime.Before(afterSubscription), "time greater than time evaluation after subscription")
}

func TestEmitterDeadline(t *testing.T) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(0))
	emitter := NewEmitter(ctx)
	deadline, deadlineSet := emitter.Deadline()
	afterDeadlineTime := time.Now()
	done := emitter.Done()
	select {
	case <-done:
		assert.True(t, deadlineSet, "deadline should have been set")
		assert.True(t, deadline.Before(afterDeadlineTime), "deadline should have occured already")
	default:
		t.Error("cancelling the context did not signal the emitter")
	}
	cancelFunc()
}

func TestEmitterSettingError(t *testing.T) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(0))
	emitter := NewEmitter(ctx)
	newFakeError := errors.New("some new error")
	emitter.Error = newFakeError
	done := emitter.Done()
	select {
	case <-done:
		assert.EqualError(t, emitter.Err(), newFakeError.Error())
	default:
		t.Error("cancelling the context did not signal the emitter")
	}
	cancelFunc()
}
