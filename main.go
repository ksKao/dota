package main

import (
	"log"
	"net/http"

	"github.com/ksKao/dota/handlers"
)

func main() {
	handlers.RegisterRootRoutes()

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
