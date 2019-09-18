package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()

	r.Post("/search", func(w http.ResponseWriter, r *http.Request) {
		type searchStruct struct {
			Search string   `json:"search,omitempty"`
			Urls   []string `json:"sites,omitempty"`
		}
		var (
			req     searchStruct
			matched []string
		)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)
		if err != nil {
			w.Write([]byte(err.Error()))
		}

		matched = searchInArray(req.Search, req.Urls)

		data, err := json.Marshal(matched)
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		w.Write(data)
	})

	log.Fatal(http.ListenAndServe(":8080", r))

}

func searchInArray(search string, urls []string) []string {
	var matchedUrls []string

	search = strings.ToLower(search)

	for _, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			log.Panic(err)
		}
		cont, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Panic(err)
		}
		if strings.Contains(strings.ToLower(string(cont)), search) {
			matchedUrls = append(matchedUrls, url)
		}
	}

	return matchedUrls
}
