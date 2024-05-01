package router

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"touyakun/models"
	"touyakun/utils"

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
			s := e.Source.(webhook.UserSource)
			switch u.Get("action") {
			case "delete":
				// 薬の一覧を取得
				medications, err := dosageModel.GetMedications(s.UserId)
				if err != nil {
					utils.ReplyTextMessage(app.bot, w, e.ReplyToken, &messaging_api.TextMessage{
						Text: "登録されている薬はありません",
					})
					return
				}
				//ユーザーにどの薬を消すかFlex Messageを使って質問
				contents := []messaging_api.FlexBubble{}
				for _, medication := range medications {
					morningAmount := 0
					afternoonAmount := 0
					eveningAmount := 0
					if medication.IsMorning {
						morningAmount = medication.Amount
					}
					if medication.IsAfternoon {
						afternoonAmount = medication.Amount
					}
					if medication.IsEvening {
						eveningAmount = medication.Amount
					}
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
				utils.ReplyFlexCarouselMessage(app.bot, w, e.ReplyToken, contents)
			case "deleteById":
				medicationId := u.Get("medication_id")
				id, err := strconv.Atoi(medicationId)
				if err != nil {
					w.WriteHeader(500)
					return
				}
				err = dosageModel.DeleteMedications(s.UserId, id)
				if err != nil {
					w.WriteHeader(500)
					return
				}
				utils.ReplyTextMessage(app.bot, w, e.ReplyToken, &messaging_api.TextMessage{
					Text: "削除しました",
				})
			case "settime":
				// どの時間を変えたいか選択する
				template := &messaging_api.ButtonsTemplate{
					Title: "通知時刻変更",
					Text:  "どの時刻を変更する？",
					Actions: []messaging_api.ActionInterface{
						&messaging_api.DatetimePickerAction{
							Label:   "朝の時刻を設定",
							Data:    "action=setMorning",
							Initial: "08:00",
							Mode:    messaging_api.DatetimePickerActionMODE_TIME,
						},
						&messaging_api.DatetimePickerAction{
							Label:   "昼の時刻を設定",
							Data:    "action=setAfternoon",
							Initial: "12:00",
							Mode:    messaging_api.DatetimePickerActionMODE_TIME,
						},
						&messaging_api.DatetimePickerAction{
							Label:   "夜の時刻を設定",
							Data:    "action=setEvening",
							Initial: "20:00",
							Mode:    messaging_api.DatetimePickerActionMODE_TIME,
						},
					},
				}
				utils.ReplyTemplateMessage(app.bot, w, e.ReplyToken, template)
				// case "setMorning":
				// 	// 朝の時間を変更する
				// case: "setAfternoon":
				// 	// 昼の時間を変更する
				// case: "setEvening":
				// 	// 夜の時間を変更する
			}
		}
	}
}
