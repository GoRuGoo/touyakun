package controllers

import (
	"net/http"
	"touyakun/models"
)

type TimeController struct {
	timeModel models.TimeModel
	w         http.ResponseWriter
}

func InitializeTimeController(t models.TimeModel, w http.ResponseWriter) *TimeController {
	return &TimeController{timeModel: t, w: w}
}

func (tc *TimeController) DeleteTime(userId string, timeId int) {
	err := tc.timeModel.DeleteTime(userId, timeId)
	if err != nil {
		tc.w.WriteHeader(500)
		tc.w.Write([]byte(err.Error()))
		return
	}

	tc.w.WriteHeader(200)
	return
}
