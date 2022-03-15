# fastfood-store

With this app, I tried to code __pub/sub pattern__ at the infrastructure level (without looking at any existing implementation so it is likely to be sloppy) and get familiar with the concept of __event-driven architecture__. In this way, I got to learn about concurrency in go (go routines, channels) as well.

## Pub/Sub Code

__fastfood-store__ has a pub/sub model to broadcast stock events belonging to a topic. A publisher publishes an event for a specific topic and subscribers can listen to this topic's events from the beginning of the event queue. The subscriber's `Listen` method takes a topic name as the first parameter, and a function that gets passed an event as a second parameter. In an infinite loop, the method checks whether the subscribed topic's event queue is full or not.

```golang
// subscriber.go
func (fss *FastfoodStoreSubscriber) Listen(topicName string, action func(e interface{})) {
	c := make(chan int)
	i := 0
	for {
		eq := fss.publisher.Topic(topicName).EventQueue
		for ; i < len(eq); i++ {
			action(eq[i])
		}
		go fss.publisher.ClearEventQueue(topicName, c)
		select {
		case i = <-c:
		}
	}
}
```

 The index value is reset to zero if the event queue is cleared by the publisher and we receive a value from the channel passed to the `ClearEventQueue` method.

```golang
// publisher.go
func (fsp *FastfoodStorePublisher) ClearEventQueue(topicName string, c chan int) {
	time.Sleep(20 * time.Second)
	fsp.Topic(topicName).EventQueue = make([]interface{}, 0)
	c <- 0
}
```

More than one event can listen to the same event queue.

```golang
// main.go
func printEventReceived(ie interface{}) {
	_, ok := ie.(publisher.StockEvent)
	if !ok {
		log.Fatal("not of stock type")
	}
}
```

Subscription works as shown below:

```golang
// main.go
go sub.Listen("stock", inventorySvc.UpdateItemStock)
go sub.Listen("stock", printEventReceived)
```

## Behavior

If you buy a combo (one item from each category makes a combo) there is a %15 discount. If user buys two items of the same category, the app counts the most expensive item as part of the combo in order to maximize user benefit.

```
----------BURGERS----------
1 Python Burger 5.99      
2 C Burger 4.99
3 Ruby Burger 6.49
4 Go Burger 5.99
5 C++ Burger 7.99
6 Java Burger 7.99

----------SIDES----------
7 Small Fries 2.49
8 Medium Fries 3.49
9 Large Fries 4.29
10 Small Caesar Salad 3.49
11 Large Caesar Salad 4.49

----------DRINKS----------
12 Small Coke 1.99
13 Medium Coke 2.49
14 Large Coke 2.99
15 Small Ginger Ale 1.99
16 Medium Ginger Ale 2.49
17 Large Ginger Ale 2.99
18 Small Chocolate Milkshake 3.99
19 Medium Chocolate Milkshake 4.49
20 Large Chocolate Milkshake 4.99

Please enter the number of the item that you would like to add to your order. Enter q to complete your order.

Enter an item number:
1

Enter an item number:
5

Enter an item number:
7

Enter an item number:
11

Enter an item number:
20

Enter an item number:
q

Placing order...

$14.85 Burger Combo
C++ Burger
Large Caesar Salad
Large Chocolate Milkshake

Non combos
Python Burger $5.99
Small Fries $2.49

Total: $23.33
Would you like to continue ordering? Type yes or no.
yes

Please enter the number of the item that you would like to add to your order. Enter q to complete your order.

Enter an item number:
14

Enter an item number:
q

Placing order...
2022/03/15 11:00:26 no combos found

Non combos
Large Coke $2.99

Total: $26.32
Would you like to continue ordering? Type yes or no.
no
Would you like to finish ordering and pay for your order? Type yes or no.
yes
That'll be $26.32. Come again!
```
