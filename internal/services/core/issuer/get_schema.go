package issuer

import (
	"github.com/imroc/req/v3"
	"github.com/rarimo/rarime-orgs-svc/resources"
	"net/http"
)

type ClaimSchema struct {
	ActionType string `json:"action_type"`
	SchemaType string `json:"schema_type"`
	SchemaUrl  string `json:"schema_url"`
}

func GetSchemaUrl(r *http.Request, actionType string) (ClaimSchema, error) {
	var result resources.ClaimSchemaResponse

	baseUrl := r.Host

	client := req.C().SetBaseURL("http://" + baseUrl)

	response, err := client.R().SetPathParam("action_type", actionType).Get("/v1/orgs/schema/{action_type}/")
	if err != nil {
		return ClaimSchema{}, err
	}

	if response.Response.StatusCode >= 299 {
		return ClaimSchema{}, err
	}

	err = response.Unmarshal(&result)
	if err != nil {
		return ClaimSchema{}, err
	}

	return ClaimSchema(result.Data.Attributes), nil
}
