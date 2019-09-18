package main

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/go-chi/chi"
)

const (
	COOKIE_NAME  = "cookie"
	COOKIE_VALUE = "geekbrains"
)

func main() {
	router := chi.NewRouter()

	router.Get("/cookie", sendCookieHandler)
	router.Post("/cookie", getCookieHandler)

	logrus.Fatal(http.ListenAndServe(":8080", router))
}

func sendCookieHandler(w http.ResponseWriter, r *http.Request) {
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := &http.Cookie{
		Name:    COOKIE_NAME,
		Value:   COOKIE_VALUE,
		Expires: expiration,
		Path:    "/"}
	http.SetCookie(w, cookie)
	logrus.Info("cookie set:" + cookie.String())
	w.Write([]byte("cookie set:" + cookie.String()))
}

func getCookieHandler(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(COOKIE_NAME); err != nil {
		//TODO
		logrus.Info(err)
		return
	} else {
		logrus.Info("got cookie:" + cookie.String())
		w.Write([]byte("got cookie: " + cookie.String()))
	}
}
