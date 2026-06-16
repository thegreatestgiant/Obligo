package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/thegreatestgiant/Charity-Tracker/internal/DB"
	"github.com/thegreatestgiant/Charity-Tracker/internal/handlers"
)

func main() {
	envFile := ".env"
	if len(os.Args) > 1 {
		envFile = os.Args[1]
	}
	godotenv.Load(envFile)

	db, err := DB.OpenDB(os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Running migrations...")
	if err := DB.RunMigrations(db); err != nil {
		log.Fatalf("Could not migrate database: %v", err)
	}
	log.Println("Migrations applied successfully.")

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
