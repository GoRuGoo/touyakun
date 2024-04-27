package main

import "touyakun/router"

func main() {
	r := router.NewRouter()
	r.Run(":8080")
}
