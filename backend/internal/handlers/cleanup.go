package handlers

import "time"

func (cfg *App) Cleanup(ch <-chan time.Time) {
	for {
		select {
		case <-ch:
			cfg.deleteExpiredJTI()
			cfg.deleteExpiredRefresh()
		}
	}
}
