package router

import (
	"database/sql"
	"log"
	"touyakun/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	r.Use(cors.New(config))
	r.Use(middleware.ErrorHandler())

	db, err := sql.Open("postgres", "host=db port=5432 user=touyakun password=password dbname=touyakun sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	initializeDosageRouter(r, db)

	initializeTimeRouter(r, db)

	return r

}
