package notificator

import (
	"context"
	"net/http"
	"time"
)

var (
	MailChipAPIBaseURL = "https://mandrillapp.com/api/1.0/"
)

type Notificator interface {
	SendVerifyEmail(ctx context.Context, email string, token string) error
}

func newNotificator(appBaseURL, apiKey string) Notificator {
	return &notificator{
		&http.Client{
			Timeout: 1 * time.Minute,
		},
		appBaseURL,
		apiKey,
	}
}

type notificator struct {
	http       *http.Client
	appBaseURL string
	apiKey     string
}
