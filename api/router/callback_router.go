package router

import (
	"database/sql"
	"errors"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
	"net/http"
	"touyakun/models"
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

	for _, event := range cb.Events {
		switch e := event.(type) {
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
		}
	}
	return
}
