package main

import (
	"log"
	"time"

	"github.com/gokcelb/point-of-sale/internal/inventory"
	"github.com/gokcelb/point-of-sale/internal/stall"
	"github.com/gokcelb/point-of-sale/pkg/publisher"
	"github.com/gokcelb/point-of-sale/pkg/subscriber"
)

func main() {
	pub := publisher.NewStockPublisher()

	inventoryRepo := inventory.NewRepository()
	inventorySvc := inventory.New(inventoryRepo)

	stallSvc := stall.NewService(inventorySvc, pub)
	stallCLI := stall.New(stallSvc)

	stallCLI.GiveCatalogue()
	stallCLI.TakeOrder()

	sub := subscriber.New(pub)
	go sub.Listen("stock", inventorySvc.UpdateItemStock)
	go sub.Listen("stock", printEventReceived)

	time.Sleep(30 * time.Second)
}

// see if two actions can listen to same events of same topic
func printEventReceived(ie interface{}) {
	e, ok := ie.(publisher.StockEvent)
	if !ok {
		log.Fatal("not of stock type")
	}
	log.Printf("new stock event received for printEventReceived function, event id: %s, item id: %d", e.ID, e.ItemID)
}
