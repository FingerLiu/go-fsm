package fsm

import "context"

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

// AddGlobalEnterHook and AddGlobalExitHook must be used in the end of building fsm chain
func (f *FSM) AddGlobalEnterHook(hook func(ctx context.Context, state string)) *FSM {
	for _, s := range f.states {
		s.SetEnterHook(hook)
	}
	return f
}

// AddGlobalEnterHook and AddGlobalExitHook must be used in the end of building fsm chain
func (f *FSM) AddGlobalExitHook(hook func(ctx context.Context, state string)) *FSM {
	for _, s := range f.states {
		s.SetExitHook(hook)
	}
	return f
}
