package main

import (
	"time"

	"github.com/gokcelb/point-of-sale/internal/inventory"
	"github.com/gokcelb/point-of-sale/internal/stall"
	"github.com/gokcelb/point-of-sale/pkg/publisher"
)

func main() {
	p := publisher.NewStockPublisher()

	inventoryRepo := inventory.NewRepository()
	inventorySvc := inventory.New(inventoryRepo, p)
	go p.Listen(inventorySvc.UpdateItemStock)

	stallSvc := stall.NewService(inventorySvc, p)
	stallCLI := stall.New(stallSvc)
	stallCLI.GiveCatalogue()

	stallCLI.TakeOrder()

	time.Sleep(10 * time.Second)
}
