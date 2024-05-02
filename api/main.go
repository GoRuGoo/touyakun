package main

import (
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
	"os"
	"touyakun/controllers"
	"touyakun/router"
)

func main() {
	//リクエスト単体で実行するのにルーターは必要ないのでcontrollerがrouter代わり
	nc := controllers.InitializeNotificationController(os.Getenv("CHANNEL_TOKEN"))

	c := cron.New()
	c.AddFunc("@every 5s", nc.NotificationController)
	c.Start()

	router.InitializeRouter()
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
