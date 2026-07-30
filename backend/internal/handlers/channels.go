package handlers

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type channel struct {
	Name string
	Val  float64
}

func (cfg *App) channelAmntOwed(ctx context.Context, user_id uuid.UUID, ch chan channel) {
	ch <- channel{Name: "owed", Val: cfg.getAmountOwed(ctx, user_id)}
}

func (cfg *App) channelAmntFulfilled(ctx context.Context, user_id uuid.UUID, ch chan channel) {
	ch <- channel{Name: "ful", Val: cfg.getAmountFulfilled(ctx, user_id)}
}

func (cfg *App) channelDonated(ctx context.Context, user_id uuid.UUID, ch chan channel) {
	ch <- channel{Name: "donated", Val: cfg.getAmountDonated(ctx, user_id)}
}

func (cfg *App) channelEarned(ctx context.Context, user_id uuid.UUID, ch chan channel) {
	ch <- channel{Name: "earned", Val: cfg.getAmountEarned(ctx, user_id)}
}

func (cfg *App) channelPercent(ctx context.Context, user_id uuid.UUID, ch chan channel) {
	ch <- channel{Name: "percent", Val: cfg.getDonationPercent(ctx, user_id)}
}

func (cfg *App) channelAll(ctx context.Context, user_id uuid.UUID) (owed float64, fulfilled float64, remaining float64, donated float64, earned float64, percent float64) {
	ch := make(chan channel, 6)

	go cfg.channelAmntOwed(ctx, user_id, ch)
	go cfg.channelAmntFulfilled(ctx, user_id, ch)
	go cfg.channelDonated(ctx, user_id, ch)
	go cfg.channelEarned(ctx, user_id, ch)
	go cfg.channelPercent(ctx, user_id, ch)

	for i := 0; i < 5; i++ {
		f := <-ch
		switch f.Name {
		case "owed":
			owed = f.Val
		case "ful":
			fulfilled = f.Val
		case "donated":
			donated = f.Val
		case "earned":
			earned = f.Val
		case "percent":
			percent = f.Val
		}
	}
	remaining = owed - ((math.Min(100, fulfilled) / 100) * owed)
	remaining = math.Round(remaining*100) / 100
	return
}

func (cfg *App) channelCsv(ctx context.Context, ch chan job, result chan jobResult, user_id uuid.UUID, percentage float64) {
	for job := range ch {

		row, index := job.data, job.rowIndex
		var parsedDate time.Time
		var parseErr error
		dateStr := row[4]
		if dateStr == "" {
			parsedDate = time.Now()
		} else {
			parsedDate, parseErr = time.Parse(time.RFC3339, dateStr)
			if parseErr != nil {
				errorMsg := fmt.Sprintf("Invalid date format on row %d: %v", index+2, parseErr)
				log.Println(errorMsg)
				result <- jobResult{message: errorMsg, code: http.StatusBadRequest}
				return
			}
		}

		amount, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			log.Printf("Invalid amount: %v", row[2])
			result <- jobResult{
				message: fmt.Sprintf("Invalid amount on row %d", index+2),
				code:    http.StatusBadRequest,
			}
			return
		}

		charity_owed := amount * percentage
		if row[1] == string(Donation) {
			charity_owed = 0
		}

		_, err = cfg.executeTemplate(
			ctx,
			`INSERT INTO Ledgers
			(user_id,
			ledger_entry,
			amount,
			description,
			charity_owed,
			transaction_date)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			user_id, row[1], amount, row[3], charity_owed, parsedDate,
		)
		if err != nil {
			log.Printf("Couldn't insert: %v", err)
			result <- jobResult{
				message: "Failed to insert record",
				code:    http.StatusInternalServerError,
			}
			return
		}
		result <- jobResult{success: true}
	}
}
