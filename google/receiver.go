package google

import (
	"context"
	"errors"
	"sort"

	googlePubSub "cloud.google.com/go/pubsub"

	"github.com/elmagician/pubsub"
)

const callbackSliceDefaultLength = 5 // default size to initialize callback filter slice.

var _ pubsub.Receiver = (*Receiver)(nil)

type Receiver struct {
	subscription *googlePubSub.Subscription

	// Local preparation
	messagesCallback map[int][]msgCallback
	checkersLength   []int
	unmatched        pubsub.MessageCallback
	onError          func(err error, msg pubsub.Message)

	// Running instance
	running         bool
	reloadOnFailure bool
	stop            func()
}

type msgCallback struct {
	envelop  pubsub.Envelop
	callback pubsub.MessageCallback
}

// OnMessage applies provided callback function when pubsub.Envelop pubsub.MessageFilter matches
// received message.
//
// Filter appliance rules
// - Must precise filter applies first
// - First in, First applied
//
// Provided again the same pubsub.Envelop will overwrite known callback.
func (r *Receiver) OnMessage(envelop pubsub.Envelop, callback pubsub.MessageCallback) {
	nbConditions := len(envelop.Filter())

	if _, ok := r.messagesCallback[nbConditions]; !ok {
		r.messagesCallback[nbConditions] = make([]msgCallback, callbackSliceDefaultLength)
		r.checkersLength = append(r.checkersLength, nbConditions)
	}

	r.messagesCallback[nbConditions] = append(
		r.messagesCallback[nbConditions], msgCallback{envelop: envelop, callback: callback},
	)
}

// OnUnmatched indicates how to react if message are not matched using known filters.
// Providing the key multiple times will overwrite callback.
// /!\ If no callback is provided for OnUnmatched, msg will
// be not acknowledge by default.
func (r *Receiver) OnUnmatched(callback pubsub.MessageCallback) {
	r.unmatched = callback
}

// OnError indicates how to react if received got an error.
// /!\ an instance cannot receive message any more after an error
// if  reloadOnFailure is not set. (cf Start)
// /!\ If no callback is provided for OnError, msg will
// be not acknowledge by default.
func (r *Receiver) OnError(callback func(err error, msg pubsub.Message)) {
	r.onError = callback
}

// Start receives message in a background process.
// It will listen until an error occurs or Stop is called.
// To reload receive process on error, set reloadOnFailure to true.
func (r *Receiver) Start(ctx context.Context, reloadOnFailure bool) {
	r.running = true
	r.reloadOnFailure = reloadOnFailure

	ctx, r.stop = context.WithCancel(ctx)

	// processReceive returns a boolean to knows if an error occurs
	// and can be relaunched.
	// receive method does not return context Errors
	processReceive := func() bool {
		if err := r.receive(ctx); err != nil {
			if r.onError != nil {
				r.onError(err, nil)
			}

			return true
		}

		return false
	}

	if r.reloadOnFailure {
		go func() {
			processReceive() // has we are not relaunching process, we do not care if an error happened here.
			r.Stop()
		}()
	} else {
		go func() {
			if hasError := processReceive(); hasError {
				r.Start(ctx, true)
			}
		}()
	}
}

// Receive receives message in a foreground process. It will stop and
// return on the first error met.
// When using Receive, OnError callback will never be called.
func (r *Receiver) Receive(ctx context.Context) error {
	ctx, r.stop = context.WithCancel(ctx)
	r.running = true

	return r.receive(ctx)
}

// Stop stops receiver.
func (r *Receiver) Stop() {
	r.running = false
	r.stop()
}

func (r *Receiver) receive(ctx context.Context) error {
	sort.Ints(r.checkersLength)

	err := r.subscription.Receive(ctx, func(ctx context.Context, gcpMessage *googlePubSub.Message) {
		msg := &message{Message: gcpMessage}

		for _, nbCheck := range r.checkersLength {
			for _, msgCallback := range r.messagesCallback[nbCheck] {
				filters := msgCallback.envelop.Filter()
				if matchesMetadata(gcpMessage.Attributes, filters) {
					newEnvelop := msgCallback.envelop.New()

					if err := newEnvelop.FromPubsubMessage(msg); err != nil {
						if r.onError != nil {
							r.onError(err, msg)
						} else {
							msg.Nack()
						}
					}

					msgCallback.callback(ctx, msg)
				}
			}
		}

		if r.unmatched != nil {
			r.unmatched(ctx, msg)
		} else {
			msg.Nack()
		}
	})

	if err != nil && errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}

func matchesMetadata(metaData map[string]string, filter map[string]string) bool {
	for key, expected := range filter {
		actual, exists := metaData[key]
		if !exists {
			return false
		}

		if actual != expected {
			return false
		}
	}

	return true
}
