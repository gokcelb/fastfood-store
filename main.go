package main

import (
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

	actions := subscriber.EventDrivenActions{inventorySvc.UpdateItemStock}
	sub := subscriber.NewStockSubscriber(pub, actions)
	go sub.Listen()

	time.Sleep(1 * time.Second)
}
