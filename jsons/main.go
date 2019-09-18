/*
 * HomeWork-2: Search string - 2: JSON server
 * Created on 15.09.19 12:12
 * Copyright (c) 2019 - Eugene Klimov
 */

package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

const (
	servAddr  = "localhost:8080"
	sitesFile = "sites.txt"
)

type google struct {
	Search        string   `json:"search"`
	Sites         []string `json:"sites"`
	CaseSensitive bool     `json:"case_sens"`
	urls          []string
	mux           sync.Mutex
}

func (g *google) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, "Hello and GoodBye! - Need POST method.\n")
		return
	}

	// decode POST data
	err := json.NewDecoder(r.Body).Decode(g)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Can't parse POST data:", err)
		return
	}

	// get search results
	err = g.searchStringURL()
	if err != nil {
		log.Println("Error while search:", err) // no need return bec. it may be for one site
	}

	// encode to json
	b, err := json.MarshalIndent(g, "", "    ") // for best view in curl
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Can't encode result data:", err)
		return
	}

	// set proper headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	w.Write(b)
}

func (g *google) searchStringURL() error {

	var eg errgroup.Group
	g.Sites = []string{} // reset results

	for _, url := range g.urls {

		if len(url) < 3 { // no fake strings
			continue
		}

		url := url
		eg.Go(func() error {

			resp, err := http.Get(url)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			s := string(body)
			if !g.CaseSensitive {
				s = strings.ToLower(s)
				g.Search = strings.ToLower(g.Search)
			}

			if strings.Contains(s, g.Search) {
				g.mux.Lock()
				g.Sites = append(g.Sites, url)
				g.mux.Unlock()
			}
			return nil
		})
	}

	return eg.Wait()
}

func (g *google) getURLs() error {

	file, err := os.Open(sitesFile)
	if err != nil {
		return err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	g.urls = strings.Split(string(b), "\n")

	return nil
}

func main() {

	googleHandler := &google{}

	// get URLs from file
	err := googleHandler.getURLs()
	if err != nil {
		log.Fatalln("Error reading sites from file:", err)
	}

	// safe shutdown
	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM) // os.Kill wrong for linters

	go func() {
		<-shutdown
		// any safe work here
		fmt.Printf("\nShutdown server at: %s\n", servAddr)
		os.Exit(0)
	}()

	// prepare server, no need smart router for simple scenario
	http.Handle("/", googleHandler)

	fmt.Println("Starting server at:", servAddr)
	log.Fatalln(http.ListenAndServe(servAddr, nil))
}

// curl --header "Content-Type: application/json" --request POST --data '{"search":"bug"}' http://localhost:8080
// curl --header "Content-Type: application/json" --request POST --data "{\"search\":\"bug\"}" http://localhost:8080
// "Бим", "Книга", "книга", "1973", "2033"
