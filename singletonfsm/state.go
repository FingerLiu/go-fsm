package singletonfsm

import "context"

type State struct {
	Name      string
	enterHook func(ctx context.Context, state string)
	exitHook  func(ctx context.Context, state string)
}

func (s *State) SetEnterHook(hook func(ctx context.Context, state string)) {
	s.enterHook = hook
}

func (s *State) SetExitHook(hook func(ctx context.Context, state string)) {
	s.exitHook = hook
}
