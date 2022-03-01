package inventory

import (
	"errors"
	"log"
	"strings"

	"github.com/gokcelb/point-of-sale/pkg/publisher"
)

var (
	ErrOutOfStock       = errors.New("item is out of stock")
	ErrCategoryCreation = errors.New("category could not be created")
)

type Repository interface {
	NumberOfItems() int
	Stock(id int) int
	UpdateStock(id int, newQuantity int)
	Item(id int) Item
	Items() []Item
}

type Service struct {
	repository Repository
}

func New(repo Repository) *Service {
	return &Service{repository: repo}
}

func (s *Service) SufficientStock(id int) bool {
	return s.repository.Stock(id) > 0
}

func (s *Service) UpdateItemStock(e publisher.StockEvent) {
	log.Printf("new event received. event id: %s, item id: %d", e.ID, e.ItemID)

	qty := s.repository.Stock(e.ItemID)
	if qty < 1 {
		log.Println("quantity below 1, quantity:", qty)
	}
	log.Println("quantity:", qty)

	s.repository.UpdateStock(e.ItemID, qty-1)
	log.Println("updated stock:", s.repository.Stock(e.ItemID))
	log.Println("stock updated successfully")
}

var keywords = map[string][]string{
	"Burgers": {"burger"},
	"Sides":   {"fries", "salad"},
	"Drinks":  {"coke", "ale", "milk"},
}

func keywordCategory(name string) string {
	for cat, keywords := range keywords {
		for _, keyword := range keywords {
			if strings.Contains(strings.ToLower(name), keyword) {
				return cat
			}
		}
	}
	return ""
}

func (s *Service) Catalogue() (map[string][]Item, error) {
	log.Println("entered Catalogue method")
	catalogue := map[string][]Item{
		"Burgers": {},
		"Sides":   {},
		"Drinks":  {},
	}
	items := s.OrganizeItems()
	log.Println("organized items")

	for _, item := range items {
		category := keywordCategory(item.Name)
		log.Println("ran keywordCategory function, category: ", category)
		if len(category) == 0 {
			log.Println("category length 0")
			return nil, ErrCategoryCreation
		}
		catalogue[category] = append(catalogue[category], item)
	}
	log.Println("done looping through items", catalogue)
	return catalogue, nil
}

func (s *Service) OrganizeItems() []Item {
	itemsList := s.repository.Items()
	log.Println("got items from repository", itemsList)

	organizedItemsList := make([]Item, 20)
	for _, item := range itemsList {
		organizedItemsList[item.ID-1] = item
	}
	return organizedItemsList
}
