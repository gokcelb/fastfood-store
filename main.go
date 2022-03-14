package main

import (
	"log"
	"time"

	"github.com/gokcelb/point-of-sale/internal/inventory"
	"github.com/gokcelb/point-of-sale/internal/pos"
	"github.com/gokcelb/point-of-sale/internal/stall"
	"github.com/gokcelb/point-of-sale/pkg/publisher"
	"github.com/gokcelb/point-of-sale/pkg/subscriber"
)

func main() {
	pub := publisher.NewStockPublisher()

	inventoryRepo := inventory.NewRepository()
	inventorySvc := inventory.NewService(inventoryRepo)

	pointOfSaleSvc := pos.NewService()

	stallSvc := stall.NewService(inventorySvc, pub, pointOfSaleSvc)
	stallCLI := stall.NewCLI(stallSvc)

	stallCLI.GiveCatalogue()
	stallCLI.TakeOrder()

	sub := subscriber.New(pub)
	go sub.Listen("stock", inventorySvc.UpdateItemStock)
	go sub.Listen("stock", printEventReceived)

	time.Sleep(30 * time.Second)
}

// see if two actions can listen to same events of same topic
func printEventReceived(ie interface{}) {
	_, ok := ie.(publisher.StockEvent)
	if !ok {
		log.Fatal("not of stock type")
	}
}
