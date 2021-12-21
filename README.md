# go-fsm
An ease to use finit state machine golang implementation.Turn any struct to a fsm with graphviz visualization supported.

# usage

# full example

take order status as a full example
https://github.com/FingerLiu/go-fsm/blob/main/fsm/visualization.go

## define fsm
```go
import (
    "context"
    "fmt"
    "github.com/FingerLiu/go-fsm/fsm"
)

func NewOrder(id int, name string, orderType OrderType) *OrderService {
	orderService := &OrderService{
		ID:        id,
		Name:      name,
		OrderType: orderType,
		Status:    OrderStatusCreated,
	}
	ctx := context.Background()
	orderFsm := fsm.NewFSM(ctx, name).
		// add state to fsm
		AddStates(OrderStatusCreated, OrderStatusCancelled,
			OrderStatusPaid, OrderStatusCheckout,
			OrderStatusDelivering, OrderStatusDelivered, OrderStatusFinished).
		//add transition from S to E with condition check C
		AddTransition(OrderStatusCreated, OrderStatusCancelled).
		AddTransition(OrderStatusCreated, OrderStatusPaid).
		AddTransition(OrderStatusPaid, OrderStatusCheckout).
		AddTransition(OrderStatusCheckout, OrderStatusDelivering).
		AddTransition(OrderStatusDelivering, OrderStatusDelivered).
		AddTransition(OrderStatusDelivered, OrderStatusFinished).
		//virtual order do not need deliver
		AddTransitionOn(OrderStatusCheckout, OrderStatusFinished, orderService.IsVirtual).
		// add transition on a condition
		AddTransitionOn(OrderStatusPaid, OrderStatusCancelled, orderService.IsPhysical).
		// add hook for a specific state(enter/exit)
		AddStateEnterHook(OrderStatusCancelled, orderService.stopDeliver).
		// global hook is triggered when state change(enter/exit) success.
		// here we use hook to save sate to order status field in database.
		AddGlobalEnterHook(orderService.saveStatus)

	orderFsm.SetState(orderService.Status)

	orderService.fsm = orderFsm

	fmt.Printf("[order] order created %v\n", orderService)
	return orderService
}

```

## transit
```go
	fmt.Println("------ start transit virtual order ------")
	order := NewOrder(1, "my_first_physical_order", OrderTypePhysical)
	order.Transit(OrderStatusPaid)
	
	// this will fail because condition check not met
	// and order status will stay in paid
	order.Transit(OrderStatusCancelled)
	fmt.Printf("[order] order status is %s\n", order.GetCurrentStatus())

	fmt.Println("------ start transit physical order ------")
	
	orderVirtual := NewOrder(1, "my_first_physical_order", OrderTypeVirtual)
	orderVirtual.Transit(OrderStatusPaid)
	// this will pass
	orderVirtual.Transit(OrderStatusCancelled)
	fmt.Printf("[order] order status is %s\n", orderVirtual.GetCurrentStatus())

```
## visualization
```go
    // you can gen dot file or a png image
    order.fsm.RenderGraphvizDot()
    order.fsm.RenderGraphvizImage("./demo.png")
```
## singleton
If you don't want instance a fsm for every object, 
you can use singletonfsm.
A singleton fsm does not have concept of current/setState,
it serves as a stateless util.

```go

```
