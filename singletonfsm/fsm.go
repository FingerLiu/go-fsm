package singletonfsm

import (
	"context"
	"errors"
	"fmt"
	"log"
)

type FSM struct {
	name            string
	states          []*State
	transitions     []*Transition
	globalEnterHook func(ctx context.Context, state string)
	globalExitHook  func(ctx context.Context, state string)
}

func NewFSM(name string) *FSM {
	return &FSM{name: name}
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
		log.Fatalf("[fsm]state already defined %s", state)
		return nil
	}
	f.states = append(f.states, &State{Name: state})
	return f
}

func (f *FSM) AddStates(state ...string) *FSM {
	for _, s := range state {
		f.AddState(s)
	}
	return f
}

func (f *FSM) AddTransition(from, to string) *FSM {
	return f.AddTransitionOn(from, to, nil)
}

func (f *FSM) AddTransitionOn(from, to string, condition func(ctx context.Context, state string) (bool, error)) *FSM {
	if !f.hasState(from) {
		log.Fatalf("\t[fsm] state not defined %s", from)
		return nil
	}
	if !f.hasState(to) {
		log.Fatalf("\t[fsm] state not defined %s", to)
		return nil
	}
	if !f.hasTransition(from, to) {
		fromState := f.getState(from)
		toState := f.getState(to)
		f.transitions = append(f.transitions, NewTransition(fromState, toState, condition))
	} else {
		log.Printf("\t[fsm] Skipped add transition due to transition exists. from %s to %s\n",
			from, to)
	}

	return f
}

/***** transit fsm  *****/

// transit from current state to the given state
func (f *FSM) Transit(ctx context.Context, from, to string) error {
	availableTransitions := f.getAvailableTransitions(from)
	for _, transition := range availableTransitions {
		if transition.To.Name == to {
			log.Printf("\t[fsm] transit status to %s\n", to)
			return f.doTransit(ctx, transition)
		}
	}
	err := errors.New(fmt.Sprintf("\t[fsm] transition from %s to %s not found",
		from, to))
	log.Println(err.Error())
	return err
}

// check condition and set state
func (f *FSM) doTransit(ctx context.Context, transition *Transition) error {
	log.Printf("\t[fsm] start condition check for transit(%s)\n", transition.Key)
	if transition.Condition != nil {
		if flag, err := transition.Condition(ctx, transition.From.Name); err != nil {
			log.Printf("\t[fsm] transit(%s) condition check err %s\n", transition.Key, err)
			return err
		} else if flag == false {
			err = errors.New(fmt.Sprintf("[fsm] transit(%s) condition not met", transition.Key))
			log.Printf("\t[fsm] transit(%s) condition check err %s\n", transition.Key, err)
			return err
		}
	} else {
		log.Printf("\t[fsm] skipped condition check for transit(%s) due to condition is nil\n", transition.Key)
	}

	f.setState(ctx, transition.To)
	return nil
}

// setState will execute state hooks and global hooks
func (f *FSM) setState(ctx context.Context, state *State) {
	f.executeGlobalEnterHook(ctx, state)
	f.executeEnterHook(ctx, state)

	f.executeExitHook(ctx, state)
	f.executeGlobalExitHook(ctx, state)
}

/***** retrieve fsm  *****/

func (f *FSM) GetAvailableStateNames(from string) []string {
	states := f.getAvailableStates(from)
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
