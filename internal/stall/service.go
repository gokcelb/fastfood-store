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
		Catalogue() (map[string][]inventory.Item, error)
	}

	DefaultService struct {
		inventoryService InventoryService
		publisher        publisher.Publisher
	}
)

func NewService(invSvc InventoryService, pub publisher.Publisher) *DefaultService {
	return &DefaultService{inventoryService: invSvc, publisher: pub}
}

func (s *DefaultService) Catalogue() (map[string][]inventory.Item, error) {
	return s.inventoryService.Catalogue()
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
