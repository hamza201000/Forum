package main

import (
	"fmt"
	"log"
	"net/http"

	forum "forum/handlers"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", forum.Handler)
	mux.HandleFunc("/register", forum.HandleRegister)
	mux.HandleFunc("/register/data", forum.HandlerDataRegister)
	mux.HandleFunc("/login", forum.HandleLogin)
	mux.HandleFunc("/post", forum.HandlePost)
	mux.HandleFunc("/addpost", forum.HandleAddPost)
	log.Println("Server running on: http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}
