package notificator

import (
	"bytes"
	"context"
	"encoding/json"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"io"
	"net/http"
	"net/url"
	"path"
)

type sendVerifyEmailBody struct {
	Key             string        `json:"key"`
	TemplateName    string        `json:"template_name"`
	TemplateContent []nameContent `json:"template_content"`
	Async           bool          `json:"async"`
	Message         message       `json:"message"`
}

type message struct {
	FromEmail       string     `json:"from_email"`
	FromName        string     `json:"from_name"`
	Subject         string     `json:"subject"`
	To              []to       `json:"to"`
	Important       bool       `json:"important"`
	ViewContentLink bool       `json:"view_content_link"`
	MergeVars       []mergeVar `json:"merge_vars"`
}

type mergeVar struct {
	Rcpt string        `json:"rcpt"`
	Vars []nameContent `json:"vars"`
}

type nameContent struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type to struct {
	Email string `json:"email"`
}

func (n notificator) SendVerifyEmail(ctx context.Context, email, token string) error {
	redirectURL, err := url.Parse(n.appBaseURL)
	if err != nil {
		return errors.Wrap(err, "failed to parse redirect url", logan.F{
			"app_base_url": n.appBaseURL,
		})
	}
	redirectURL.Path = path.Join(redirectURL.Path, "/verify-email")
	q := redirectURL.Query()
	q.Add("q", token)
	redirectURL.RawQuery = q.Encode()

	data := sendVerifyEmailBody{
		Key:             n.apiKey,
		TemplateName:    "rarime-verify-template",
		TemplateContent: make([]nameContent, 0),
		Message: message{
			FromEmail: "info@rarimo.com",
			FromName:  "Rarimo",
			Subject:   "Please verify your email address",
			To:        []to{{Email: email}},
			Important: true,
			MergeVars: []mergeVar{{
				Rcpt: email,
				Vars: []nameContent{
					{
						Name:    "RARIME_VERIFY_URL",
						Content: redirectURL.String(),
					},
					{
						Name:    "RECIPIENT_EMAIL_ADDRESS",
						Content: email,
					},
				},
			}},
		},
	}
	body, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal send verify email body", logan.F{
			"email": email,
		})
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		MailChipAPIBaseURL,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create send verify email request", logan.F{
			"email": email,
		})
	}

	req.URL.Path = path.Join(req.URL.Path, "/messages/send-template")
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.http.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to send verify email", logan.F{
			"email": email,
		})
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.From(errors.New("failed to send verify email"), logan.F{
			"email":       email,
			"status_code": resp.StatusCode,
		})
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read response body", logan.F{
			"email": email,
		})
	}

	var response []struct {
		Status       string  `json:"status"`
		RejectReason *string `json:"reject_reason"`
	}

	if err = json.Unmarshal(respBody, &response); err != nil {
		return errors.Wrap(err, "failed to unmarshal response body", logan.F{
			"email": email,
		})
	}
	if len(response) == 0 {
		return errors.From(errors.New("no response"), logan.F{
			"email": email,
		})
	}
	if response[0].Status != "sent" {
		return errors.From(errors.New("failed to send verify email"), logan.F{
			"email":         email,
			"status":        response[0].Status,
			"reject_reason": response[0].RejectReason,
		})
	}

	return nil
}
