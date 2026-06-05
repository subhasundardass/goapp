package core

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
)

// Event is the envelope every event carries.
type Event struct {
	// Topic is the event name, e.g. "user.created", "order.placed".
	Topic string
	// Payload holds arbitrary event data; cast to the expected type in the handler.
	Payload any
	// Source is the name of the module that published the event.
	Source string
}

// EventHandler is a function that handles an event.
// It receives a context so handlers can respect cancellation / deadlines.
type EventHandler func(ctx context.Context, event Event) error

// subscription holds a handler and its mode.
type subscription struct {
	handler EventHandler
	async   bool // true = fire-and-forget goroutine; false = synchronous
}

// EventBus is a hybrid publish/subscribe bus.
// Sync subscribers are called in the publishing goroutine (blocking).
// Async subscribers are dispatched in a new goroutine (non-blocking).
type EventBus struct {
	mu     sync.RWMutex
	subs   map[string][]subscription
	logger *slog.Logger
}

// NewEventBus creates a new EventBus.
func NewEventBus(logger *slog.Logger) *EventBus {
	return &EventBus{
		subs:   make(map[string][]subscription),
		logger: logger,
	}
}

// Subscribe registers a synchronous handler for topic.
// The handler blocks the publisher until it returns.
func (b *EventBus) Subscribe(topic string, handler EventHandler) {
	b.add(topic, handler, false)
}

// SubscribeAsync registers an asynchronous handler for topic.
// The handler runs in its own goroutine; errors are logged, not returned.
func (b *EventBus) SubscribeAsync(topic string, handler EventHandler) {
	b.add(topic, handler, true)
}

func (b *EventBus) add(topic string, handler EventHandler, async bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subs[topic] = append(b.subs[topic], subscription{handler: handler, async: async})
}

// Publish sends event to all registered subscribers for event.Topic.
// Sync handlers are called inline; async handlers are dispatched as goroutines.
// Returns the first synchronous handler error, if any.
func (b *EventBus) Publish(ctx context.Context, event Event) error {
	b.mu.RLock()
	subs := make([]subscription, len(b.subs[event.Topic]))
	copy(subs, b.subs[event.Topic])
	b.mu.RUnlock()

	for _, sub := range subs {
		if sub.async {
			// Capture loop variable safely.
			s := sub
			go func() {
				if err := s.handler(ctx, event); err != nil {
					b.logger.Error("async event handler error",
						"topic", event.Topic,
						"source", event.Source,
						"error", err,
					)
				}
			}()
		} else {
			if err := sub.handler(ctx, event); err != nil {
				return fmt.Errorf("event %q handler error: %w", event.Topic, err)
			}
		}
	}
	return nil
}

// Topics returns all topics that have at least one subscriber (useful for debugging).
func (b *EventBus) Topics() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	topics := make([]string, 0, len(b.subs))
	for t := range b.subs {
		topics = append(topics, t)
	}
	return topics
}
