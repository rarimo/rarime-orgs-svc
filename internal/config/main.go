package config

import (
	"github.com/rarimo/rarime-orgs-svc/internal/data"
	"github.com/rarimo/rarime-orgs-svc/internal/data/pg"
	"github.com/rarimo/rarime-orgs-svc/internal/notificator"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type Config interface {
	comfig.Logger
	comfig.Listenerer
	pgdb.Databaser
	OrgsConfiger
	IssuerConfiger
	notificator.Notificatorer

	Storage() data.Storage
	Orgs() OrgsConfig
	Issuer() IssuerConfig
}

type config struct {
	comfig.Logger
	pgdb.Databaser
	types.Copuser
	comfig.Listenerer
	notificator.Notificatorer
	OrgsConfiger
	IssuerConfiger

	getter kv.Getter
}

func New(getter kv.Getter) Config {
	return &config{
		getter:        getter,
		Listenerer:    comfig.NewListenerer(getter),
		Logger:        comfig.NewLogger(getter, comfig.LoggerOpts{}),
		Databaser:     pgdb.NewDatabaser(getter),
		Notificatorer: notificator.NewNotificatorer(getter),
		getter:         getter,
		Listenerer:     comfig.NewListenerer(getter),
		Logger:         comfig.NewLogger(getter, comfig.LoggerOpts{}),
		OrgsConfiger:   NewOrgsConfiger(getter),
		IssuerConfiger: NewIssuerConfiger(getter),
		Databaser:      pgdb.NewDatabaser(getter),
	}
}

func (c *config) Storage() data.Storage {
	return pg.New(c.DB().Clone())
}

func (c *config) Orgs() OrgsConfig {
	return c.OrgsConfiger.OrgsConfig()
}

func (c *config) Issuer() IssuerConfig {
	return c.IssuerConfiger.IssuerConfig()
}
