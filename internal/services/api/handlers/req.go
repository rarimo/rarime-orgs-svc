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

type reqByIDRequest struct {
	ID uuid.UUID

	IncludeOrganization bool `include:"organization"`
	IncludeGroup        bool `include:"group"`
}

func newReqByIDRequest(r *http.Request) (*reqByIDRequest, error) {
	var request reqByIDRequest

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode url")
	}

	request.ID, err = reqIDFromRequest(r)
	if err != nil {
		return nil, err
	}

	return &request, nil
}

func RequestByID(w http.ResponseWriter, r *http.Request) {
	// TODO: add auth check for user with any role but that is in the users list for specified organization or group admin, owner, superadmin
	req, err := newReqByIDRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	request, err := Storage(r).RequestQ().RequestByIDCtx(r.Context(), req.ID, false)
	if err != nil {
		panic(errors.Wrap(err, "failed to get request", logan.F{
			"req_id": req.ID,
		}))
	}
	if request == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	respRequest, err := populateRequest(*request)
	if err != nil {
		Log(r).WithError(err).Error("failed to populate request")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	resp := resources.RequestResponse{
		Data:     respRequest,
		Included: resources.Included{},
	}

	if req.IncludeOrganization {
		org, err := Storage(r).OrganizationQ().OrganizationByIDCtx(r.Context(), request.OrgID, false)
		if err != nil {
			panic(errors.Wrap(err, "failed to get organization", logan.F{
				"org_id": request.OrgID,
			}))
		}
		if org == nil {
			panic(errors.Wrap(err, "expected that org exists", logan.F{
				"req_id": req.ID,
				"org_id": request.OrgID,
			}))
		}

		respOrg := populateOrg(*org)
		resp.Included.Add(&respOrg)
	}
	if req.IncludeGroup {
		group, err := Storage(r).GroupQ().GroupByIDCtx(r.Context(), request.GroupID, false)
		if err != nil {
			panic(errors.Wrap(err, "failed to get group", logan.F{
				"org_id":   request.OrgID,
				"group_id": request.GroupID,
			}))
		}
		if group == nil {
			panic(errors.Wrap(err, "expected that group exists", logan.F{
				"req_id":   req.ID,
				"org_id":   request.OrgID,
				"group_id": request.GroupID,
			}))
		}

		respGroup := populateGroup(*group)
		resp.Included.Add(&respGroup)
	}

	ape.Render(w, resp)
}
