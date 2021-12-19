# go-fsm
An ease to use finit state machine golang implementation.Turn any struct to a fsm with graphviz visualization supported.

# usage

```golang
import github.com/fingerliu/go-fsm/fsm

const (
    OrderStatusCreated = "created"
    OrderStatusCancelled = "cancelled"
    OrderStatusPaid = "paid"
    OrderStatusCheckout = "checkout"
    OrderStatusDelivering = "delivering"
    OrderStatusDelivered = "delivered"
    OrderStatusFinished = "finished"
    
    // you can not cancel a virtual order once it is paid.
    
    // normal flow for a physical order maybe: created -> paid -> checkout -> delivering -> delivered -> finished
    OrderTypePhysical = "physical"
    
    // normal flow for a virtual order maybe: created -> paid -> finished
    OrderTypeVirtual = "virtual"
)

struct OrderService {
    Name string
    Type string
    Status string
    fsm fsm.FSM
}

func NewOrder() *OrderService {


    orderService := OrderService{}
    orderFsm := fsm.NewFSM().
    
        // add state to fsm
        AddStates(OrderStatusCreated,OrderStatusCancelled,OrderStatusPaid,OrderStatusCheckout,OrderStatusDelivering,OrderStatusDelivered,OrderStatusFinished).
    
        //add trasition from S to E with condition check C
        AddTransition(OrderStatusCreated, OrderStatusCancelled).
        
        // add trasition on a condition
        AddTransitionOn(OrderStatusPaid, OrderStatusCancelled, orderService.IsPhysical)
        
        // add hook for a specific state(enter/exit)
        AddStateEnterHook(OrderStatusCancelled, orderService.StopDeliver)
        
        // global hook is triggerred when state change(enter/exit) success.
        // here we use hook to save sate to order status field in database.
        AddGloalEnterHook(orderService.saveStatus)

    orderService.fsm = orderFsm
    
    return &orderService
}


```
