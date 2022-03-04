package publisher

type (
	Publisher interface {
		Publish(topicName string, event StockEvent)
		EventQueue(topicName string) EventQueue
	}

	FastfoodStorePublisher struct {
		topics []Topic
	}

	Topic struct {
		name       string
		eventQueue EventQueue
	}

	EventQueue []StockEvent

	StockEvent struct {
		ID     string
		ItemID int
	}
)

func NewStockPublisher() Publisher {
	return &FastfoodStorePublisher{[]Topic{}}
}

func (fsp *FastfoodStorePublisher) EventQueue(topicName string) EventQueue {
	return fsp.topic(topicName).eventQueue
}

func (fsp *FastfoodStorePublisher) topic(topicName string) *Topic {
	for _, topic := range fsp.topics {
		if topic.name == topicName {
			return &topic
		}
	}
	return newTopic(topicName)
}

func newTopic(topicName string) *Topic {
	return &Topic{
		name:       topicName,
		eventQueue: EventQueue{},
	}
}

func (fsp *FastfoodStorePublisher) Publish(topicName string, event StockEvent) {
	eq := fsp.topic(topicName).eventQueue
	eq = append(eq, event)
}
