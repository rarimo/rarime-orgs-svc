package handlers

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/rarimo/rarime-orgs-svc/internal/data"
	"github.com/rarimo/rarime-orgs-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/urlval"
	"net/http"
)

type requestListRequest struct {
	OrgID   uuid.UUID
	GroupID uuid.UUID

	UserDID *string                  `filter:"user_did"`
	Status  *resources.RequestStatus `filter:"status"`

	PageCursor uint64     `page:"cursor"`
	PageLimit  uint64     `page:"limit" default:"10"`
	Sorts      pgdb.Sorts `url:"sort" default:"-time"`

	IncludeOrganization bool `include:"organization"`
	IncludeGroup        bool `include:"group"`
}

func newRequestListRequest(r *http.Request) (*requestListRequest, error) {
	var req requestListRequest

	err := urlval.Decode(r.URL.Query(), &req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode url")
	}

	req.OrgID, err = orgIDFromRequest(r)
	if err != nil {
		return nil, err
	}

	req.GroupID, err = groupIDFromRequest(r)
	if err != nil {
		return nil, err
	}

	if req.UserDID != nil {
		err = validation.Validate(req.UserDID, ValidationDID)
		if err != nil {
			return nil, validation.Errors{
				"filter[user_did]": errors.Wrap(err, "invalid user_did filter value"),
			}
		}
	}

	return &req, nil
}

func RequestList(w http.ResponseWriter, r *http.Request) {
	req, err := newRequestListRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	org, err := Storage(r).OrganizationQ().OrganizationByIDCtx(r.Context(), req.OrgID, false)
	if err != nil {
		panic(errors.Wrap(err, "failed to get organization", logan.F{
			"org_id": req.OrgID,
		}))
	}
	if org == nil {
		ape.RenderErr(w, NotFound(fmt.Sprintf("Organization with ID: %s not exist", req.OrgID), "id"))
		return
	}

	group, err := Storage(r).GroupQ().GroupByIDCtx(r.Context(), req.GroupID, false)
	if err != nil {
		panic(errors.Wrap(err, "failed to get group", logan.F{
			"group_id": req.GroupID,
		}))
	}
	if group == nil {
		ape.RenderErr(w, NotFound(fmt.Sprintf("Group with ID: %s not exist", req.GroupID), "group_id"))
		return
	}

	requests, err := Storage(r).RequestQ().SelectCtx(r.Context(), createRequestsSelector(req))
	if err != nil {
		panic(errors.Wrap(err, "failed to select requests"))
	}

	resp := resources.RequestListResponse{
		Data:     make([]resources.Request, 0, len(requests)),
		Included: resources.Included{},
		Links: &resources.Links{
			Self: fmt.Sprintf("%s?%s", r.URL.Path, urlval.MustEncode(req)),
		},
	}
	if len(requests) == 0 {
		ape.Render(w, resp)
		return
	}

	for _, request := range requests {
		populatedRequest, err := populateRequest(request)
		if err != nil {
			Log(r).WithError(err).Error("failed to populate request")
			ape.RenderErr(w, problems.InternalError())
			return
		}
		resp.Data = append(resp.Data, populatedRequest)
	}

	cursor := req.PageCursor + req.PageLimit
	req.PageCursor = cursor
	resp.Links.Next = fmt.Sprintf("%s?%s", r.URL.Path, urlval.MustEncode(req))

	_ = resp.PutMeta(map[string]interface{}{
		"next_cursor": cursor,
	})

	if req.IncludeOrganization {
		respOrg := populateOrg(*org)
		resp.Included.Add(&respOrg)
	}
	if req.IncludeGroup {
		respGroup := populateGroup(*group)
		resp.Included.Add(&respGroup)
	}

	ape.Render(w, resp)
}

func createRequestsSelector(req *requestListRequest) data.RequestsSelector {
	sel := data.RequestsSelector{
		OrgID:      &req.OrgID,
		GroupID:    &req.GroupID,
		UserDID:    req.UserDID,
		PageCursor: req.PageCursor,
		PageLimit:  req.PageLimit,
		Sort:       req.Sorts,
	}

	if req.Status != nil {
		sel.Status = req.Status.IntP()
	}

	return sel
}
