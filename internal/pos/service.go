package pos

import (
	"log"
	"sort"

	"github.com/gokcelb/point-of-sale/internal/inventory"
)

type PointOfSale struct {
	total float64
}

func NewService() *PointOfSale {
	return &PointOfSale{total: 0}
}

func (pos *PointOfSale) TotalPrice() float64 {
	return pos.total
}

func (pos *PointOfSale) NonComboPrices(nonCombos []inventory.Item) (nonComboPrices []float64) {
	for _, item := range nonCombos {
		nonComboPrices = append(nonComboPrices, item.Price)
		pos.total += item.Price
	}
	return
}

func (pos *PointOfSale) ComboPrices(burgerCombos [][]inventory.Item) (comboPrices []float64) {
	for _, combo := range burgerCombos {
		var comboPrice float64
		for _, item := range combo {
			comboPrice += item.Price
		}
		discountedComboPrice := comboPrice - comboPrice*15/100
		comboPrices = append(comboPrices, discountedComboPrice)
		pos.total += discountedComboPrice
	}
	return
}

func (pos *PointOfSale) CombosAndNonCombos(orderedItems []inventory.Item) ([][]inventory.Item, []inventory.Item) {
	var (
		burgers []inventory.Item
		sides   []inventory.Item
		drinks  []inventory.Item
	)

	for _, item := range orderedItems {
		if item.Category == "burgers" {
			burgers = append(burgers, item)
		} else if item.Category == "sides" {
			sides = append(sides, item)
		} else {
			drinks = append(drinks, item)
		}
	}

	if len(burgers) == 0 || len(sides) == 0 || len(drinks) == 0 {
		log.Println("no combos found")
		return nil, orderedItems
	}

	var burgerCombos [][]inventory.Item
	for _, category := range [][]inventory.Item{burgers, sides, drinks} {
		pos.sortItemsByPrice(category)
	}

	shortest := pos.findShortestLength([][]inventory.Item{burgers, sides, drinks}, len(burgers))
	for i := 0; i < shortest; i++ {
		burgerCombos = append(burgerCombos, []inventory.Item{burgers[i], sides[i], drinks[i]})
	}

	var nonCombos []inventory.Item
	for _, category := range [][]inventory.Item{burgers, sides, drinks} {
		for _, item := range category {
			if pos.itemIsCombo(item, burgerCombos) {
				continue
			}
			nonCombos = append(nonCombos, item)
		}
	}

	return burgerCombos, nonCombos
}

func (*PointOfSale) itemIsCombo(item inventory.Item, combos [][]inventory.Item) bool {
	for _, combo := range combos {
		for _, comboItem := range combo {
			if item.ID == comboItem.ID {
				return true
			}
		}
	}
	return false
}

func (*PointOfSale) findShortestLength(a [][]inventory.Item, presumedShortest int) int {
	for _, v := range a {
		if len(v) < presumedShortest {
			presumedShortest = len(v)
		}
	}
	return presumedShortest
}

func (*PointOfSale) sortItemsByPrice(a []inventory.Item) {
	sort.Slice(a, func(i, j int) bool {
		return a[i].Price < a[j].Price
	})
}
