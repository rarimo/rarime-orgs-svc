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
		r.Route("/orgs", func(r chi.Router) {
			r.Get("/", handlers.OrgList)
			r.Post("/", handlers.OrgCreate)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", handlers.OrgByID)
				//r.With(handlers.AuthMiddleware()).Get("/users", handlers.OrgUserList)
				r.Get("/users", handlers.OrgUserList)
				//r.Post("/", handlers.OrgVerify)
				//r.Get("/verification-code", handlers.OrgVerificationCode)
			})
			//r.Route("/groups", func(r chi.Router) {
			//r.Get("/", handlers.GroupList)
			//r.With(handlers.AuthMiddleware()).Post("/", handlers.GroupCreate)
			//r.Route("/{group_id}", func(r chi.Router) {
			//	r.Get("/", handlers.GroupByID)
			//	r.Route("/emails", func(r chi.Router) {
			//		r.Group(func(r chi.Router) {
			//			r.Use(handlers.AuthMiddleware())
			//			r.Post("/", handlers.InvitationEmailCreate)
			//			r.Patch("/", handlers.InvitationEmailAccept)
			//		})
			//	})
			//	r.Route("/requests", func(r chi.Router) {
			//		r.Group(func(r chi.Router) {
			//			r.Use(handlers.AuthMiddleware())
			//			r.Get("/", handlers.RequestList)
			//			r.Route("/{request_id}", func(r chi.Router) {
			//				r.Patch("/", handlers.RequestFill)
			//				r.Post("/", handlers.RequestVerify)
			//			})
			//		})
			//	})
			//})
			//})
		})
	})

	cfg.Log().WithFields(logan.F{
		"service": "api",
		"addr":    cfg.Listener().Addr(),
	}).Info("starting api")

	ape.Serve(ctx, r, cfg, ape.ServeOpts{})
}
