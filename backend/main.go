package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/thegreatestgiant/Charity-Tracker/internal/db"
	"github.com/thegreatestgiant/Charity-Tracker/internal/handlers"
)

func main() {
	envFile := ".env"
	if len(os.Args) > 1 {
		envFile = os.Args[1]
	}
	godotenv.Load(envFile)

	db, err := db.OpenDB(os.Getenv("DB_DEV_URL"))
	if err != nil {
		log.Fatal(err)
	}

	cfg := &handlers.App{
		DB:       db,
		JWT:      []byte(os.Getenv("JWT_SECRET")),
		Lifetime: time.Hour * 24,
	}

	ticker := time.NewTicker(cfg.Lifetime)
	go func() {
		for range ticker.C {
			cfg.Cleanup()
		}
	}()

	handlers.StartServer(cfg)
}
