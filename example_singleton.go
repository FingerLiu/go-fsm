package main

import (
	"context"
	fsm "github.com/FingerLiu/go-fsm/singletonfsm"
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

// OrderV2Service is a typical stateless service layer,
// it does not save any state, state should be passed through params or context
type OrderV2Service struct {
	fsm *fsm.FSM
}

type OrderV2 struct {
	ID        int
	Name      string
	OrderType OrderType
	Status    string
}

func NewOrderV2(id int, name string, orderType OrderType) *OrderV2 {
	orderV2 := &OrderV2{
		ID:        id,
		Name:      name,
		OrderType: orderType,
		Status:    OrderStatusCreated,
	}
	return orderV2
}

func NewOrderV2Service() *OrderV2Service {
	orderServiceV2 := &OrderV2Service{}
	orderFsm := fsm.NewFSM("orderV2").
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
		AddTransitionOn(OrderStatusCheckout, OrderStatusFinished, orderServiceV2.IsVirtual).
		// add transition on a condition
		AddTransitionOn(OrderStatusPaid, OrderStatusCancelled, orderServiceV2.IsPhysical).
		// add hook for a specific state(enter/exit)
		AddStateEnterHook(OrderStatusCancelled, orderServiceV2.stopDeliver).
		// global hook is triggered when state change(enter/exit) success.
		// here we use hook to save sate to order status field in database.
		AddGlobalEnterHook(orderServiceV2.saveStatus)

	orderServiceV2.fsm = orderFsm

	log.Printf("[order] order created %v\n", orderServiceV2)
	return orderServiceV2
}

// transit to destination state
func (o *OrderV2Service) Transit(ctx context.Context, from, to string) error {
	return o.fsm.Transit(ctx, from, to)
}

// output graphviz visualization
func (o *OrderV2Service) VisualizeFsm(filename string) {
	o.fsm.RenderGraphvizImage(filename)
}

func (o *OrderV2Service) saveStatus(ctx context.Context, status string) {
	// TODO save to database
	log.Printf("[order]saved order to db with status %s\n", status)
	orderV2 := ctx.Value("order").(*OrderV2)
	orderV2.Status = status
	return
}

func (o *OrderV2Service) IsPhysical(ctx context.Context, status string) (bool, error) {
	orderV2 := ctx.Value("order").(*OrderV2)
	res := (orderV2.OrderType) == OrderTypePhysical
	log.Printf("[order] IsPhysical return %v\n", res)
	return res, nil
}

func (o *OrderV2Service) IsVirtual(ctx context.Context, status string) (bool, error) {
	orderV2 := ctx.Value("order").(*OrderV2)
	res := (orderV2.OrderType) == OrderTypeVirtual
	log.Printf("[order] IsVirtual return %v\n", res)
	return res, nil
}

func (o *OrderV2Service) stopDeliver(ctx context.Context, status string) {
	// TODO call deliver sub system to stop
	log.Printf("[order] stop deliver order with status %s\n", status)
	return
}

func main() {
	log.Println("start singletonfsm test with order")
	orderV2Service := NewOrderV2Service()
	orderPhysical := NewOrderV2(1, "my_first_physical_order", OrderTypePhysical)
	orderVirtual := NewOrderV2(1, "my_first_physical_order", OrderTypeVirtual)

	log.Println("------ start transit physical order ------")
	ctx := context.Background()
	ctx1 := context.WithValue(ctx, "order", orderPhysical)
	orderV2Service.Transit(ctx1, orderPhysical.Status, OrderStatusPaid)
	orderV2Service.Transit(ctx1, orderPhysical.Status, OrderStatusCancelled)
	log.Printf("[order] order status is %s\n", orderPhysical.Status)

	log.Println("------ start transit virtual order ------")
	log.Println("[order] start fsm test with order")
	log.Printf("[order] order status is %s\n", orderVirtual.Status)
	ctx2 := context.WithValue(ctx, "order", orderVirtual)
	orderV2Service.Transit(ctx2, orderVirtual.Status, OrderStatusPaid)
	orderV2Service.Transit(ctx2, orderVirtual.Status, OrderStatusCancelled)
	log.Printf("[order] order status is %s\n", orderVirtual.Status)

	//orderV2Service.fsm.RenderGraphvizDot()
	orderV2Service.VisualizeFsm("./demoV2.png")
}
