package ports

import "github.com/bryanarroyaveortiz/fashion-blue/internal/domain/events"

type EventPublisher interface {
	Publish(event events.OrderEvent)
}

type EventSubscriber interface {
	Subscribe(eventType events.OrderEventType, ch chan events.OrderEvent)
	Unsubscribe(eventType events.OrderEventType, ch chan events.OrderEvent)
}
