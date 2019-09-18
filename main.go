package main

import (
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/go-chi/chi"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
)

const (
	ARG_VALUE  = "value"
	COOKIE_KEY = "cookie"
)

type result struct {
	url string
	err error
}

type searchAnswer struct {
	Search string `json:"search"`
	Sites []string `json:"sites"`
}

type searchRequest struct {
	Search string `json:"search"`
	Urls []string `json:"urls"`
}

//выаод информации по контакту c работы
func (searchAnswer *searchAnswer) AddURL(url string) {
	searchAnswer.Sites = append(searchAnswer.Sites, url)
}

func main() {
	router := chi.NewRouter()

	router.Route("/", func(r chi.Router) {
		r.Post("/search", SearchHandler)
		r.Get("/cookie", GetCookieHandler)
		r.Post("/cookie/{value}", SetCookieHandler)
	})

	go func() {
		logrus.Info("server was starts")
		err := http.ListenAndServe(":8080", router)
		log.Fatal(err)
	}()

	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Kill, os.Interrupt)
	<-stopChan

	logrus.Info("shutting down")
}

func GetCookieHandler(w http.ResponseWriter, r *http.Request) {
	cookieValue := "not found"
	cookie, _ := r.Cookie(COOKIE_KEY)
	if cookie == nil {
		respondWithJSON(w, http.StatusNotFound, cookieValue)
		return
	}
	cookieValue = cookie.Value

	respondWithJSON(w, http.StatusOK, cookieValue)
}

func SetCookieHandler(w http.ResponseWriter, r *http.Request) {
	value := chi.URLParam(r, ARG_VALUE)

	cookie := &http.Cookie{
		Name: COOKIE_KEY,
	}
	cookie.Value = value
	http.SetCookie(w, cookie)

	respondWithJSON(w, http.StatusOK, value)
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var sr searchRequest
	err := decoder.Decode(&sr)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	var resultSearchAnswer = searchAnswer{sr.Search, []string{}}
	findResults := findStringInUrls(sr.Search, sr.Urls)
	for _, findResult := range findResults {
		if findResult.url != "" {
			resultSearchAnswer.AddURL(findResult.url)
		}
	}

	respondWithJSON(w, http.StatusOK, resultSearchAnswer)
}

// respondWithJSON write json response format
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func findStringInUrls(search string, urls []string) []result {
	var results []result

	findInUrls := make(chan result)
	for _, url := range urls {
		go findStringInUrl(search, url, findInUrls)
	}

	for i := 0; i < len(urls); i++ {
		result := <- findInUrls
		results = append(results, result)
	}

	return results
}

func findStringInUrl(search string, url string, ch chan <- result) {
	body, err := getBody(url)
	if err != nil {
		ch <- result {"", err}
		return
	}

	if strings.Contains(body, search) {
		ch <- result {url, nil}
	} else {
		ch <- result {"", nil}
	}
}

func getBody(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}