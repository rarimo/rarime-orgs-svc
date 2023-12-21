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

func newOrgCreateRequest(r *http.Request) (*resources.OrganizationCreate, error) {
	var req resources.OrganizationCreate

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(err, "failed to decode body")
	}

	if valid := json.Valid(req.Attributes.Metadata); !valid {
		return nil, validation.Errors{
			"data/attributes/metadata": errors.New("invalid metadata json"),
		}
	}

	return &req, validation.Errors{
		"data/attributes/domain":    validation.Validate(req.Attributes.Domain, validation.Required),
		"data/attributes/owner_did": validation.Validate(req.Attributes.OwnerDid, validation.Required),
	}
}

func OrgCreate(w http.ResponseWriter, r *http.Request) {
	req, err := newOrgCreateRequest(r)

	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	owner := data.User{
		ID:        uuid.New(),
		Did:       req.Attributes.OwnerDid,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	org := data.Organization{
		ID:           uuid.New(),
		Owner:        owner.ID,
		MembersCount: 1,
		Domain:       req.Attributes.Domain,
		Metadata:     xo.Jsonb(req.Attributes.Metadata),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	err = Storage(r).Transaction(func() error {
		err = Storage(r).UserQ().InsertCtx(r.Context(), &owner)
		if err != nil {
			return errors.Wrap(err, "failed to insert new user", logan.F{
				"did": req.Attributes.OwnerDid,
			})
		}
		err = Storage(r).OrganizationQ().InsertCtx(r.Context(), &org)
		if err != nil {
			return errors.Wrap(err, "failed to insert new org", logan.F{
				"owner_did": req.Attributes.OwnerDid,
				"domain":    req.Attributes.Domain,
			})
		}

		return nil
	})
	if err != nil {
		Log(r).WithError(err).Error("failed to create new organization")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	resp := resources.OrganizationResponse{
		Data:     populateOrg(org),
		Included: resources.Included{},
	}

	respUser := populateUser(owner)
	resp.Included.Add(&respUser)

	ape.Render(w, resp)
}

func populateUser(user data.User) resources.User {
	res := resources.User{
		Key: resources.Key{
			ID:   user.ID.String(),
			Type: resources.USERS,
		},
		Attributes: resources.UserAttributes{
			Role:      resources.UserRole(user.Role),
			OrgId:     user.OrgID.String(),
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		Relationships: &resources.UserRelationships{
			Organization: resources.Relation{
				Data: &resources.Key{
					ID:   user.OrgID.String(),
					Type: resources.ORGANIZATIONS,
				},
			},
		},
	}

	if user.Did != "" {
		res.Attributes.Did = &user.Did
	}

	return res
}
