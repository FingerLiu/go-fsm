# go-fsm
An ease to use finit state machine golang implementation.Turn any struct to a fsm with graphviz visualization supported.

# usage

TODO

# full example

take order status as a full example

```golang
package main

import (
	"context"
	"github.com/fingerliu/go-fsm/fsm"
)

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
	OrderTypePhysical = "physical"

	// normal flow for a virtual order maybe: created -> paid -> finished
	OrderTypeVirtual = "virtual"
)

type OrderService struct {
	Name   string
	Type   string
	Status string
	fsm    *fsm.FSM
}

func NewOrder() *OrderService {
	orderService := &OrderService{}
	ctx := context.Background()
	orderFsm := fsm.NewFSM(ctx).

		// add state to fsm
		AddStates(OrderStatusCreated, OrderStatusCancelled, OrderStatusPaid, OrderStatusCheckout, OrderStatusDelivering, OrderStatusDelivered, OrderStatusFinished).

		//add trasition from S to E with condition check C
		AddTransition(OrderStatusCreated, OrderStatusCancelled).

		// add trasition on a condition
		AddTransitionOn(OrderStatusPaid, OrderStatusCancelled, orderService.IsPhysical).

		// add hook for a specific state(enter/exit)
		AddStateEnterHook(OrderStatusCancelled, orderService.stopDeliver).

		// global hook is triggerred when state change(enter/exit) success.
		// here we use hook to save sate to order status field in database.
		AddGlobalEnterHook(orderService.saveStatus)

	orderService.fsm = orderFsm

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

func (o *OrderService) saveStatus(status string) error {
	// TODO save to database
	return nil
}

func (o *OrderService) IsPhysical(status string) bool {
	return o.Type == OrderTypeVirtual
}

func (o *OrderService) stopDeliver(status string) error {
	// TODO save to database
	return nil
}

```
