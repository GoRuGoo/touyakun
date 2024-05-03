package main

import (
	"database/sql"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
	"os"
	"touyakun/controllers"
	"touyakun/router"
)

func main() {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=touyakun password=password dbname=touyakun sslmode=disable")
	if err != nil {
		log.Println(err)
	}

	//リクエスト単体で実行するのにルーターは必要ないのでcontrollerがrouter代わり
	nc := controllers.InitializeNotificationController(os.Getenv("CHANNEL_TOKEN"))
	c := cron.New()
	c.AddFunc("* * * * *", func() { nc.NotificationController(db) })
	c.Start()

	router.InitializeRouter(db)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
