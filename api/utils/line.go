package utils

import (
	"fmt"
	"net/http"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
)

func ReplyFlexCarouselMessage(bot *messaging_api.MessagingApiAPI, w http.ResponseWriter, replyToken string, contents []messaging_api.FlexBubble) {
	err := replyMessage(bot, w, replyToken, &messaging_api.FlexMessage{
		Contents: &messaging_api.FlexCarousel{
			Contents: contents,
		},
		AltText: "Flex message alt text",
	})
	if err != nil {
		w.WriteHeader(500)
		return
	}
}

func ReplyTextMessage(bot *messaging_api.MessagingApiAPI, w http.ResponseWriter, replyToken string, textMessage *messaging_api.TextMessage) {
	replyMessage(bot, w, replyToken, *textMessage)
}

func ReplyTemplateMessage(bot *messaging_api.MessagingApiAPI, w http.ResponseWriter, replyToken string, template messaging_api.TemplateInterface) {
	replyMessage(bot, w, replyToken, messaging_api.TemplateMessage{AltText: "TemplateMessage alt text", Template: template})
}

func replyMessage(bot *messaging_api.MessagingApiAPI, w http.ResponseWriter, replyToken string, message messaging_api.MessageInterface) error {
	_, err := bot.ReplyMessage(
		&messaging_api.ReplyMessageRequest{
			ReplyToken: replyToken,
			Messages:   []messaging_api.MessageInterface{message},
		},
	)
	if err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
		return err
	}
	return nil
}
