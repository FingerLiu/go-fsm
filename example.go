package main

import (
	"context"
	"fmt"
	"github.com/FingerLiu/go-fsm/fsm"
)

type OrderType string

const (
	OrderStatusCreated    = "created"
	OrderStatusCancelled  = "cancelled"
	OrderStatusPaid       = "paid"
	OrderStatusCheckout   = "checkout"
	OrderStatusDelivering = "delivering"
	OrderStatusDelivered  = "delivered"
	OrderStatusFinished   = "finished"

	// you can not cancel a virtual order once it is paid.

	// normal flow for a physical order maybe: created -> paid -> checkout -> delivering -> delivered -> finished
	OrderTypePhysical OrderType = "physical"

	// normal flow for a virtual order maybe: created -> paid -> finished
	OrderTypeVirtual OrderType = "virtual"
)

type OrderService struct {
	ID        int
	Name      string
	OrderType OrderType
	Status    string
	fsm       *fsm.FSM
}

func NewOrder(id int, name string, orderType OrderType) *OrderService {
	orderService := &OrderService{
		ID:        id,
		Name:      name,
		OrderType: orderType,
		Status:    OrderStatusCreated,
	}
	ctx := context.Background()
	orderFsm := fsm.NewFSM(ctx).
		// add state to fsm
		AddStates(OrderStatusCreated, OrderStatusCancelled,
			OrderStatusPaid, OrderStatusCheckout,
			OrderStatusDelivering, OrderStatusDelivered, OrderStatusFinished).
		//add transition from S to E with condition check C
		AddTransition(OrderStatusCreated, OrderStatusCancelled).
		// add transition on a condition
		AddTransitionOn(OrderStatusPaid, OrderStatusCancelled, orderService.IsPhysical).
		// add hook for a specific state(enter/exit)
		AddStateEnterHook(OrderStatusCancelled, orderService.stopDeliver).
		// global hook is triggered when state change(enter/exit) success.
		// here we use hook to save sate to order status field in database.
		AddGlobalEnterHook(orderService.saveStatus)

	orderService.fsm = orderFsm
	fmt.Printf("order created %v\n", orderService)
	return orderService
}

// force set state
func (o *OrderService) SetState(status string) error {
	return nil
}

// transit to destination state
func (o *OrderService) Transit(status string) error {
	return o.fsm.Transit(status)
}

func (o *OrderService) GetCurrentStatus() string {
	return o.fsm.GetCurrentState()
}

// output graphviz visualization
func (o *OrderService) VisualizeFsm() string {
	return o.fsm.Visualization()
}

func (o *OrderService) saveStatus(ctx context.Context, status string) {
	// TODO save to database
	return
}

func (o *OrderService) IsPhysical(ctx context.Context, status string) (bool, error) {
	return (o.OrderType) == OrderTypeVirtual, nil
}

func (o *OrderService) stopDeliver(ctx context.Context, status string) {
	// TODO save to database
	return
}

func main() {
	fmt.Println("start fsm test with order")
	order := NewOrder(1, "my_first_physical_order", OrderTypePhysical)
	order.SetState(OrderStatusCancelled)
}
