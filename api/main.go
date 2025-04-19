package main

import (
	"fmt"
	"github.com/ycombinator/usta-match-scheduler/api/routing"
	"net/http"
)

func main() {
	router := routing.NewRouter()

	server := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	fmt.Println("Serving API on localhost:8000...")
	server.ListenAndServe()
}
