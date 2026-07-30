package handlers

import (
	"context"
	"time"
)

func (cfg *App) Cleanup(ch <-chan time.Time) {
	for {
		select {
		case <-ch:
			cfg.deleteExpiredJTI(context.Background())
			cfg.deleteExpiredRefresh(context.Background())
		}
	}
}
