package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("test")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json") // JSONを返すことを宣言
		response := map[string]string{"message": "Hello, World!!!!"}
		json.NewEncoder(w).Encode(response) // レスポンスをJSON形式にエンコードして書き込み
	})

	http.ListenAndServe(":8080", nil)
}

