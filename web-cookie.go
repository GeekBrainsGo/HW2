package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/Sirupsen/logrus"
)

const (
	COOKIE_KEY = "IPaddress"
)

func main() {
	stopchan := make(chan os.Signal)

	// logrus.SetReportCaller(true)

	router := http.NewServeMux()

	router.HandleFunc("/", setCookie)
	router.HandleFunc("/getCookie", getCookie)

	go func() {
		logrus.Info("Сервер запущен")
		err := http.ListenAndServe(":8080", router)
		logrus.Fatal(err)
	}()

	signal.Notify(stopchan, os.Kill, os.Interrupt)

	<-stopchan

	logrus.Info("Сервер остановлен!")

}

func setCookie(wr http.ResponseWriter, req *http.Request) {
	cookie := &http.Cookie{
		Name: COOKIE_KEY,
	}
	cookie.Value = req.RemoteAddr
	http.SetCookie(wr, cookie)
	fmt.Fprintf(wr, "Установили cookie: %s", []byte(cookie.Value))
}

func getCookie(wr http.ResponseWriter, req *http.Request) {
	cookie, _ := req.Cookie(COOKIE_KEY)
	if cookie == nil {
		http.Error(wr, "Отсутствует cookie", http.StatusForbidden)
		logrus.Info("Отсутствует cookie")
		return
	}
	userkey := cookie.Value
	fmt.Fprintf(wr, "Сохраненный cookie: %s", []byte(userkey))
}
