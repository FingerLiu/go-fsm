package fsm

import (
	"context"
	"fmt"
)

type Transition struct {
	From      *State
	To        *State
	Key       string
	Condition func(ctx context.Context, currentState string) (bool, error)
}

func NewTransition(from, to *State, condition func(ctx context.Context, currentState string) (bool, error)) *Transition {
	return &Transition{
		From:      from,
		To:        to,
		Key:       GenTransitionKey(from.Name, to.Name),
		Condition: condition,
	}
}

func GenTransitionKey(from, to string) string {
	return fmt.Sprintf("%s->%s", from, to)
}
