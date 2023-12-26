package handlers

import (
	"encoding/json"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/rarimo/rarime-orgs-svc/internal/data"
	"github.com/rarimo/rarime-orgs-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/urlval"
	"net/http"
	"strconv"
)

type orgListRequest struct {
	RawOwner *string `filter:"owner"`
	Owner    *uuid.UUID
	UserDID  *string                       `filter:"user_did"`
	Status   *resources.OrganizationStatus `filter:"status"`

	PageCursor uint64     `page:"cursor"`
	PageLimit  uint64     `page:"limit" default:"10"`
	Sorts      pgdb.Sorts `url:"sort" default:"-time"`
}

func newOrgListRequest(r *http.Request) (*orgListRequest, error) {
	var req orgListRequest

	err := urlval.Decode(r.URL.Query(), &req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode url")
	}

	if req.RawOwner != nil {
		owner, err := uuid.Parse(*req.RawOwner)
		if err != nil {
			return nil, validation.Errors{
				"owner": errors.Wrap(err, "failed to parse owner"),
			}
		}
		req.Owner = &owner
	}

	return &req, nil
}

func OrgList(w http.ResponseWriter, r *http.Request) {
	req, err := newOrgListRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	orgs, err := Storage(r).OrganizationQ().SelectCtx(r.Context(), createOrgsSelector(req))
	if err != nil {
		panic(errors.Wrap(err, "failed to select organizations"))
	}

	resp := resources.OrganizationListResponse{
		Data:     make([]resources.Organization, 0, len(orgs)),
		Included: resources.Included{},
		Links: &resources.Links{
			Self: fmt.Sprintf("%s?%s", r.URL.Path, urlval.MustEncode(req)),
		},
	}

	if len(orgs) == 0 {
		ape.Render(w, resp)
		return
	}

	for _, org := range orgs {
		resp.Data = append(resp.Data, populateOrg(org))
	}

	cursor := req.PageCursor + req.PageLimit
	req.PageCursor = cursor
	resp.Links.Next = fmt.Sprintf("%s?%s", r.URL.Path, urlval.MustEncode(req))

	_ = resp.PutMeta(map[string]interface{}{
		"next_cursor": cursor,
	})

	ape.Render(w, resp)
}

func populateOrg(org data.Organization) resources.Organization {
	attrs := resources.OrganizationAttributes{
		Domain:            org.Domain,
		Owner:             org.Owner.String(),
		Status:            resources.OrganizationStatus(org.Status),
		Metadata:          json.RawMessage(org.Metadata),
		MembersCount:      strconv.FormatInt(int64(org.MembersCount), 10),
		IssuedClaimsCount: strconv.FormatInt(org.IssuedClaimsCount, 10),
		UpdatedAt:         org.UpdatedAt,
		CreatedAt:         org.CreatedAt,
	}

	if org.Did.Valid {
		attrs.Did = &org.Did.String
	}

	if org.VerificationCode.Valid {
		attrs.VerificationCode = &org.VerificationCode.String
	}

	result := resources.Organization{
		Key: resources.Key{
			ID:   org.ID.String(),
			Type: resources.ORGANIZATIONS,
		},
		Attributes: attrs,
		Relationships: &resources.OrganizationRelationships{
			Owner: &resources.Relation{
				Data: &resources.Key{
					ID:   org.Owner.String(),
					Type: resources.USERS,
				},
			},
		},
	}

	return result
}

func createOrgsSelector(request *orgListRequest) data.OrgsSelector {
	sel := data.OrgsSelector{
		Owner:      request.Owner,
		UserDID:    request.UserDID,
		PageCursor: request.PageCursor,
		PageLimit:  request.PageLimit,
		Sort:       request.Sorts,
	}

	if request.Status != nil {
		sel.Status = request.Status.Intp()
	}

	return sel
}
