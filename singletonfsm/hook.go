package singletonfsm

import (
	"context"
	"fmt"
)

func (f *FSM) AddStateEnterHook(state string, hook func(ctx context.Context, state string)) *FSM {
	s := f.getState(state)
	s.SetEnterHook(hook)
	return f
}

func (f *FSM) AddStateExitHook(state string, hook func(ctx context.Context, state string)) *FSM {
	s := f.getState(state)
	s.SetExitHook(hook)
	return f
}

// AddGlobalEnterHook and AddGlobalExitHook will be executed on every state
func (f *FSM) AddGlobalEnterHook(hook func(ctx context.Context, state string)) *FSM {
	f.globalEnterHook = hook
	return f
}

// AddGlobalEnterHook and AddGlobalExitHook will be executed on every state
func (f *FSM) AddGlobalExitHook(hook func(ctx context.Context, state string)) *FSM {
	f.globalExitHook = hook
	return f
}

func (f *FSM) executeHook(state *State, hook func(ctx context.Context, state string)) {
	if hook != nil {
		fmt.Printf("\t[fsm] start execute hook for state %s\n", state.Name)
		hook(f.ctx, state.Name)
	}
}

func (f *FSM) executeGlobalEnterHook(state *State) {
	if f.globalEnterHook != nil {
		f.executeHook(state, f.globalEnterHook)
	}
}

func (f *FSM) executeGlobalExitHook(state *State) {
	f.executeHook(state, f.globalExitHook)
}

func (f *FSM) executeEnterHook(state *State) {
	f.executeHook(state, state.enterHook)
}

func (f *FSM) executeExitHook(state *State) {
	f.executeHook(state, state.exitHook)
}
