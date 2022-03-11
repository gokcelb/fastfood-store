package stall

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gokcelb/point-of-sale/internal/inventory"
)

type Service interface {
	OrganizedItems() []inventory.Item
	ProcessOrder(orderNos []int)
}

type CLI struct {
	svc Service
}

func NewCLI(service Service) *CLI {
	return &CLI{service}
}

func (cli *CLI) GiveCatalogue() {
	categories := []string{"Burgers", "Sides", "Drinks"}
	var categoryStart bool
	// avoid initializing with zero value because
	// it disrupts if logic for burger category
	var categoriesIdx int = -1

	for _, item := range cli.svc.OrganizedItems() {
		name := strings.ToLower(item.Name)
		if strings.Contains(name, "burger") {
			categoryStart = categoriesIdx != 0
			categoriesIdx = 0
		} else if strings.Contains(name, "fries") || strings.Contains(name, "salad") {
			categoryStart = categoriesIdx != 1
			categoriesIdx = 1
		} else {
			categoryStart = categoriesIdx != 2
			categoriesIdx = 2
		}

		if categoryStart {
			fmt.Printf("\n----------%s----------\n", categories[categoriesIdx])
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
			break
		}

		orderNo, err := strconv.Atoi(order)
		if err != nil {
			fmt.Println("Please enter a valid number")
			continue
		}
		orderNos = append(orderNos, orderNo)
	}
	cli.svc.ProcessOrder(orderNos)
}
