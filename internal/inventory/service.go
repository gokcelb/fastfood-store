package inventory

import (
	"errors"
	"log"
	"strings"
)

var (
	OutOfStockErr       = errors.New("item is out of stock")
	CategoryCreationErr = errors.New("category could not be created")
)

type Repository interface {
	NumberOfItems() int
	Stock(id int) int
	UpdateStock(id int, newQuantity int)
	Item(id int) Item
	Items() []Item
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo}
}

func (s *Service) Item(id int) (Item, error) {
	log.Println("entered Item method")
	qty := s.repo.Stock(id)
	if qty < 1 {
		log.Println("quantity below 0")
		return Item{}, OutOfStockErr
	}
	log.Println("quantity above 0")

	s.repo.UpdateStock(id, qty-1)
	return s.repo.Item(id), nil
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
			return nil, CategoryCreationErr
		}
		catalogue[category] = append(catalogue[category], item)
	}
	log.Println("done looping through items", catalogue)
	return catalogue, nil
}

func (s *Service) OrganizeItems() []Item {
	itemsList := s.repo.Items()
	log.Println("got items from repository", itemsList)

	organizedItemsList := make([]Item, 20)
	for _, item := range itemsList {
		organizedItemsList[item.ID-1] = item
	}
	return organizedItemsList
}
