package pos

import (
	"sort"

	"github.com/gokcelb/point-of-sale/internal/inventory"
)

type PointOfSale struct{}

func (*PointOfSale) ComboPrices(burgerCombos [][]inventory.Item) []float64 {
	comboPrices := make([]float64, len(burgerCombos))
	for _, combo := range burgerCombos {
		var comboPrice float64
		for _, item := range combo {
			comboPrice += item.Price
		}
		discountedComboPrice := comboPrice - comboPrice*15/100
		comboPrices = append(comboPrices, discountedComboPrice)
	}
	return comboPrices
}

func (pos *PointOfSale) Combos(orderedItems []inventory.Item) (burgerCombos [][]inventory.Item) {
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
		return nil
	}

	pos.sortItemsByPrice(burgers)
	pos.sortItemsByPrice(sides)
	pos.sortItemsByPrice(drinks)

	shortest := pos.findShortestLength([][]inventory.Item{burgers, sides, drinks}, len(burgers))

	for i := 0; i < shortest; i++ {
		burgerCombos = append(burgerCombos, []inventory.Item{burgers[i], sides[i], drinks[i]})
	}
	return burgerCombos
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
