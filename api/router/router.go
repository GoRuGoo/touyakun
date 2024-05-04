package router

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

func InitializeRouter(db *sql.DB) {
	app, err := NewLINEConfig(os.Getenv("CHANNEL_SECRET"), os.Getenv("CHANNEL_TOKEN"), db)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", app.CallBackRouter)
}
