package handlers

import (
	"encoding/json"
	"fmt"
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
	"time"
)

func newInvitationEmailCreateRequest(r *http.Request) (*uuid.UUID, *uuid.UUID, *resources.InvitationCreateEmailRequest, error) {
	var req resources.InvitationCreateEmailRequest

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

	if req.Data.Attributes.Rules != nil {
		if valid := json.Valid(*req.Data.Attributes.Rules); !valid {
			return nil, nil, nil, validation.Errors{
				"data/attributes/rules": errors.New("invalid rules json"),
			}
		}
	}

	return &orgID, &groupID, &req, validation.Errors{
		"data/attributes/email": validation.Validate(req.Data.Attributes.Email, validation.Required, rules.EmailFormat),
	}
}

func InvitationEmailCreate(w http.ResponseWriter, r *http.Request) {
	// TODO: add auth check for owner or group admin or superadmin permissions.
	orgID, groupID, req, err := newInvitationEmailCreateRequest(r)
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

	request := data.Request{
		ID:        uuid.New(),
		OrgID:     org.ID,
		GroupID:   group.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if req.Data.Attributes.Rules != nil {
		request.Metadata = xo.Jsonb(*req.Data.Attributes.Rules)
	}

	inv := data.EmailInvitation{
		ID:        uuid.New(),
		ReqID:     request.ID,
		OrgID:     org.ID,
		GroupID:   group.ID,
		Email:     req.Data.Attributes.Email,
		Otp:       data.GenerateOTP(),
		CreatedAt: time.Now().UTC(),
	}

	err = Storage(r).Transaction(func() error {
		err = Storage(r).RequestQ().InsertCtx(r.Context(), &request)
		if err != nil {
			return errors.Wrap(err, "failed to insert new request", logan.F{
				"org_id":   org.ID,
				"group_id": group.ID,
			})
		}
		err = Storage(r).EmailInvitationQ().InsertCtx(r.Context(), &inv)
		if err != nil {
			return errors.Wrap(err, "failed to insert new email invitation", logan.F{
				"org_id":   org.ID,
				"group_id": group.ID,
			})
		}

		return nil
	})
	if err != nil {
		Log(r).WithError(err).Error("failed to create new request")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	Log(r).WithFields(logan.F{
		"otp":      inv.Otp,
		"email":    req.Data.Attributes.Email,
		"inv_id":   inv.ID,
		"org_id":   org.ID,
		"group_id": group.ID,
	}).Debug("new invitation email")

	// TODO: add email sending logic here
	resp := resources.InvitationEmailResponse{
		Data:     populateInvitationEmail(inv),
		Included: resources.Included{},
	}
	respRequest := populateRequest(request)
	resp.Included.Add(&respRequest)

	ape.Render(w, resp)
}

func populateInvitationEmail(invite data.EmailInvitation) resources.InvitationEmail {
	return resources.InvitationEmail{
		Key: resources.Key{
			ID:   invite.ID.String(),
			Type: resources.INVITATIONS_EMAIL,
		},
		Attributes: resources.InvitationEmailAttributes{
			Email:     invite.Email,
			GroupId:   invite.GroupID.String(),
			OrgId:     invite.OrgID.String(),
			ReqId:     invite.ReqID.String(),
			CreatedAt: invite.CreatedAt,
		},
		Relationships: &resources.InvitationEmailRelationships{
			Request: resources.Relation{
				Data: &resources.Key{
					ID:   invite.ReqID.String(),
					Type: resources.REQUESTS,
				},
			},
		},
	}
}

func populateRequest(request data.Request) resources.Request {
	req := resources.Request{
		Key: resources.Key{
			ID:   request.ID.String(),
			Type: resources.REQUESTS,
		},
		Attributes: resources.RequestAttributes{
			GroupId:   request.GroupID.String(),
			OrgId:     request.OrgID.String(),
			Metadata:  json.RawMessage(request.Metadata),
			Status:    resources.RequestStatus(request.Status),
			UpdatedAt: request.UpdatedAt,
			CreatedAt: request.CreatedAt,
		},
		Relationships: &resources.RequestRelationships{
			Group: &resources.Relation{
				Data: &resources.Key{
					ID:   request.GroupID.String(),
					Type: resources.GROUPS,
				},
			},
			Organization: &resources.Relation{
				Data: &resources.Key{
					ID:   request.OrgID.String(),
					Type: resources.ORGANIZATIONS,
				},
			},
		},
	}

	if request.UserDid.Valid {
		req.Attributes.UserDid = &request.UserDid.String
	}

	return req
}
