package router

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"touyakun/controllers"
	"touyakun/middleware"
	"touyakun/models"
)

func initializeTimeRouter(r *gin.Engine,db *sql.DB){
	time := r.Group("time")

	timeModel := models.InitializeTimeRepo(db)
	timeController := controllers.InitializeTimeController(timeModel)

	time.DELETE("/:id", middleware.AuthHandler(), timeController.DeleteTime)
}