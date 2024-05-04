package controllers

import (
	"context"
	"database/sql"
	"github.com/line/line-bot-sdk-go/v8/linebot"
	"golang.org/x/time/rate"
	"log"
	"time"
	"touyakun/models"
)

type NotificationConfig struct {
	bot *linebot.Client
}

func InitializeNotificationController(channelAccessToken, channelSecret string) *NotificationConfig {
	bot, _ := linebot.New(channelSecret, channelAccessToken)
	return &NotificationConfig{bot: bot}
}

func (nc NotificationConfig) NotificationController(db *sql.DB) {
	notificationRepo := models.InitializeNotificationRepo(db)

	nowTime := time.Now()

	notificationList, err := notificationRepo.GetNotificationList(nowTime)
	if err != nil {
		log.Println(err)
	}

	l := rate.NewLimiter(2000.0, 1)
	ctx := context.Background()
	for _, notification := range notificationList {
		if err := l.Wait(ctx); err != nil {
			log.Println(err)
		}
		sendMedicationNotificationForSpecifiedLineUser(notification.LineUserId, notification.DosageName, notification.DosageAmout, nc.bot)
	}
}

func sendMedicationNotificationForSpecifiedLineUser(lineUserId string, dosageName string, dosageAmount string, bot *linebot.Client) {
	message := linebot.NewTextMessage("服薬の時間です!\n" + dosageName + "を" + dosageAmount + "錠飲んでください！ 🎉")
	if _, err := bot.PushMessage(lineUserId, message).Do(); err != nil {
		log.Println(err)
	}
}
