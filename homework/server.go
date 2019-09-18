package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

// Server - Объект сервера
type Server struct {
	lg        *logrus.Logger
	cookieKey string
}

// NewServer - создаёт новый экземпляр сервера
func NewServer(lg *logrus.Logger) *Server {
	return &Server{
		lg:        lg,
		cookieKey: "HW2_COOKIE",
	}
}

// Start - запускает сервер
func (serv *Server) Start() error {
	r := chi.NewRouter()
	r.Use(serv.RequestTracerMiddleware)
	serv.ConfigureHandlers(r)
	serv.lg.Info("server is running")
	return http.ListenAndServe(":8080", r)
}

// ConfigureHandlers - настраивает хендлеры и их пути
func (serv *Server) ConfigureHandlers(r *chi.Mux) {
	r.Route("/", func(r chi.Router) {
		r.Post("/search", serv.HandlePostSearch)
		r.Route("/auth", func(r chi.Router) {
			r.Post("/", serv.HandlePostAuth)
			r.Get("/", serv.HandleGetAuth)
			r.Delete("/", serv.HandleDeleteAuth)
		})
	})
}

// SendErr - отправляет ошибку пользователю и логирует её
func (serv *Server) SendErr(w http.ResponseWriter, err error, code int, obj ...interface{}) {
	serv.lg.WithField("data", obj).WithError(err).Error("server error")
	w.WriteHeader(code)
	errModel := ErrorModel{
		Code:     code,
		Err:      err.Error(),
		Desc:     "server error",
		Internal: obj,
	}
	data, _ := json.Marshal(errModel)
	w.Write(data)
}

// SendInternalErr - отправляет 500 ошибку
func (serv *Server) SendInternalErr(w http.ResponseWriter, err error, obj ...interface{}) {
	serv.SendErr(w, err, http.StatusInternalServerError, obj)
}

// RequestTracerMiddleware - отслеживает и логирует входящие запросы
func (serv *Server) RequestTracerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		serv.lg.
			WithFields(map[string]interface{}{
				"url":    r.URL.String(),
				"cookie": r.Header.Get("Cookie"),
				"body":   string(body),
			}).
			Debug("request")
		next.ServeHTTP(w, r)
	})
}
