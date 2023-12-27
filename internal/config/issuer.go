package config

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

type IssuerConfiger interface {
	IssuerConfig() IssuerConfig
}

func NewIssuerConfiger(getter kv.Getter) IssuerConfiger {
	return &issuer{
		getter: getter,
	}
}

type issuer struct {
	issuerOnce comfig.Once
	getter     kv.Getter
}

type IssuerConfig struct {
	BaseUrl      string `fig:"base_url,required"`
	AuthUsername string `fig:"auth_username,required"`
	AuthPassword string `fig:"auth_password,required"`
}

func (i *issuer) IssuerConfig() IssuerConfig {
	return i.issuerOnce.Do(func() interface{} {
		var result IssuerConfig

		err := figure.
			Out(&result).
			From(kv.MustGetStringMap(i.getter, "issuer")).
			Please()
		if err != nil {
			panic(err)
		}

		return result
	}).(IssuerConfig)
}
