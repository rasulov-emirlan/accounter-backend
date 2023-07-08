package shutdown

import "context"

type CloseFunc func(ctx context.Context) error

type Scheduler struct {
	stack []CloseFunc
}

func NewScheduler() Scheduler {
	return Scheduler{
		stack: make([]CloseFunc, 0),
	}
}

func (s *Scheduler) Add(f CloseFunc) {
	s.stack = append(s.stack, f)
}

func (s *Scheduler) Close(ctx context.Context) error {
	// last in first out
	for i := len(s.stack) - 1; i >= 0; i-- {
		if err := s.stack[i](ctx); err != nil {
			return err
		}
	}

	return nil
}
