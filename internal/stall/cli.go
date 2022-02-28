package stall

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gokcelb/point-of-sale/internal/inventory"
)

type Service interface {
	Catalogue() (map[string][]inventory.Item, error)
	ProcessOrder(orderNos []int)
}

type CLI struct {
	svc Service
}

func New(service Service) *CLI {
	return &CLI{service}
}

func (cli *CLI) GiveCatalogue() {
	catalogue, err := cli.svc.Catalogue()
	if err != nil {
		log.Fatal(err)
	}

	for cat, items := range catalogue {
		fmt.Printf("\n---------- %s ----------\n", cat)
		for _, item := range items {
			fmt.Printf("%d. %s $%0.2f\n", item.ID, item.Name, item.Price)
		}
	}
}

func (s *CLI) TakeOrder() {
	fmt.Println("Please enter the number of items that you would like to add to your order. Enter q to complete your order.")
	var (
		order    string
		orderNos []int
	)
	for {
		fmt.Println("Enter an item number:")
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
	s.svc.ProcessOrder(orderNos)
}
