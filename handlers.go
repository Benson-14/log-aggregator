package main

import (
	"encoding/json"
	"net/http"

	"github.com/Benson-14/log-aggregator/parser"
	"github.com/Benson-14/log-aggregator/query"
	"github.com/Benson-14/log-aggregator/storage"
)

type app struct {
	store    storage.Storage
	parser   parser.Parser
	executor *query.Executor
}

func (a *app) handleIngest(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Line string `json:"line"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	entry, err := a.parser.Parse(body.Line)
	if err != nil {
		http.Error(w, "failed to parse log line: "+err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if err := a.store.Append(entry); err != nil {
		http.Error(w, "failed to store entry: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (a *app) handleSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		http.Error(w, "missing query parameter: q", http.StatusBadRequest)
		return
	}

	results := a.executor.Run(q)
	writeJSON(w, http.StatusOK, results)
}

func (a *app) handleGetLogs(w http.ResponseWriter, r *http.Request) {
	entries := a.store.All()
	writeJSON(w, http.StatusOK, entries)
}

func (a *app) handleStats(w http.ResponseWriter, r *http.Request) {
	stats := map[string]any{
		"total":     a.store.Len(),
		"by_level":  a.store.CountByLevel(),
		"by_source": a.store.CountBySource(),
	}
	writeJSON(w, http.StatusOK, stats)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
