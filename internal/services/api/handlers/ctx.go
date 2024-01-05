package handlers

import (
	"context"
	"net/http"

	"github.com/rarimo/rarime-orgs-svc/internal/config"
	"github.com/rarimo/rarime-orgs-svc/internal/data"
	"github.com/rarimo/rarime-orgs-svc/internal/notificator"
	"github.com/rarimo/rarime-orgs-svc/internal/services/core/issuer"
	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	storageCtxKey
	orgsConfigCtxKey
	issuerCtxKey
	notificatorCtxKey
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func CtxStorage(storage data.Storage) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, storageCtxKey, storage)
	}
}

func Storage(r *http.Request) data.Storage {
	return r.Context().Value(storageCtxKey).(data.Storage)
}

func CtxOrgsConfig(cfg config.OrgsConfig) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, orgsConfigCtxKey, cfg)
	}
}

func OrgsConfig(r *http.Request) config.OrgsConfig {
	return r.Context().Value(orgsConfigCtxKey).(config.OrgsConfig)
}

func CtxIssuer(iss *issuer.Issuer) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, issuerCtxKey, iss)
	}
}

func Issuer(r *http.Request) *issuer.Issuer {
	return r.Context().Value(issuerCtxKey).(*issuer.Issuer)
}

func CtxNotificator(notificator notificator.Notificator) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, notificatorCtxKey, notificator)
	}
}

func Notificator(r *http.Request) notificator.Notificator {
	return r.Context().Value(notificatorCtxKey).(notificator.Notificator)
}
