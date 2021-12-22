package main

import (
	"context"
	"github.com/FingerLiu/go-fsm/fsm"
	"log"
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

type Order struct {
	ID        int
	Name      string
	OrderType OrderType
	Status    string
	fsm       *fsm.FSM
}

func NewOrder(id int, name string, orderType OrderType) *Order {
	order := &Order{
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
		AddTransitionOn(OrderStatusCheckout, OrderStatusFinished, order.IsVirtual).
		// add transition on a condition
		AddTransitionOn(OrderStatusPaid, OrderStatusCancelled, order.IsPhysical).
		// add hook for a specific state(enter/exit)
		AddStateEnterHook(OrderStatusCancelled, order.stopDeliver).
		// global hook is triggered when state change(enter/exit) success.
		// here we use hook to save sate to order status field in database.
		AddGlobalEnterHook(order.saveStatus)

	orderFsm.SetState(order.Status)

	order.fsm = orderFsm

	log.Printf("[order] order created %v\n", order)
	return order
}

// force set state
func (o *Order) SetState(status string) error {
	return o.fsm.SetState(status)
}

// transit to destination state
func (o *Order) Transit(status string) error {
	return o.fsm.Transit(status)
}

func (o *Order) GetCurrentStatus() string {
	return o.fsm.GetCurrentState()
}

// output graphviz visualization
func (o *Order) VisualizeFsm(filename string) {
	o.fsm.RenderGraphvizImage(filename)
}

func (o *Order) saveStatus(ctx context.Context, status string) {
	// TODO save to database
	log.Printf("[order]saved order to db with status %s\n", status)
	return
}

func (o *Order) IsPhysical(ctx context.Context, status string) (bool, error) {
	res := (o.OrderType) == OrderTypePhysical
	log.Printf("[order] IsPhysical return %v\n", res)
	return res, nil
}

func (o *Order) IsVirtual(ctx context.Context, status string) (bool, error) {
	res := (o.OrderType) == OrderTypeVirtual
	log.Printf("[order] IsVirtual return %v\n", res)
	return res, nil
}

func (o *Order) stopDeliver(ctx context.Context, status string) {
	// TODO call deliver sub system to stop
	log.Printf("[order] stop deliver order with status %s\n", status)
	return
}

func main() {
	log.Println("start fsm test with order")
	order := NewOrder(1, "my_first_physical_order", OrderTypePhysical)
	log.Printf("[order] order status is %s\n", order.GetCurrentStatus())
	order.Transit(OrderStatusPaid)
	order.Transit(OrderStatusCancelled)
	log.Printf("[order] order status is %s\n", order.GetCurrentStatus())

	log.Println("------ start transit virtual order ------")
	log.Println("[order] start fsm test with order")
	orderVirtual := NewOrder(1, "my_first_physical_order", OrderTypeVirtual)
	log.Printf("[order] order status is %s\n", orderVirtual.GetCurrentStatus())
	orderVirtual.Transit(OrderStatusPaid)
	orderVirtual.Transit(OrderStatusCancelled)
	log.Printf("[order] order status is %s\n", orderVirtual.GetCurrentStatus())

	//order.fsm.RenderGraphvizDot()
	order.VisualizeFsm("./demo.png")
}
