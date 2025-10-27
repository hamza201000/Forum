package main

import (
	"fmt"
	forum "forum/handlers"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/", forum.Handler)
	http.HandleFunc("/register",forum.Handleregister)
	log.Println("Server running on: http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	

}
