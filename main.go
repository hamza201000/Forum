package main

import (
	"fmt"
	"log"
	"net/http"

	forum "forum/handlers"
)

func main() {
	http.HandleFunc("/", forum.Handler)
	http.HandleFunc("/register", forum.HandleRegister)
	log.Println("Server running on: http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}
