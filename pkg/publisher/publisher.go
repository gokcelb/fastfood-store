package publisher

import (
	"time"
)

type (
	Publisher interface {
		Publish(topicName string, event interface{})
		Topic(topicName string) *Topic
		ClearEventQueue(topicName string, c chan int)
	}

	FastfoodStorePublisher struct {
		topics []*Topic
	}

	Topic struct {
		Name       string
		EventQueue []interface{}
	}

	StockEvent struct {
		ID     string
		ItemID int
	}
)

func NewStockPublisher() Publisher {
	return &FastfoodStorePublisher{}
}

func (fsp *FastfoodStorePublisher) Topic(topicName string) *Topic {
	for _, Topic := range fsp.topics {
		if Topic.Name == topicName {
			return Topic
		}
	}
	return fsp.newTopic(topicName)
}

func (fsp *FastfoodStorePublisher) newTopic(topicName string) *Topic {
	newTopic := &Topic{
		Name: topicName,
	}
	fsp.topics = append(fsp.topics, newTopic)
	return newTopic
}

func (fsp *FastfoodStorePublisher) Publish(topicName string, event interface{}) {
	fsp.Topic(topicName).EventQueue = append(fsp.Topic(topicName).EventQueue, event)
}

func (fsp *FastfoodStorePublisher) ClearEventQueue(topicName string, c chan int) {
	time.Sleep(20 * time.Second)
	fsp.Topic(topicName).EventQueue = make([]interface{}, 0)
	c <- 0
}
