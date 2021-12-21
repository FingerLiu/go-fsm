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

// force set state
func (o *OrderService) SetState(status string) error {
	return o.fsm.SetState(status)
}

// transit to destination state
func (o *OrderService) Transit(status string) error {
	return o.fsm.Transit(status)
}

func (o *OrderService) GetCurrentStatus() string {
	return o.fsm.GetCurrentState()
}

// output graphviz visualization
func (o *OrderService) VisualizeFsm(filename string) {
	o.fsm.RenderGraphvizImage(filename)
}

func (o *OrderService) saveStatus(ctx context.Context, status string) {
	// TODO save to database
	fmt.Printf("[order]saved order to db with status %s\n", status)
	return
}

func (o *OrderService) IsPhysical(ctx context.Context, status string) (bool, error) {
	res := (o.OrderType) == OrderTypePhysical
	fmt.Printf("[order] IsPhysical return %v\n", res)
	return res, nil
}

func (o *OrderService) IsVirtual(ctx context.Context, status string) (bool, error) {
	res := (o.OrderType) == OrderTypeVirtual
	fmt.Printf("[order] IsVirtual return %v\n", res)
	return res, nil
}

func (o *OrderService) stopDeliver(ctx context.Context, status string) {
	// TODO call deliver sub system to stop
	fmt.Printf("[order] stop deliver order with status %s\n", status)
	return
}

func main() {
	fmt.Println("start fsm test with order")
	order := NewOrder(1, "my_first_physical_order", OrderTypePhysical)
	fmt.Printf("[order] order status is %s\n", order.GetCurrentStatus())
	order.Transit(OrderStatusPaid)
	order.Transit(OrderStatusCancelled)
	fmt.Printf("[order] order status is %s\n", order.GetCurrentStatus())

	fmt.Println("------ start transit virtual order ------")
	fmt.Println("[order] start fsm test with order")
	orderVirtual := NewOrder(1, "my_first_physical_order", OrderTypeVirtual)
	fmt.Printf("[order] order status is %s\n", orderVirtual.GetCurrentStatus())
	orderVirtual.Transit(OrderStatusPaid)
	orderVirtual.Transit(OrderStatusCancelled)
	fmt.Printf("[order] order status is %s\n", orderVirtual.GetCurrentStatus())

	order.fsm.RenderGraphvizDot()
	order.VisualizeFsm("./demo.png")
}
