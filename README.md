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
)

struct OrderService {
    Name string
    Status string
    fsm fsm.FSM
}

func NewOrder() *OrderService {
    orderFsm := fsm.NewFSM().
    
    // add state to fsm
    AddStates().
    //add trasition  from S to E with condition check C
    AddTransition().

    orderService := OrderService {
        fsm: orderFsm
    }
    return &orderService
}


```
