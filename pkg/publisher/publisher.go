package publisher

type Publisher interface {
	Publish(e StockEvent)
	Event() StockEvent
	EventQueueEmpty() bool
}

func NewStockPublisher() Publisher {
	return &StockPublisher{[]StockEvent{}}
}

type StockPublisher struct {
	eventQueue []StockEvent
}

type StockEvent struct {
	ID     string
	ItemID int
}

func (sp *StockPublisher) Publish(event StockEvent) {
	sp.eventQueue = append(sp.eventQueue, event)
}

func (sp *StockPublisher) Event() (event StockEvent) {
	event, sp.eventQueue = sp.eventQueue[0], sp.eventQueue[1:]
	return
}

func (sp *StockPublisher) EventQueueEmpty() bool {
	return len(sp.eventQueue) == 0
}
