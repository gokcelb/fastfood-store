package stall

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gokcelb/point-of-sale/internal/inventory"
)

type Service interface {
	OrganizedItems() []inventory.Item
	ProcessOrderStocks(orderNos []int)
	ProcessOrders(orderNos []int) (combosWPrices, nonCombosWPrices nested)
	TotalPrice() float64
}

type CLI struct {
	svc Service
}

func NewCLI(service Service) *CLI {
	return &CLI{service}
}

func (cli *CLI) GiveCatalogue() {
	var (
		currCategory  string
		categoryStart bool
	)
	for _, item := range cli.svc.OrganizedItems() {
		categoryStart = currCategory != item.Category
		if categoryStart {
			fmt.Printf("\n----------%s----------\n", strings.ToUpper(item.Category))
			currCategory = item.Category
		}
		fmt.Println(item.ID, item.Name, item.Price)
	}
}

func (cli *CLI) TakeOrder() {
	fmt.Println("\nPlease enter the number of items that you would like to add to your order. Enter q to complete your order.")
	var (
		order    string
		orderNos []int
	)
	for {
		fmt.Println("\nEnter an item number:")
		fmt.Scanln(&order)
		if order == "q" {
			fmt.Println("\nPlacing order...")
			break
		}

		orderNo, err := strconv.Atoi(order)
		if err != nil {
			fmt.Println("Please enter a valid number")
			continue
		}
		orderNos = append(orderNos, orderNo)
	}
	cli.svc.ProcessOrderStocks(orderNos)
	combosWPrices, nonCombosWPrices := cli.svc.ProcessOrders(orderNos)
	cli.giveBill(combosWPrices, nonCombosWPrices)
}

func (cli *CLI) giveBill(combosWithPrices, nonCombosWithPrices nested) {
	var combo []inventory.Item
	var comboPrice float64
	for _, val := range combosWithPrices {
		if v, ok := val["combo"].([]inventory.Item); ok {
			combo = v
		}
		if v, ok := val["price"].(float64); ok {
			comboPrice = v
		}
		fmt.Printf("\n$%0.2f Burger Combo\n", comboPrice)
		for _, item := range combo {
			fmt.Printf("%s\n", item.Name)
		}
	}

	fmt.Println("\nNon combos")
	var nonCombo inventory.Item
	var nonComboPrice float64
	for _, val := range nonCombosWithPrices {
		if v, ok := val["item"].(inventory.Item); ok {
			nonCombo = v
		}
		if v, ok := val["price"].(float64); ok {
			nonComboPrice = v
		}
		fmt.Printf("%s $%0.2f\n", nonCombo.Name, nonComboPrice)
	}

	fmt.Printf("Total: $%0.2f\n", cli.svc.TotalPrice())
}
