package main

import (
	"github.com/gokcelb/point-of-sale/internal/inventory"
	"github.com/gokcelb/point-of-sale/internal/stall"
)

func main() {
	inventoryRepo := inventory.NewRepository()
	inventorySvc := inventory.New(inventoryRepo)

	stallSvc := stall.NewService(inventorySvc)
	stallCLI := stall.New(stallSvc)
	stallCLI.GiveCatalogue()
}
