package handlers

import (
	"github.com/google/uuid"
	"github.com/rarimo/rarime-orgs-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/urlval"
	"net/http"
)

type orgByIDRequest struct {
	ID uuid.UUID

	IncludeOwner bool `include:"owner"`
}

func newOrgByIDRequest(r *http.Request) (*orgByIDRequest, error) {
	var request orgByIDRequest

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode url")
	}

	request.ID, err = orgIDFromRequest(r)
	if err != nil {
		return nil, err
	}

	return &request, nil
}

func OrgByID(w http.ResponseWriter, r *http.Request) {
	req, err := newOrgByIDRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	org, err := Storage(r).OrganizationQ().OrganizationByIDCtx(r.Context(), req.ID, false)
	if err != nil {
		panic(errors.Wrap(err, "failed to get organization", logan.F{
			"org_id": req.ID,
		}))
	}
	if org == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	resp := resources.OrganizationResponse{
		Data:     populateOrg(*org),
		Included: resources.Included{},
	}

	if req.IncludeOwner {
		f := logan.F{
			"org_id":  req.ID,
			"user_id": org.Owner,
		}
		owner, err := Storage(r).UserQ().UserByIDCtx(r.Context(), org.Owner, false)
		if err != nil {
			panic(errors.Wrap(err, "failed to get owner", f))
		}
		if owner == nil {
			panic(errors.From(errors.New("expected owner exist"), f))
		}
		respOwner := populateUser(*owner)
		resp.Included.Add(&respOwner)
	}

	ape.Render(w, resp)
}
