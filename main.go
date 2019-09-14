package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/Sirupsen/logrus"
	"github.com/go-chi/chi"
	uuid "github.com/satori/go.uuid"
)

const (
	ARG_KEY    = "key"
	ARG_VALUE  = "value"
	COOKIE_KEY = "cookie"
)

type DataBase map[string]string

var UserDB map[string]DataBase

func main() {
	stopchan := make(chan os.Signal)
	UserDB = map[string]DataBase{}

	logrus.SetReportCaller(true)

	router := chi.NewRouter()

	router.Route("/", func(r chi.Router) {
		r.Get("/{key}", GetIndexHandler)
		r.Post("/{key}", PostIndexHandler)
	})

	go func() {
		logrus.Info("server was starts")
		err := http.ListenAndServe(":8080", router)
		log.Fatal(err)
	}()

	signal.Notify(stopchan, os.Kill, os.Interrupt)
	<-stopchan

	logrus.Info("shutting down")
}

func GetIndexHandler(w http.ResponseWriter, r *http.Request) {
	userKey := CookieControl(w, r)

	key := chi.URLParam(r, ARG_KEY)
	DB, exists := UserDB[userKey]
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	value, exists := DB[key]

	logVal := map[string]interface{}{
		ARG_KEY:   key,
		ARG_VALUE: value,
	}

	if exists {
		logrus.WithFields(logVal).Info("value exists")
		w.Write([]byte(value))
	} else {
		logrus.WithFields(logVal).Info("value does not exists")
		w.WriteHeader(http.StatusNotFound)
	}
}

func PostIndexHandler(w http.ResponseWriter, r *http.Request) {
	userKey := CookieControl(w, r)

	key := chi.URLParam(r, ARG_KEY)
	value := r.FormValue(ARG_VALUE)

	logVal := map[string]interface{}{
		ARG_KEY:   key,
		ARG_VALUE: value,
	}

	UserDB[userKey] = DataBase{}
	UserDB[userKey][key] = value
	logrus.WithFields(logVal).Info("value stored")

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(value))
}

func CookieControl(w http.ResponseWriter, r *http.Request) string {
	cookie, _ := r.Cookie(COOKIE_KEY)
	if cookie == nil {
		cookie = &http.Cookie{
			Name: COOKIE_KEY,
		}
	}

	userKey := cookie.Value

	if userKey != "" {
		return userKey
	} else {
		cookie.Value = uuid.Must(uuid.NewV4()).String()
		http.SetCookie(w, cookie)
		return cookie.Value
	}
}
