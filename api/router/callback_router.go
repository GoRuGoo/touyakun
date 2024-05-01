package router

import (
	"database/sql"
	"errors"
	"fmt"
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
		case webhook.PostbackEvent:
			data := e.Postback.Data
			// dataはURLパラメータの形式で書かれているので、パースする
			// 例: "action=buy&itemid=123"
			u, err := url.ParseQuery(data)
			if err != nil {
				w.WriteHeader(500)
				return
			}
			// TODO: medication_idがついていたら削除処理を行う
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
				morningAmount := 0
				afternoonAmount := 0
				eveningAmount := 0
				if medications.morningFlg {
					morningAmount := medications.Amount
				}
				if medications.afternoonFlg {
					afternoonAmount := medications.Amount
				}
				if medications.eveningFlg {
					eveningAmount := medications.Amount
				}
				contents := []messaging_api.FlexBubble{}
				for _, medication := range medications {
					contents = append(contents, messaging_api.FlexBubble{
						Body: &messaging_api.FlexBox{
							Layout: messaging_api.FlexBoxLAYOUT_VERTICAL,
							Contents: []messaging_api.FlexComponentInterface{
								&messaging_api.FlexText{
									Text:   medication.Name,
									Weight: messaging_api.FlexTextWEIGHT_BOLD,
								},
								&messaging_api.FlexText{
									Text: fmt.Sprintf("朝%d錠 昼%d錠 夜%d錠", morningAmount, afternoonAmount, eveningAmount),
								},
								&messaging_api.FlexText{
									Text: fmt.Sprintf("服用期間: %d日分", medication.Duration),
								},
								&messaging_api.FlexButton{
									Action: &messaging_api.PostbackAction{
										Label: "削除",
										Data:  fmt.Sprintf("action=deleteById&medication_id=%d", medication.Id),
									},
								},
							},
						}})
				}
				app.bot.ReplyMessage(
					&messaging_api.ReplyMessageRequest{
						ReplyToken: e.ReplyToken,
						Messages: []messaging_api.MessageInterface{
							&messaging_api.FlexMessage{
								Contents: &messaging_api.FlexCarousel{
									Contents: contents,
								},
								AltText: "Flex message alt text",
							}},
					},
				)
			case "deleteById":
				medicationId := u.Get("medication_id")
			}
		}
	}
	return
}
