package issuer

import (
	"encoding/json"
	"io/ioutil"

	"github.com/rarimo/rarime-orgs-svc/internal/config"

	"github.com/imroc/req/v3"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/logan/v3"
)

type Issuer interface {
	IssueClaim(
		userDid string,
		credentialSubject interface{},
	) (*IssueClaimResponse, error)
	SchemaType() string
	SchemaURL() string
}

type issuer struct {
	client       *req.Client
	authUsername string
	authPassword string
	schemaType   string
	schemaURL    string
}

func New(log *logan.Entry, config *config.IssuerConfig) Issuer {
	return &issuer{
		client: req.C().
			SetBaseURL(config.BaseUrl).
			SetLogger(log),
		schemaType:   config.SchemaType,
		schemaURL:    config.SchemaUrl,
		authUsername: config.AuthUsername,
		authPassword: config.AuthPassword,
	}
}

func (is *issuer) IssueClaim(
	userDid string,
	credentialsSubject interface{},
) (*IssueClaimResponse, error) {
	var result UUIDResponse

	response, err := is.client.R().
		SetBasicAuth(is.authUsername, is.authPassword).
		SetBodyJsonMarshal(credentialsSubject).
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
	body, err := ioutil.ReadAll(response.Body)
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

func (is *issuer) SchemaType() string {
	return is.schemaType
}

func (is *issuer) SchemaURL() string {
	return is.schemaURL
}
