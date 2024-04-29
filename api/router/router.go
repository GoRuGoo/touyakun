package router

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func InitializeRouter() {
	db, err := sql.Open("postgres", "host=db port=5432 user=touyakun password=password dbname=touyakun sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	app, err := NewLINEConfig(db)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/callback", app.CallBackRouter)
}
