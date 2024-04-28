package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"touyakun/models"
)

type DosageController struct {
	dosageModel models.DosageModel
}

func InitializeDosageController(d models.DosageModel) *DosageController {
	return &DosageController{dosageModel: d}
}

func (dc *DosageController) GetMedications(c *gin.Context) {
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
	fmt.Println(varidatedAuthKey)
	medications, err := dc.dosageModel.GetMedications("hoge")
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusNotFound, "failed to get medications", "failed to get medications"})
		return
	}

	c.JSON(200, gin.H{"medication_list": medications})
	return
}
