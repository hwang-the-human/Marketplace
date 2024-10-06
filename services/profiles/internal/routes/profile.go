package routes

import (
	"github.com/go-chi/chi"
	"marketplace/services/profiles/internal/handlers"
)

func ProfileRouter(profileHandler *handlers.ProfileHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/profiles", func(r chi.Router) {
		r.Get("/{id}", profileHandler.GetProfile)
		r.Post("/", profileHandler.CreateProfile)
		r.Put("/{id}", profileHandler.UpdateProfile)
	})

	return r
}
