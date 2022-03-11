package stall

import (
	"log"

	"github.com/gokcelb/point-of-sale/internal/inventory"
	"github.com/gokcelb/point-of-sale/pkg/publisher"
	"github.com/google/uuid"
)

type (
	InventoryService interface {
		SufficientStock(orderID int) bool
		UpdateItemStock(e interface{})
		OrganizedItems() []inventory.Item
	}

	DefaultService struct {
		inventoryService InventoryService
		publisher        publisher.Publisher
	}
)

func NewService(invSvc InventoryService, pub publisher.Publisher) *DefaultService {
	return &DefaultService{inventoryService: invSvc, publisher: pub}
}

func (s *DefaultService) OrganizedItems() []inventory.Item {
	return s.inventoryService.OrganizedItems()
}

func (s *DefaultService) ProcessOrder(orderIDs []int) {
	for _, orderID := range orderIDs {
		if !s.inventoryService.SufficientStock(orderID) {
			log.Println("skipped processing order for item id:", orderID)
			continue
		}
		s.publisher.Publish("stock", publisher.StockEvent{
			ID:     uuid.NewString(),
			ItemID: orderID,
		})
		log.Println("published new stock event")
	}
}
