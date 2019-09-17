package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	chi "github.com/go-chi/chi"
	middleware "github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
)

const (
	COOKIE_NAME = "name"
)

type WebServer struct {
	srv *http.Server
}

type SearchTask struct {
	SearchString string   `json:"search"`
	Sites        []string `json:"sites"`
}

type ErrorResponse struct {
	Status     string `json:"status"`
	StatusText string `json:"status_text"`
}

func searchText(resp http.ResponseWriter, req *http.Request) {

	jsonData, err := ioutil.ReadAll(req.Body)
	if err != nil {
		createErrorResponse(resp, err.Error(), 500)
		return
	}
	var taskData SearchTask
	err = json.Unmarshal(jsonData, &taskData)
	if err != nil {
		createErrorResponse(resp, err.Error(), 500)
		return
	}

	findedLinks, err := findLinks(taskData.SearchString, taskData.Sites)
	if err != nil {
		createErrorResponse(resp, err.Error(), 500)
		return
	}

	respText, err := json.Marshal(findedLinks)
	if err != nil {
		createErrorResponse(resp, err.Error(), 500)
		return
	}

	resp.Write(respText)

}

func createErrorResponse(resp http.ResponseWriter, statusText string, statusCode int) {
	respJSON, _ := json.Marshal(&ErrorResponse{"error", statusText})
	http.Error(resp, string(respJSON), 500)
}

func findLinks(findText string, findLinks []string) (findedLinks []string, err error) {

	for _, link := range findLinks {

		resp, err := http.Get(link)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if strings.Contains(string(content), findText) {
			findedLinks = append(findedLinks, link)
		}

	}

	return
}

func getCookieHandler(resp http.ResponseWriter, req *http.Request) {

	cookieData, err := req.Cookie(COOKIE_NAME)
	if err != nil {
		createErrorResponse(resp, err.Error(), 500)
		return
	}

	resp.Write([]byte(cookieData.Value))

}

func setCookieHandler(resp http.ResponseWriter, req *http.Request) {

	cookieValue := chi.URLParam(req, "value")

	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: COOKIE_NAME, Value: cookieValue, Path: "/", Expires: expiration}
	http.SetCookie(resp, &cookie)

	resp.Write([]byte("OK"))

}

func CreateWebServer(host string, port string) *WebServer {

	routes := chi.NewRouter()

	routes.Use(middleware.Logger)

	routes.Post("/search", searchText)

	routes.Get("/cookie", getCookieHandler)
	routes.Get("/cookie/{value}", setCookieHandler)

	webServer := &WebServer{}
	webServer.srv = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: routes,
	}

	return webServer
}

func (server *WebServer) Start() error {
	log.Info("Starting server...")
	return server.srv.ListenAndServe()
}

func (server *WebServer) Stop(ctx context.Context) error {
	log.Info("Shutingdown server...")
	return server.srv.Shutdown(ctx)
}

func main() {

	var err error

	newServer := CreateWebServer("127.0.0.1", "8888")

	// if interrupted, shutdown and exit
	signals := make(chan os.Signal, 2)
	signal.Notify(signals, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func() {
		<-signals
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err = newServer.Stop(ctx); err != nil {
			log.Fatal("Failed to stop webServer: " + err.Error())
		}
		os.Exit(0)
	}()

	if err = newServer.Start(); err != nil {
		if err == http.ErrServerClosed {
			log.Info("webServer closed")
		} else {
			log.Fatal("Failed to start webServer: " + err.Error())
		}
	}

}
