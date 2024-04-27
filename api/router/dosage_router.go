package router

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"touyakun/controllers"
	"touyakun/middleware"
	"touyakun/models"
)

func initializeDosageRouter(r *gin.Engine, db *sql.DB) {
	dosage := r.Group("dosage")

	dosageModel := models.InitializeDosageRepo(db)
	dosageController := controllers.InitializeDosageController(dosageModel)

	dosage.GET("/medications", middleware.AuthHandler(), dosageController.GetMedications)
}
