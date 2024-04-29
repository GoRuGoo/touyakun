package main

import (
	"log"
	"net/http"
	"touyakun/router"
)

func main() {
	router.InitializeRouter()
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
