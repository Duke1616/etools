package schedule

import (
	"context"
	"fmt"
)

func (s *Scheduler) Schedule(ctx context.Context) error {
	for {
		dbCtx, cancel := context.WithTimeout(ctx, s.dbTimeout)
		j, err := s.svc.Preempt(dbCtx)
		cancel()
		if err != nil {
			continue
		}

		fmt.Print(j.Id)
	}

}
