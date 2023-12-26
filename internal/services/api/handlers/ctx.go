package handlers

import (
	"context"
	"github.com/rarimo/rarime-orgs-svc/internal/data"
	"github.com/rarimo/rarime-orgs-svc/internal/notificator"
	"gitlab.com/distributed_lab/logan/v3"
	"net/http"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	storageCtxKey
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

func CtxNotificator(notificator notificator.Notificator) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, notificatorCtxKey, notificator)
	}
}

func Notificator(r *http.Request) notificator.Notificator {
	return r.Context().Value(notificatorCtxKey).(notificator.Notificator)
}
