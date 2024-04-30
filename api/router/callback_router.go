package router

import (
	"database/sql"
	"errors"
	"net/http"
	"net/url"
	"touyakun/controllers"
	"touyakun/models"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

type LINEConfig struct {
	channelSecret string
	bot           *messaging_api.MessagingApiAPI
	blob          *messaging_api.MessagingApiBlobAPI
	db            *sql.DB
}

func NewLINEConfig(channelSecret, channelToken string, db *sql.DB) (*LINEConfig, error) {
	bot, err := messaging_api.NewMessagingApiAPI(channelToken)
	if err != nil {
		return nil, err
	}
	blob, err := messaging_api.NewMessagingApiBlobAPI(channelToken)
	if err != nil {
		return nil, err
	}

	return &LINEConfig{
		channelSecret: channelSecret,
		bot:           bot,
		blob:          blob,
		db:            db,
	}, nil
}

func (app *LINEConfig) CallBackRouter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	cb, err := webhook.ParseRequest(app.channelSecret, r)
	if err != nil {
		if errors.Is(err, webhook.ErrInvalidSignature) {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	userModel := models.InitializeUserRepo(app.db)
	dosageModel := models.InitializeDosageRepo(app.db)
	userController := controllers.InitializeUserController(userModel, w)

	for _, event := range cb.Events {
		switch e := event.(type){
		case webhook.FollowEvent:
			switch s := e.Source.(type) {
			case webhook.UserSource:
				// ユーザーが友達追加した時の処理
				isNotExist, err := userModel.IsNotExistUser(s.UserId)
				if err != nil {
					w.WriteHeader(500)
					w.Write([]byte(err.Error()))
					return
				}

				if !isNotExist {
					w.WriteHeader(409)
					w.Write([]byte("user already exists"))
					return
				}

				err = userModel.RegisterUser(s.UserId)
				if err != nil {
					w.WriteHeader(500)
					w.Write([]byte(err.Error()))
					return
				}
			}
		case webhook.UnfollowEvent:
			switch s := e.Source.(type) {
			case webhook.UserSource:
				// ユーザーが友達追加した時の処理
				isNotExist, err := userModel.IsNotExistUser(s.UserId)
				if err != nil {
					w.WriteHeader(500)
					w.Write([]byte(err.Error()))
					return
				}

				if isNotExist {
					w.WriteHeader(404)
					w.Write([]byte("user does not exist"))
					return
				}

				err = userModel.DeleteUser(s.UserId)
				if err != nil {
					w.WriteHeader(500)
					w.Write([]byte(err.Error()))
					return
				}
			}
		case webhook.PostbackEvent:
			data := e.Postback.Data
			// dataはURLパラメータの形式で書かれているので、パースする
			// 例: "action=buy&itemid=123"
			u, err := url.ParseQuery(data)
			if err != nil {
				w.WriteHeader(500)
				return
			}
			switch u.Get("action") {
			case "delete":
				// 薬の一覧を取得
				s := e.Source.(webhook.UserSource)
				medications, err := dosageModel.GetMedications(s.UserId)
				if err != nil {
					w.WriteHeader(500)
					return
				}
				//ユーザーにどの薬を消すかFlex Messageを使って質問
				contents:=&messaging_api.FlexCarousel{
					Contents: []messaging_api.FlexBubble{
						{
							Body:
						},
					},
				}
				app.bot.ReplyMessage(
					&messaging_api.ReplyMessageRequest{
						ReplyToken: e.ReplyToken,
						Messages: []messaging_api.MessageInterface{&messaging_api.FlexMessage{
							Contents: contents,
							AltText:  "Flex message alt text",
						}},
					},
				)
			}
		}
	}
	return
}
