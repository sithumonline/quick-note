package router

import (
	"github.com/go-chi/chi"

	"github.com/sithumonline/quick-note/api/handler"
)

func HealthRoute() chi.Router {
	r := chi.NewRouter()

	r.Get("/", handler.GetHealth)

	return r
}
