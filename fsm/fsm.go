package fsm

import (
	"context"
	"errors"
	"fmt"
)

type FSM struct {
	states []*State
	transitions []*Transition
	hooks []interface{}
	currentState *State
	ctx context.Context
}

func NewFSM(ctx context.Context) *FSM {
	return &FSM{ctx: ctx}
}

// build fsm
func (f *FSM) hasState(state string) bool {
	for _, s := range f.states {
		if s.Name == state {
			return true
		}
	}
	return false
}

func (f *FSM) getState(state string) *State {
	for _, s := range f.states {
		if s.Name == state {
			return s
		}
	}
	return nil
}

func (f *FSM) hasTransition(from, to string) bool {
	key := GenTransitionKey(from, to)
	for _, s := range f.transitions {
		if s.Key == key {
			return true
		}
	}
	return false
}

func (f *FSM) getTransition(from, to string) *Transition {
	key := GenTransitionKey(from, to)
	for _, s := range f.transitions {
		if s.Key == key {
			return s
		}
	}
	return nil
}

func (f *FSM) AddState(state string) *FSM {
	if f.hasState(state) {
		fmt.Errorf("state already defined %s", state)
		return nil
	}
	f.states = append(f.states, &State{Name: state})
	return f
}

func (f *FSM) AddStates(state... string) *FSM {
	for _, s := range state {
		f.AddState(s)
	}
	return f
}

func (f *FSM) AddTransition(from, to string) *FSM {
	return f.AddTransitionOn(from, to, nil)
}

func (f *FSM) AddTransitionOn(from, to string, condition func(ctx context.Context, state string)(bool, error)) *FSM {
	if !f.hasState(from) {
		fmt.Errorf("state not defined %s", from)
		return nil
	}
	if !f.hasState(to) {
		fmt.Errorf("state not defined %s", to)
		return nil
	}
	if !f.hasTransition(from, to) {
		fromState := f.getState(from)
		toState := f.getState(to)
		f.transitions = append(f.transitions, NewTransition(fromState, toState, condition))
	} else {
		fmt.Printf("Skipped add transition due to transition exists. from %s to %s\n",
			from, to)
	}

	return f
}

/***** transit fsm  *****/

// force set state without transit check
// will return err if state not in f.states
func (f *FSM) SetState(state string) error {
	s := f.getState(state)
	if s == nil {
		err := errors.New(fmt.Sprintf("state not defined %s", state))
		fmt.Print(err)
		return err
	}
	f.setState(s)
	return nil
}

// transit from current state to the given state
func (f *FSM) Transit(state string) error {
	availableTransitions := f.getAvailableTransitions(f.currentState.Name)
	for _, transition := range availableTransitions {
		if transition.To.Name == state {
			return f.doTransit(transition)
		}
	}
	return nil
}

// check condition and set state
func (f *FSM) doTransit(transition *Transition) error {
	fmt.Printf("start condition check for transit(%s)\n", transition.Key)
	if flag, err := transition.Condition(f.ctx, f.GetCurrentState()); err != nil {
		fmt.Printf("transit(%s) condition check err %s\n", transition.Key, err)
		return err
	}  else if flag == false {
		err = errors.New(fmt.Sprintf("transit(%s) condition not met", transition.Key))
		return err
	}

	f.setState(transition.To)
	return nil
}

// setState will execute state hooks and global hooks
func (f *FSM) setState(state *State)  {
	f.currentState = state
}

/***** retrieve fsm  *****/

func (f *FSM) GetCurrentState() string {
	return f.currentState.Name
}

func (f *FSM) GetAvailableStateNames() []string {
	state := f.GetCurrentState()
	states := f.getAvailableStates(state)
	names := make([]string, 0)
	for _, s := range states {
		names = append(names, s.Name)
	}
	return names
}

// only check transition link, do not check condition
func (f *FSM) getAvailableStates(state string) []*State {
	states := make([]*State, 0)
	for _, transition := range f.transitions {
		if transition.From.Name == state {
			states = append(states, transition.To)
		}
	}
	return states
}

// only check transition link, do not check condition
func (f *FSM) getAvailableTransitions(state string) []*Transition {
	transitions := make([]*Transition, 0)
	for _, transition := range f.transitions {
		if transition.From.Name == state {
			transitions = append(transitions, transition)
		}
	}
	return transitions
}