package fsm

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
	currentState    *State
	ctx             context.Context
}

func NewFSM(ctx context.Context, name string) *FSM {
	return &FSM{ctx: ctx, name: name}
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

// force set state without transit check
// will return err if state not in f.states
func (f *FSM) SetState(state string) error {
	s := f.getState(state)
	if s == nil {
		err := errors.New(fmt.Sprintf("\t[fsm] state not defined %s", state))
		log.Print(err)
		return err
	}
	log.Printf("\t[fsm] set status to %s\n", state)
	f.setState(s)
	return nil
}

// transit from current state to the given state
func (f *FSM) Transit(state string) error {
	availableTransitions := f.getAvailableTransitions(f.currentState.Name)
	for _, transition := range availableTransitions {
		if transition.To.Name == state {
			log.Printf("\t[fsm] transit status to %s\n", state)
			return f.doTransit(transition)
		}
	}
	err := errors.New(fmt.Sprintf("\t[fsm] transition from %s to %s not found", f.currentState.Name, state))
	log.Println(err.Error())
	return err
}

// check condition and set state
func (f *FSM) doTransit(transition *Transition) error {
	log.Printf("\t[fsm] start condition check for transit(%s)\n", transition.Key)
	if transition.Condition != nil {
		if flag, err := transition.Condition(f.ctx, f.GetCurrentState()); err != nil {
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

	f.setState(transition.To)
	return nil
}

// setState will execute state hooks and global hooks
func (f *FSM) setState(state *State) {
	f.executeGlobalEnterHook(state)
	f.executeEnterHook(state)
	f.currentState = state
	f.executeExitHook(state)
	f.executeGlobalExitHook(state)
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
