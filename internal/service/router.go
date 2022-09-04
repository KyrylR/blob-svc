package service

import (
	"blob-svc/internal/data/pg"
	"blob-svc/internal/service/handlers"
	"blob-svc/internal/service/helpers"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
)

func (s *service) router() chi.Router {
	r := chi.NewRouter()
	log := s.log.WithFields(map[string]interface{}{
		"service": "blob=svc-api",
	})

	r.Use(
		ape.RecoverMiddleware(log),
		ape.LoganMiddleware(log),
		ape.CtxMiddleware(
			helpers.CtxLog(log),
			helpers.CtxBlobsQ(pg.NewBlobsQ(s.db)),
		),
	)
	r.Route("/integrations/blob-svc", func(r chi.Router) {
		r.Post("/", handlers.CreateBlob)
		r.Get("/", handlers.GetBlobList)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handlers.GetBlob)
			r.Delete("/", handlers.DeleteBlob)
		})
	})

	return r
}
