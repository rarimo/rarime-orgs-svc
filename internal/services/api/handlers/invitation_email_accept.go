package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	rules "github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"github.com/rarimo/rarime-orgs-svc/internal/data"
	"github.com/rarimo/rarime-orgs-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
	"time"
)

func newInvitationEmailAcceptRequest(r *http.Request) (*uuid.UUID, *uuid.UUID, *resources.InvitationAcceptEmailRequest, error) {
	var req resources.InvitationAcceptEmailRequest

	orgID, err := orgIDFromRequest(r)
	if err != nil {
		return nil, nil, nil, err
	}

	groupID, err := groupIDFromRequest(r)
	if err != nil {
		return nil, nil, nil, err
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to decode body")
	}

	return &orgID, &groupID, &req, validation.Errors{
		"data/id":                  validation.Validate(req.Data.ID, validation.Required, rules.UUID),
		"data/attributes/otp":      validation.Validate(req.Data.Attributes.Otp, validation.Required, validation.Length(data.MaxDigits, data.MaxDigits), rules.Digit),
		"data/attributes/user_did": validation.Validate(req.Data.Attributes.UserDid, validation.Required, ValidationDID),
	}.Filter()
}

func InvitationEmailAccept(w http.ResponseWriter, r *http.Request) {
	orgID, groupID, req, err := newInvitationEmailAcceptRequest(r)
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
		ape.RenderErr(w, NotFound(fmt.Sprintf("Organization with ID: %s not exist", orgID), "id"))
		return
	}
	if org.Status != resources.OrganizationStatus_Verified.Int16() {
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"id": errors.Errorf("organization: %s is not verified", org.ID),
		})...)
		return
	}

	group, err := Storage(r).GroupQ().GroupByIDCtx(r.Context(), *groupID, false)
	if err != nil {
		panic(errors.Wrap(err, "failed to get group", logan.F{
			"group_id": groupID,
		}))
	}
	if group == nil {
		ape.RenderErr(w, NotFound(fmt.Sprintf("Group with ID: %s not exist", orgID), "group_id"))
		return
	}

	inv, err := Storage(r).EmailInvitationQ().EmailInvitationByIDCtx(r.Context(), uuid.MustParse(req.Data.ID), false)
	if err != nil {
		panic(errors.Wrap(err, "failed to get invitation", logan.F{
			"inv_id": req.Data.ID,
		}))
	}
	if inv == nil {
		ape.RenderErr(w, NotFound(fmt.Sprintf("Invitation with ID: %s not exist", req.Data.ID), "id"))
		return
	}

	request, err := Storage(r).RequestQ().RequestByIDCtx(r.Context(), inv.ReqID, true)
	if err != nil {
		panic(errors.Wrap(err, "failed to get request", logan.F{
			"req_id": inv.ReqID,
		}))
	}
	if request == nil {
		panic(errors.Wrap(err, "expected that request exists", logan.F{
			"inv_id":   inv.ID,
			"req_id":   inv.ReqID,
			"org_id":   org.ID,
			"group_id": group.ID,
		}))
	}
	if req.Data.Attributes.Otp != inv.Otp {
		ape.RenderErr(w, problems.Forbidden())
		return
	}
	// TODO: add claims issuing logic here
	user := data.User{
		ID:        uuid.New(),
		Did:       req.Data.Attributes.UserDid,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	request.UserDid = sql.NullString{
		String: user.Did,
		Valid:  true,
	}
	request.Status = resources.RequestStatus_Accepted.Int16()
	request.UpdatedAt = time.Now().UTC()

	err = Storage(r).Transaction(func() error {
		err = Storage(r).UserQ().InsertCtx(r.Context(), &user)
		if err != nil {
			return errors.Wrap(err, "failed to insert new user", logan.F{
				"did": req.Data.Attributes.UserDid,
			})
		}
		err = Storage(r).RequestQ().UpdateCtx(r.Context(), request)
		if err != nil {
			return errors.Wrap(err, "failed to update request", logan.F{
				"req_id": request.ID,
			})
		}

		return nil
	})
	if err != nil {
		Log(r).WithError(err).Error("failed to create new user and update request")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	resp := resources.InvitationEmailResponse{
		Data:     populateInvitationEmail(*inv),
		Included: resources.Included{},
	}
	respRequest := populateRequest(*request)
	resp.Included.Add(&respRequest)

	ape.Render(w, resp)
}
