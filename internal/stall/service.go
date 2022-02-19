package stall

import "github.com/gokcelb/point-of-sale/internal/inventory"

type InventoryService interface {
	Item(orderID int) (inventory.Item, error)
	Catalogue() (map[string][]inventory.Item, error)
}

type DefaultService struct {
	invSvc InventoryService
}

func NewService(invSvc InventoryService) *DefaultService {
	return &DefaultService{invSvc}
}

func (s *DefaultService) Catalogue() (map[string][]inventory.Item, error) {
	return s.invSvc.Catalogue()
}
