package router

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"touyakun/controllers"
	"touyakun/models"
)

func initializeDosageRouter(r *gin.Engine, db *sql.DB) {
	dosage := r.Group("/dosage")

	dosageModel := models.InitializeDosageRepo(db)
	dosageController := controllers.InitializeDosageController(dosageModel)

	dosage.GET("/medications", dosageController.GetMedications)
}
