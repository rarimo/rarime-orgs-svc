package handlers

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/rarimo/rarime-orgs-svc/internal/data"
	"github.com/rarimo/rarime-orgs-svc/resources"
	"github.com/rarimo/xo/types/xo"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
	"time"
)

func newGroupCreateRequest(r *http.Request) (*uuid.UUID, *resources.GroupCreateRequest, error) {
	var req resources.GroupCreateRequest

	orgID, err := orgIDFromRequest(r)
	if err != nil {
		return nil, nil, err
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, nil, errors.Wrap(err, "failed to decode body")
	}

	if valid := json.Valid(req.Data.Attributes.Metadata); !valid {
		return nil, nil, validation.Errors{
			"data/attributes/metadata": errors.New("invalid metadata json"),
		}
	}

	if valid := json.Valid(req.Data.Attributes.Rules); !valid {
		return nil, nil, validation.Errors{
			"data/attributes/rules": errors.New("invalid rules json"),
		}
	}

	return &orgID, &req, nil
}

func GroupCreate(w http.ResponseWriter, r *http.Request) {
	orgID, req, err := newGroupCreateRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	org, err := Storage(r).OrganizationQ().OrganizationByIDCtx(r.Context(), *orgID, false)
	if err != nil {
		panic(errors.Wrap(err, "failed to get organization", logan.F{
			"org_id": orgID,
		}))
	}
	if org == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}
	if org.Status != resources.OrganizationStatus_Verified.Int16() {
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"id": errors.Errorf("organization: %s is not verified", org.ID),
		})...)
		return
	}

	// TODO: add auth for role "owner" or "superadmin"

	group := data.Group{
		ID:        uuid.New(),
		OrgID:     *orgID,
		Metadata:  xo.Jsonb(req.Data.Attributes.Metadata),
		Rules:     xo.Jsonb(req.Data.Attributes.Rules),
		CreatedAt: time.Now().UTC(),
	}

	if err := Storage(r).GroupQ().InsertCtx(r.Context(), &group); err != nil {
		panic(errors.Wrap(err, "failed to create group", logan.F{
			"group": group,
		}))
	}

	ape.Render(w, resources.GroupResponse{
		Data:     populateGroup(group),
		Included: resources.Included{},
	})
}
