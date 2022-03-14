package subscriber

import (
	"github.com/gokcelb/point-of-sale/pkg/publisher"
)

type (
	Subscriber interface {
		Listen(topicName string, action func(e interface{}))
	}

	FastfoodStoreSubscriber struct {
		publisher publisher.Publisher
	}
)

func New(pub publisher.Publisher) Subscriber {
	return &FastfoodStoreSubscriber{pub}
}

func (fss *FastfoodStoreSubscriber) Listen(topicName string, action func(e interface{})) {
	c := make(chan int)
	i := 0
	for {
		eq := fss.publisher.Topic(topicName).EventQueue
		for ; i < len(eq); i++ {
			action(eq[i])
		}
		go fss.publisher.ClearEventQueue(topicName, c)
		select {
		case i = <-c:
		}
	}
}
