package controllers

import (
	"github.com/gin-gonic/gin"
	"touyakun/models"
)

type DosageController struct {
	dosageModel models.DosageModel
}

func InitializeDosageController(d models.DosageModel) *DosageController {
	return &DosageController{dosageModel: d}
}

func (dc *DosageController) GetMedications(c *gin.Context) {
	test := dc.dosageModel.GetMedications()
	c.JSON(200, gin.H{"test": test})
	return
}
