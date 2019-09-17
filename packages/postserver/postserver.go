// Package postserver implement basic server with json post functionality.
package postserver

/*
	Basics Go.
	Rishat Ishbulatov, dated Sep 16, 2019.

	Using function to search from the past practice, build a server
	that will accept JSON with search request in the POST request and
	return the response as an array of strings in JSON.
*/

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/brie3/HW1/links"
)

// Query stands for search query.
type Query struct {
	Search string   `json:"search"`
	Sites  []string `json:"sites"`
}

// Start starts server with basic post json query handler.
func Start() {
	// stop on ^c
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	rounter := http.NewServeMux()
	rounter.HandleFunc("/", Search)

	// start server
	srv := &http.Server{Addr: ":8000", Handler: rounter}
	go func() {
		log.Println("server starting")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("listen and serve failed: %v", err)
		}
	}()
	<-quit

	// gracefully stop server
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("server stopped")
}

// Search writes json response with links to pages on which the search query was found.
func Search(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("only post method supported.\n"))
		return
	}
	var q Query
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	links, err := links.Find(q.Search, q.Sites)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	jsonlinks, err := json.Marshal(links)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonlinks)
}
