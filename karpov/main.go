/*

 */
package main

import (
	"net/http"

	"github.com/art-frela/HW2/karpov/infra"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(infra.FilterContentType)

	r.Post("/search", infra.SearchText)
	http.ListenAndServe(":3333", r)
}
