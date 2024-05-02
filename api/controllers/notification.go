package controllers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type NotificationConfig struct {
	channelAccessToken string
}

func InitializeNotificationController(channelAccessToken string) *NotificationConfig {
	return &NotificationConfig{channelAccessToken: channelAccessToken}
}

func (nc NotificationConfig) NotificationController() {

	requestBody, err := json.Marshal(map[string]interface{}{
		"to":       os.Getenv("GORU_ID"),
		"messages": []map[string]string{{"type": "text", "text": "goru"}},
	})

	req, _ := http.NewRequest("POST", "https://api.line.me/v2/bot/message/push", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+nc.channelAccessToken)

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
