package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

// HandlePostSearch - ищет текст на сайтах
func (serv *Server) HandlePostSearch(w http.ResponseWriter, r *http.Request) {
	inData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	input := SitesSearchReqModel{}
	json.Unmarshal(inData, &input)

	sites, err := SiteSearch(input.Query, input.Sites)
	if err != nil {
		serv.lg.WithError(err).Fatal("can't get sites info")
	}

	output := SitesSearchRespModel{
		Sites: sites,
	}
	outData, err := json.Marshal(output)
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	w.Write(outData)
}

// HandlePostAuth - регистрирует нового пользователя
func (serv *Server) HandlePostAuth(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:  serv.cookieKey,
		Value: uuid.NewV4().String(),
	}
	http.SetCookie(w, &cookie)
	w.Write([]byte(fmt.Sprintf("your new cookie: %s", cookie.Value)))
}

// HandleGetAuth - выводит куку обратившегося пользователя
func (serv *Server) HandleGetAuth(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(serv.cookieKey)
	if err != nil {
		if err == http.ErrNoCookie {
			w.Write([]byte("you haven't our cookies"))
			return
		}
		serv.SendInternalErr(w, err)
		return
	}

	w.Write([]byte(fmt.Sprintf("your cookie key: %s", cookie.Value)))
}

// HandleDeleteAuth - удаляет куки пользователя
func (serv *Server) HandleDeleteAuth(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:    serv.cookieKey,
		Value:   "",
		Expires: time.Unix(0, 0),
	}
	http.SetCookie(w, &cookie)
	w.Write([]byte("your cookies was flushed"))
}
