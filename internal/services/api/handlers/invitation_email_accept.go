package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	rules "github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"github.com/rarimo/rarime-orgs-svc/internal/data"
	"github.com/rarimo/rarime-orgs-svc/internal/services/core/issuer"
	"github.com/rarimo/rarime-orgs-svc/internal/services/core/issuer/schemas"
	"github.com/rarimo/rarime-orgs-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
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
		ape.RenderErr(w, NotFound(fmt.Sprintf("Group with ID: %s not exist", *groupID), "group_id"))
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

	user := data.User{
		ID:        uuid.New(),
		Did:       req.Data.Attributes.UserDid,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	emailClaim, roleClaim, err := issueInvitationAcceptClaims(r, user, inv.Email)
	if err != nil {
		Log(r).WithError(err).Error("Failed to issue claims on invitation accept")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	request.UserDid = sql.NullString{
		String: user.Did,
		Valid:  true,
	}
	request.Status = resources.RequestStatus_Accepted.Int16()
	request.UpdatedAt = time.Now().UTC()

	if err = updateOnInvitationAccept(r, user, *request); err != nil {
		Log(r).WithError(err).Error("Failed to update data on invitation accept")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	resp := resources.InvitationEmailResponse{Data: populateInvitationEmail(*inv)}
	respRequest := populateRequest(*request)
	resp.Included.Add(&respRequest)
	resp.Included.Add(&resources.ClaimOffer{
		Id:   emailClaim.Data.ID,
		Type: string(emailClaim.Data.Type),
	})
	resp.Included.Add(&resources.ClaimOffer{
		Id:   roleClaim.Data.ID,
		Type: string(roleClaim.Data.Type),
	})

	ape.Render(w, resp)
}

func issueInvitationAcceptClaims(
	r *http.Request,
	user data.User,
	userEmail string,
) (email, role *issuer.IssueClaimResponse, err error) {

	schema, err := getSchemaByActionType(r, schemas.ActionEmployeeNickname)
	if err != nil {
		return nil, nil, fmt.Errorf("get schema by action type: %w", err)
	}

	var (
		req = schemas.CreateCredentialRequest{
			CredentialSchema: schema.SchemaUrl,
			CredentialSubject: schemas.EmployeeNickname{
				EmployeeSocialMedia: schemas.SocialMediaEmail,
				EmployeeNickname:    userEmail,
			},
			Type: schema.SchemaType,
		}
		errCtx = fmt.Sprintf("[did=%s; type=%s; schema_url=%s; subject=%v]",
			user.Did, req.Type, req.CredentialSchema, req.CredentialSubject)
	)

	email, err = Issuer(r).IssueClaim(user.Did, req)
	if err != nil {
		return nil, nil, fmt.Errorf("issue claim %s: %w", errCtx, err)
	}

	req.CredentialSubject = schemas.UserRole{Role: schemas.RoleUndefined}
	role, err = Issuer(r).IssueClaim(user.Did, req)
	if err != nil {
		return nil, nil, fmt.Errorf("issue claim %s: %w", errCtx, err)
	}

	return
}

func updateOnInvitationAccept(r *http.Request, user data.User, request data.Request) error {
	return Storage(r).Transaction(func() error {
		err := Storage(r).UserQ().InsertCtx(r.Context(), &user)
		if err != nil {
			return fmt.Errorf("insert user: %w", err)
		}

		err = Storage(r).RequestQ().UpdateCtx(r.Context(), &request)
		if err != nil {
			return fmt.Errorf("update request: %w", err)
		}

		return nil
	})
}

func getSchemaByActionType(r *http.Request, actionType string) (*data.ClaimSchema, error) {
	schema, err := Storage(r).ClaimSchemaQ().SchemaByActionTypeCtx(r.Context(), actionType)
	if err != nil {
		return nil, fmt.Errorf("get schema by action_type=%s: %w", actionType, err)
	}
	if schema == nil {
		return nil, fmt.Errorf("schema not found by action_type=%s", actionType)
	}
	return schema, nil
}
