package handlers

import (
	"fmt"
	"log"
	"net/http"
)

func (cfg *App) PingDB(w http.ResponseWriter, r *http.Request) {
	err := cfg.DB.Ping()
	if err != nil {
		fmt.Fprintf(w, "Failed to reach DB. Err: %v", err)
		log.Printf("DB not reachable: %v", err)
		return
	}

	fmt.Fprintln(w, `{"status": "ok"}`)
	log.Default().Println("Status ok")
}
