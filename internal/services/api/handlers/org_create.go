package handlers

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	rules "github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"github.com/rarimo/rarime-orgs-svc/internal/data"
	"github.com/rarimo/rarime-orgs-svc/resources"
	"github.com/rarimo/xo/types/xo"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
	"strings"
	"time"
)

func newOrgCreateRequest(r *http.Request) (*resources.OrganizationCreateRequest, error) {
	var req resources.OrganizationCreateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(err, "failed to decode body")
	}

	if valid := json.Valid(req.Data.Attributes.Metadata); !valid {
		return nil, validation.Errors{
			"data/attributes/metadata": errors.New("invalid metadata json"),
		}
	}

	return &req, validation.Errors{
		"data/attributes/domain":    validation.Validate(req.Data.Attributes.Domain, validation.Required, rules.URL),
		"data/attributes/owner_did": validation.Validate(req.Data.Attributes.OwnerDid, validation.Required, ValidationDID),
	}.Filter()
}

func OrgCreate(w http.ResponseWriter, r *http.Request) {
	req, err := newOrgCreateRequest(r)

	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	req.Data.Attributes.Domain = removeScheme(req.Data.Attributes.Domain)

	orgID := uuid.New()

	owner := data.User{
		ID:        uuid.New(),
		OrgID:     orgID,
		Did:       req.Data.Attributes.OwnerDid,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	org := data.Organization{
		ID:           orgID,
		Owner:        owner.ID,
		MembersCount: 1,
		Domain:       req.Data.Attributes.Domain,
		Metadata:     xo.Jsonb(req.Data.Attributes.Metadata),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	err = Storage(r).Transaction(func() error {
		err = Storage(r).UserQ().InsertCtx(r.Context(), &owner)
		if err != nil {
			return errors.Wrap(err, "failed to insert new user", logan.F{
				"did": req.Data.Attributes.OwnerDid,
			})
		}
		err = Storage(r).OrganizationQ().InsertCtx(r.Context(), &org)
		if err != nil {
			return errors.Wrap(err, "failed to insert new org", logan.F{
				"owner_did": req.Data.Attributes.OwnerDid,
				"domain":    req.Data.Attributes.Domain,
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

// removeScheme removes the scheme from a domain name.
func removeScheme(domain string) string {
	if strings.HasPrefix(domain, "http://") {
		return strings.TrimPrefix(domain, "http://")
	} else if strings.HasPrefix(domain, "https://") {
		return strings.TrimPrefix(domain, "https://")
	}
	return domain
}
