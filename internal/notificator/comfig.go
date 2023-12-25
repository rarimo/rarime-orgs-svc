package notificator

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type Notificatorer interface {
	Notificator() Notificator
}

func NewNotificatorer(getter kv.Getter) Notificatorer {
	return &notificatorer{
		getter: getter,
	}
}

type notificatorer struct {
	getter kv.Getter
	comfig.Once
}

func (b *notificatorer) Notificator() Notificator {
	return b.Do(func() interface{} {
		var config struct {
			AppBaseURL string `fig:"app_base_url,required"`
			ApiKey     string `fig:"api_key,required"`
		}
		err := figure.
			Out(&config).
			From(kv.MustGetStringMap(b.getter, "notificator")).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out notificator"))
		}
		return newNotificator(config.AppBaseURL, config.ApiKey)
	}).(Notificator)
}
