package router

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
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
	timeModel := models.InitializeTimeRepo(app.db)

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
				times, err := timeModel.GetMedicationRemindTimeList(s.UserId)
				if err != nil {
					w.WriteHeader(500)
					return
				}
				// どの時間を変えたいか選択する
				template := &messaging_api.ButtonsTemplate{
					Title: "通知時刻変更",
					Text:  "どの時刻を変更する？ 変更したいところをタップ！",
					Actions: []messaging_api.ActionInterface{
						&messaging_api.DatetimePickerAction{
							Label:   "朝：(現在の設定:" + times.MorningTime + ")",
							Data:    "action=registerTime&time=morning",
							Initial: times.MorningTime,
							Mode:    messaging_api.DatetimePickerActionMODE_TIME,
						},
						&messaging_api.DatetimePickerAction{
							Label:   "昼：(現在の設定:" + times.AfternoonTime + ")",
							Data:    "action=registerTime&time=afternoon",
							Initial: times.AfternoonTime,
							Mode:    messaging_api.DatetimePickerActionMODE_TIME,
						},
						&messaging_api.DatetimePickerAction{
							Label:   "夜：(現在の設定:" + times.EveningTime + ")",
							Data:    "action=registerTime&time=evening",
							Initial: times.EveningTime,
							Mode:    messaging_api.DatetimePickerActionMODE_TIME,
						},
					},
				}
				utils.ReplyTemplateMessage(app.bot, w, e.ReplyToken, template)
			case "registerTime":
				timeParam, found := e.Postback.Params["time"]
				if !found {
					w.WriteHeader(500)
					return
				}
				parsedTime, err := time.Parse("15:04", timeParam)
				if err != nil {
					w.WriteHeader(500)
					return
				}
				timeText := ""
				switch u.Get("time") {
				case "morning":
					err = timeModel.RegisterMorningTime(s.UserId, parsedTime)
					timeText = "朝"
				case "afternoon":
					err = timeModel.RegisterAfternoonTime(s.UserId, parsedTime)
					timeText = "昼"
				case "evening":
					err = timeModel.RegisterEveningTime(s.UserId, parsedTime)
					timeText = "夜"
				}
				if err != nil {
					w.WriteHeader(500)
					return
				}
				utils.ReplyTextMessage(app.bot, w, e.ReplyToken, &messaging_api.TextMessage{
					Text: timeText + "の時刻を" + parsedTime.Format("15:04") + "に設定しました！",
				})
			case "register":
				//薬の画像を送らせる
				utils.ReplyTextMessage(app.bot, w, e.ReplyToken, &messaging_api.TextMessage{
					Text: "薬の画像を送ってください",
					QuickReply: &messaging_api.QuickReply{
						Items: []messaging_api.QuickReplyItem{
							{
								Action: messaging_api.CameraAction{
									Label: "カメラで撮影",
								},
							},
							{
								Action: messaging_api.CameraRollAction{
									Label: "カメラロールから選択",
								},
							},
						}}})
				// 将来この処理のあとにのみ画像を受け取るようにするなら、userモデルなどに画像を受け取る状態を持たせる

			}

		case webhook.MessageEvent:
			switch message := e.Message.(type) {
			case webhook.ImageMessageContent:
				// 薬の情報をAPIから取得
				content, err := app.blob.GetMessageContentTranscodingByMessageId(message.Id)
				if err != nil {
					w.WriteHeader(500)
					return
				}
				fmt.Println(content)
				// app.bot.ShowLoadingAnimation(&messaging_api.ShowLoadingAnimationRequest{
				// 	ChatId:         s.UserId,
				// 	LoadingSeconds: 10,
				// })
			}
		}
	}
}
