package utils

import (
	"net/http"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
)

func ReplyFlexCarouselMessage(bot *messaging_api.MessagingApiAPI, w http.ResponseWriter, replyToken string, contents []messaging_api.FlexBubble) {
	err := replyMessage(bot, w, replyToken, []messaging_api.MessageInterface{&messaging_api.FlexMessage{
		Contents: &messaging_api.FlexCarousel{
			Contents: contents,
		},
		AltText: "Flex message alt text",
	}})
	if err != nil {
		w.WriteHeader(500)
		return
	}
}

func ReplyTextMessage(bot *messaging_api.MessagingApiAPI, w http.ResponseWriter, replyToken string, textMessage *messaging_api.TextMessage) {
	err := replyMessage(bot, w, replyToken, []messaging_api.MessageInterface{
		textMessage,
	})
	if err != nil {
		w.WriteHeader(500)
		return
	}
}

func replyMessage(bot *messaging_api.MessagingApiAPI, w http.ResponseWriter, replyToken string, messages []messaging_api.MessageInterface) error {
	_, err := bot.ReplyMessage(
		&messaging_api.ReplyMessageRequest{
			ReplyToken: replyToken,
			Messages:   messages,
		},
	)
	if err != nil {
		w.WriteHeader(500)
		return err
	}
	return nil
}
