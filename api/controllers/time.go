package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"touyakun/models"
)

type TimeController struct {
	timeModel models.TimeModel
}

func InitializeTimeController(t models.TimeModel) *TimeController {
	return &TimeController{timeModel: t}
}

func (tc *TimeController) DeleteTime(c *gin.Context) {
	authKey, exist := c.Get("auth_key")
	if !exist {
		c.Error(errors.New("auth_key not found")).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusUnauthorized, "auth_key not found", "auth_key not found"})
		return
	}

	varidatedAuthKey, ok := authKey.(string)
	if !ok {
		c.Error(errors.New("auth_key is not a string")).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusUnauthorized, "auth_key is not a string", "auth_key is not a string"})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusBadRequest, "id is not a number", "id is not a number"})
		return
	}
	err = tc.timeModel.DeleteTime(varidatedAuthKey, id)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusNotFound, "failed to delete time", "failed to delete time"})
		return
	}

	c.JSON(200, gin.H{"message": "time deleted"})
	return
}
