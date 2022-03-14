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
		Item(orderID int) inventory.Item
		UpdateItemStock(e interface{})
		OrganizedItems() []inventory.Item
	}

	PointOfSaleService interface {
		CombosAndNonCombos(orderedItems []inventory.Item) ([][]inventory.Item, []inventory.Item)
		ComboPrices(burgerCombos [][]inventory.Item) []float64
		NonComboPrices(nonCombos []inventory.Item) []float64
		TotalPrice() float64
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

type nested map[int]map[string]interface{}

func (n nested) structure(orders interface{}, orderPrices []float64) {
	if v, ok := orders.([][]inventory.Item); ok {
		log.Println("type combos -> structure combos")
		n.structureCombos(v, orderPrices)
	} else if v, ok := orders.([]inventory.Item); ok {
		log.Println("type items -> structure non combos")
		n.structureNonCombos(v, orderPrices)
	}
}

func (n nested) structureCombos(v [][]inventory.Item, p []float64) {
	for i := range v {
		n[i] = map[string]interface{}{}
		n[i]["combo"] = v[i]
		n[i]["price"] = p[i]
	}
	log.Println("structure combos", n)
}

func (n nested) structureNonCombos(v []inventory.Item, p []float64) {
	for i := range v {
		n[i] = map[string]interface{}{}
		n[i]["item"] = v[i]
		n[i]["price"] = p[i]
	}
	log.Println("structure non combos", n)
}

func (s *DefaultService) ProcessOrders(orderIDs []int) (nested, nested) {
	var orderedItems []inventory.Item
	for _, orderID := range orderIDs {
		item := s.inventoryService.Item(orderID)
		orderedItems = append(orderedItems, item)
	}

	combos, nonCombos := s.pointOfSaleService.CombosAndNonCombos(orderedItems)
	nonComboPrices := s.pointOfSaleService.NonComboPrices(nonCombos)
	comboPrices := s.pointOfSaleService.ComboPrices(combos)

	nonCombosWithPrices := make(nested)
	nonCombosWithPrices.structure(nonCombos, nonComboPrices)
	log.Println(nonCombosWithPrices)

	if combos == nil {
		log.Println("combos nil")
		return nil, nonCombosWithPrices
	}

	combosWithPrices := make(nested)
	combosWithPrices.structure(combos, comboPrices)
	log.Println(combosWithPrices)

	return combosWithPrices, nonCombosWithPrices
}

func (s *DefaultService) TotalPrice() float64 {
	return s.pointOfSaleService.TotalPrice()
}
