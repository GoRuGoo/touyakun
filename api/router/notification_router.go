package router

import (
	"log"
	"net/http"
)

type NotificationConfig struct {
	channelAccessToken string
}

func NewNotificationConfig(channelAccessToken string) *NotificationConfig {
	return &NotificationConfig{channelAccessToken: channelAccessToken}
}

func (nc NotificationConfig) HandleSendMessage(handler func(r *http.Request)) {
	req, _ := http.NewRequest("POST", "https://api.line.me/v2/bot/message/push", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+nc.channelAccessToken)

	handler(req)
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
