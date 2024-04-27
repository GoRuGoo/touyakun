package router

import (
	"database/sql"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	r.Use(cors.New(config))

	db, err := sql.Open("postgres", "host=localhost port=5432 user=test password=password dbname=test sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	initializeDosageRouter(r, db)

	return r
}
