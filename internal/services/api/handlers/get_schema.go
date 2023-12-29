package handlers

import (
	"github.com/go-chi/chi"
	"github.com/rarimo/rarime-orgs-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/urlval"
	"net/http"
)

type schemaByActionTypeRequest struct {
	ActionType string `json:"action_type"`
}

func newSchemaByActionTypeRequest(r *http.Request) (*schemaByActionTypeRequest, error) {
	var request schemaByActionTypeRequest

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode url")
	}

	request.ActionType = chi.URLParam(r, "action_type")
	if request.ActionType == "" {
		return nil, errors.New("action_type is required")
	}

	return &request, nil
}

func SchemaByActionType(w http.ResponseWriter, r *http.Request) {
	req, err := newSchemaByActionTypeRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	schema, err := Storage(r).ClaimSchemaQ().SchemaByActionTypeCtx(r.Context(), req.ActionType)

	if err != nil {
		panic(errors.Wrap(err, "failed to get schema by action type", logan.F{
			"action_type": req.ActionType,
		}))
	}
	if schema == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	resp := resources.ClaimSchemaResponse{
		Data: resources.ClaimSchema{
			Key: resources.Key{},
			Attributes: resources.ClaimSchemaAttributes{
				ActionType: req.ActionType,
				SchemaType: schema.SchemaType,
				SchemaUrl:  schema.SchemaUrl,
			},
		},
	}

	ape.Render(w, resp)
}
