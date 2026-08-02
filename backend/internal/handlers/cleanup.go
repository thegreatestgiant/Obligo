package handlers

import (
	"context"
	"time"
)

func (cfg *App) Cleanup(ch <-chan time.Time) {
	for range ch {
		cfg.deleteExpiredJTI(context.Background())
		cfg.deleteExpiredRefresh(context.Background())
	}
}
