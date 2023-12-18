package api

import (
	"context"
	"gitlab.com/distributed_lab/logan/v3"
	"time"

	"github.com/go-chi/chi"
	"github.com/rarimo/rarime-orgs-svc/internal/config"
	"github.com/rarimo/rarime-orgs-svc/internal/services/api/handlers"
	"gitlab.com/distributed_lab/ape"
)

func Run(ctx context.Context, cfg config.Config) {
	r := chi.NewRouter()

	const slowRequestDurationThreshold = time.Second
	ape.DefaultMiddlewares(r, cfg.Log(), slowRequestDurationThreshold)

	r.Use(
		ape.CtxMiddleware(
			handlers.CtxLog(cfg.Log()),
			handlers.CtxStorage(cfg.Storage()),
		),
	)

	r.Route("/v1", func(r chi.Router) {
		r.Route("/proofs", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(handlers.AuthMiddleware())
				r.Post("/", handlers.ProofCreate)
			})
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", handlers.ProofByID)
			})
		})
	})

	cfg.Log().WithFields(logan.F{
		"service": "api",
		"addr":    cfg.Listener().Addr(),
	}).Info("starting api")

	ape.Serve(ctx, r, cfg, ape.ServeOpts{})
}
