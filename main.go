package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Benson-14/log-aggregator/parser"
	"github.com/Benson-14/log-aggregator/query"
	"github.com/Benson-14/log-aggregator/storage"
)

func main() {
	store := storage.NewMemoryStorage()
	executor := query.NewExecutor(store)

	a := &app{
		store:    store,
		executor: executor,
		parser:   parser.NewJSONParser(),
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /logs", a.handleIngest)
	mux.HandleFunc("GET /logs", a.handleGetLogs)
	mux.HandleFunc("GET /logs/search", a.handleSearch)
	mux.HandleFunc("GET /logs/stats", a.handleStats)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Server running on http://localhost:8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
