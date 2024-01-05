package issuer

import (
	"encoding/json"
	"io"

	"github.com/rarimo/rarime-orgs-svc/internal/config"
	"github.com/rarimo/rarime-orgs-svc/internal/services/core/issuer/schemas"

	"github.com/imroc/req/v3"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/logan/v3"
)

type Issuer struct {
	client       *req.Client
	authUsername string
	authPassword string
}

func New(log *logan.Entry, config config.IssuerConfig) *Issuer {
	return &Issuer{
		client: req.C().
			SetBaseURL(config.BaseUrl).
			SetLogger(log),
		authUsername: config.AuthUsername,
		authPassword: config.AuthPassword,
	}
}

func (is *Issuer) IssueClaim(userDid string, req schemas.CreateCredentialRequest) (*IssueClaimResponse, error) {
	var result UUIDResponse

	response, err := is.client.R().
		SetBasicAuth(is.authUsername, is.authPassword).
		SetBodyJsonMarshal(req).
		SetSuccessResult(result).
		SetPathParam("identifier", userDid).
		Post("/{identifier}/claims")
	if err != nil {
		return nil, errors.Wrap(err, "failed to send post request")
	}

	if response.StatusCode >= 299 {
		return nil, errors.Wrapf(ErrUnexpectedStatusCode, response.String())
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal json")
	}

	resp := result.IssueClaimResponse()
	return &resp, nil
}
