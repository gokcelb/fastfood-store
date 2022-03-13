package inventory

import (
	"errors"
	"log"

	"github.com/gokcelb/point-of-sale/pkg/publisher"
)

var (
	ErrOutOfStock       = errors.New("item is out of stock")
	ErrCategoryCreation = errors.New("category could not be created")
)

type Repository interface {
	Stock(id int) int
	UpdateStock(id int, newQuantity int)
	Item(id int) Item
	Items() []Item
}

type Service struct {
	repository Repository
}

func NewService(repo Repository) *Service {
	return &Service{repository: repo}
}

func (s *Service) SufficientStock(id int) bool {
	return s.repository.Stock(id) > 0
}

func (s *Service) Item(id int) Item {
	return s.repository.Item(id)
}

func (s *Service) UpdateItemStock(ei interface{}) {
	e, ok := ei.(publisher.StockEvent)
	if !ok {
		log.Fatal("not of stock type")
	}
	log.Printf("new event received for UpdateItemStock function, event id: %s, item id: %d", e.ID, e.ItemID)

	qty := s.repository.Stock(e.ItemID)

	s.repository.UpdateStock(e.ItemID, qty-1)
	log.Printf("updated stock: %d, old stock: %d", s.repository.Stock(e.ItemID), qty)
}

func (s *Service) OrganizedItems() []Item {
	itemsList := s.repository.Items()

	organizedItemsList := make([]Item, 20)
	for _, item := range itemsList {
		organizedItemsList[item.ID-1] = item
	}
	return organizedItemsList
}
