package config

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

type OrgsConfiger interface {
	OrgsConfig() OrgsConfig
}

func NewOrgsConfiger(getter kv.Getter) OrgsConfiger {
	return &Orgs{
		getter: getter,
	}
}

type Orgs struct {
	OrgsOnce comfig.Once
	getter   kv.Getter
}

type OrgsConfig struct {
	VerifyDomain bool `fig:"verify_domain,required"`
}

func (e *Orgs) OrgsConfig() OrgsConfig {
	return e.OrgsOnce.Do(func() interface{} {
		var result OrgsConfig

		err := figure.
			Out(&result).
			From(kv.MustGetStringMap(e.getter, "orgs")).
			Please()
		if err != nil {
			panic(err)
		}

		return result
	}).(OrgsConfig)
}
