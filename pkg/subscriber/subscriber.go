package subscriber

import (
	"time"

	"github.com/gokcelb/point-of-sale/pkg/publisher"
)

type Subscriber interface {
	Listen()
}

func NewStockSubscriber(pub publisher.Publisher, act EventDrivenActions) Subscriber {
	return &StockSubscriber{publisher: pub, actions: act}
}

type StockSubscriber struct {
	publisher publisher.Publisher
	actions   EventDrivenActions
}

type EventDrivenActions []func(e publisher.StockEvent)

func (ss *StockSubscriber) Listen() {
	for {
		if !ss.publisher.EventQueueEmpty() {
			for _, action := range ss.actions {
				action(ss.publisher.Event())
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}
