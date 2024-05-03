package controllers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"time"
	"touyakun/models"
)

type NotificationConfig struct {
	channelAccessToken string
}

func InitializeNotificationController(channelAccessToken string) *NotificationConfig {
	return &NotificationConfig{channelAccessToken: channelAccessToken}
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
		sendMedicationNotificationForSpecifiedLineUser(notification.LineUserId, notification.DosageName, notification.DosageAmout, nc.channelAccessToken)
	}
}

func sendMedicationNotificationForSpecifiedLineUser(lineUserId string, dosageName string, dosageAmount string, channelAccessToken string) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"to":       lineUserId,
		"messages": []map[string]string{{"type": "text", "text":"服薬の時間です!\n"+ dosageName + "を" + dosageAmount + "錠飲んでください！ 🎉"}},
	})

	req, _ := http.NewRequest("POST", "https://api.line.me/v2/bot/message/push", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+channelAccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("failed to send message")
	}

	log.Println("success to send message")
}
