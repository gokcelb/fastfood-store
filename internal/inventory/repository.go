package inventory

import "math/rand"

// TODO: create a publisher/subscriber mechanism so that services can subscribe to product price and stock changes

type MemDB map[int]Item

func NewRepository() MemDB {
	return map[int]Item{
		1:  {ID: 1, Name: "Python Burger", Price: 5.99, Quantity: rand.Intn(16) + 1, Category: "burgers"},
		2:  {ID: 2, Name: "C Burger", Price: 4.99, Quantity: rand.Intn(16) + 1, Category: "burgers"},
		3:  {ID: 3, Name: "Ruby Burger", Price: 6.49, Quantity: rand.Intn(16) + 1, Category: "burgers"},
		4:  {ID: 4, Name: "Go Burger", Price: 5.99, Quantity: rand.Intn(16) + 1, Category: "burgers"},
		5:  {ID: 5, Name: "C++ Burger", Price: 7.99, Quantity: rand.Intn(16) + 1, Category: "burgers"},
		6:  {ID: 6, Name: "Java Burger", Price: 7.99, Quantity: rand.Intn(16) + 1, Category: "burgers"},
		7:  {ID: 7, Name: "Small Fries", Price: 2.49, Quantity: rand.Intn(16) + 1, Category: "sides"},
		8:  {ID: 8, Name: "Medium Fries", Price: 3.49, Quantity: rand.Intn(16) + 1, Category: "sides"},
		9:  {ID: 9, Name: "Large Fries", Price: 4.29, Quantity: rand.Intn(16) + 1, Category: "sides"},
		10: {ID: 10, Name: "Small Caesar Salad", Price: 3.49, Quantity: rand.Intn(16) + 1, Category: "sides"},
		11: {ID: 11, Name: "Large Caesar Salad", Price: 4.49, Quantity: rand.Intn(16) + 1, Category: "sides"},
		12: {ID: 12, Name: "Small Coke", Price: 1.99, Quantity: rand.Intn(16) + 1, Category: "drinks"},
		13: {ID: 13, Name: "Medium Coke", Price: 2.49, Quantity: rand.Intn(16) + 1, Category: "drinks"},
		14: {ID: 14, Name: "Large Coke", Price: 2.99, Quantity: rand.Intn(16) + 1, Category: "drinks"},
		15: {ID: 15, Name: "Small Ginger Ale", Price: 1.99, Quantity: rand.Intn(16) + 1, Category: "drinks"},
		16: {ID: 16, Name: "Medium Ginger Ale", Price: 2.49, Quantity: rand.Intn(16) + 1, Category: "drinks"},
		17: {ID: 17, Name: "Large Ginger Ale", Price: 2.99, Quantity: rand.Intn(16) + 1, Category: "drinks"},
		18: {ID: 18, Name: "Small Chocolate Milkshake", Price: 3.99, Quantity: rand.Intn(16) + 1, Category: "drinks"},
		19: {ID: 19, Name: "Medium Chocolate Milkshake", Price: 4.49, Quantity: rand.Intn(16) + 1, Category: "drinks"},
		20: {ID: 20, Name: "Large Chocolate Milkshake", Price: 4.99, Quantity: rand.Intn(16) + 1, Category: "drinks"},
	}
}

func (db MemDB) Stock(id int) int {
	return db[id].Quantity
}

func (db MemDB) UpdateStock(id int, newQuantity int) {
	item := db[id]
	item.Quantity = newQuantity
	db[id] = item
}

func (db MemDB) Item(id int) Item {
	return db[id]
}

func (db MemDB) Items() (items []Item) {
	for _, item := range db {
		items = append(items, item)
	}
	return
}
