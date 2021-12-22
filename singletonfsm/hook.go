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

func (f *FSM) executeHook(ctx context.Context, state *State, hook func(ctx context.Context, state string)) {
	if hook != nil {
		fmt.Printf("\t[fsm] start execute hook for state %s\n", state.Name)
		hook(ctx, state.Name)
	}
}

func (f *FSM) executeGlobalEnterHook(ctx context.Context, state *State) {
	if f.globalEnterHook != nil {
		f.executeHook(ctx, state, f.globalEnterHook)
	}
}

func (f *FSM) executeGlobalExitHook(ctx context.Context, state *State) {
	f.executeHook(ctx, state, f.globalExitHook)
}

func (f *FSM) executeEnterHook(ctx context.Context, state *State) {
	f.executeHook(ctx, state, state.enterHook)
}

func (f *FSM) executeExitHook(ctx context.Context, state *State) {
	f.executeHook(ctx, state, state.exitHook)
}
