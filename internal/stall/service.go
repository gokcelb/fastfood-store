package stall

import (
	"fmt"
	"log"

	"github.com/gokcelb/point-of-sale/internal/inventory"
	"github.com/gokcelb/point-of-sale/pkg/publisher"
	"github.com/google/uuid"
)

type (
	InventoryService interface {
		SufficientStock(orderID int) bool
		Item(orderID int) inventory.Item
		UpdateItemStock(e interface{})
		OrganizedItems() []inventory.Item
	}

	PointOfSaleService interface {
		CombosAndNonCombos(orderedItems []inventory.Item) ([][]inventory.Item, []inventory.Item)
		ComboPrices(burgerCombos [][]inventory.Item) []float64
		NonComboPrices(nonCombos []inventory.Item) []float64
	}

	DefaultService struct {
		inventoryService   InventoryService
		pointOfSaleService PointOfSaleService
		publisher          publisher.Publisher
	}
)

func NewService(invSvc InventoryService, pub publisher.Publisher, pos PointOfSaleService) *DefaultService {
	return &DefaultService{inventoryService: invSvc, publisher: pub, pointOfSaleService: pos}
}

func (s *DefaultService) OrganizedItems() []inventory.Item {
	return s.inventoryService.OrganizedItems()
}

func (s *DefaultService) ProcessOrderStocks(orderIDs []int) {
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

func (s *DefaultService) ProcessOrders(orderIDs []int) (map[int]map[string]interface{}, map[int]map[string]interface{}) {
	log.Println("process orders ran")
	var orderedItems []inventory.Item
	for _, orderID := range orderIDs {
		item := s.inventoryService.Item(orderID)
		orderedItems = append(orderedItems, item)
	}

	combos, nonCombos := s.pointOfSaleService.CombosAndNonCombos(orderedItems)
	nonCombosWPrices := make(map[int]map[string]interface{})
	nonComboPrices := s.pointOfSaleService.NonComboPrices(nonCombos)
	for i := range nonCombos {
		log.Println("map assigning NON COMBOS...")
		nonCombosWPrices[i] = map[string]interface{}{}
		nonCombosWPrices[i]["item"] = nonCombos[i]
		nonCombosWPrices[i]["price"] = nonComboPrices[i]
	}

	if combos == nil {
		log.Println("combos nil")
		return nil, nonCombosWPrices
	}

	combosWPrices := make(map[int]map[string]interface{})
	comboPrices := s.pointOfSaleService.ComboPrices(combos)
	log.Println("combo prices:", comboPrices)
	for i := range combos {
		log.Println("map assigning COMBOS...")
		combosWPrices[i] = map[string]interface{}{}
		combosWPrices[i]["combo"] = combos[i]
		combosWPrices[i]["price"] = comboPrices[i]
	}

	fmt.Printf("\n\n%+v\n%+v\n\n", combosWPrices, nonCombosWPrices)
	return combosWPrices, nonCombosWPrices
}
