package internal

import (
	"fmt"
	"log"
	"net/http"
)

func Server() {
	http.HandleFunc("/command/latest", HandleLatestCommand)
	http.HandleFunc("/command/ahead", Comments)
	port := ":8080"
	fmt.Printf("Bot is running on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
