package handlers

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/rarimo/rarime-orgs-svc/internal/services/core/issuer"
	"github.com/rarimo/rarime-orgs-svc/internal/services/core/issuer/models"
	"github.com/rarimo/rarime-orgs-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"

	"net"
	"net/http"
	"time"
)

type orgVerifyRequest struct {
	ID uuid.UUID
}

func newVerifyRequest(r *http.Request) (*orgVerifyRequest, error) {
	var req orgVerifyRequest
	var err error

	req.ID, err = orgIDFromRequest(r)
	if err != nil {
		return nil, err
	}

	return &req, nil
}

func OrgVerify(w http.ResponseWriter, r *http.Request) {
	cfgIssuer := IssuerConfig(r)
	cfgRarime := OrgsConfig(r)

	req, err := newVerifyRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	org, err := Storage(r).OrganizationQ().OrganizationByIDCtx(r.Context(), req.ID, true)
	if err != nil {
		panic(errors.Wrap(err, "failed to get organization", logan.F{
			"org_id": req.ID,
		}))
	}
	if org == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}
	if resources.OrganizationStatus(org.Status) == resources.OrganizationStatus_Verified {
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"status": errors.New("organization already verified"),
		})...)
		return
	}

	if cfgRarime.VerifyDomain {
		if !verifyCodeInTxtRecords("rarimo."+org.Domain, org.VerificationCode.String) {
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"status": errors.New("domain verification failed"),
			})...)
			return
		}
	}

	user, err := Storage(r).UserQ().UserByIDCtx(r.Context(), org.Owner, true)

	credentialSubject := models.NewEmptyDomainOwnershipCredentialSubject()
	credentialSubject.IdentityID = user.Did
	credentialSubject.Domain = org.Domain

	schema, err := Storage(r).ClaimSchemaQ().SchemaByActionTypeCtx(r.Context(), issuer.DomainOwnershipActionType)

	if err != nil {
		panic(errors.Wrap(err, "failed to get schema by action type", logan.F{
			"action_type": issuer.DomainOwnershipActionType,
		}))
	}
	if schema == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	iss := issuer.New(Log(r), &cfgIssuer, schema.SchemaType, schema.SchemaUrl)

	credentialReq := models.DomainOwnershipCredentialRequest{
		CredentialSchema:  schema.SchemaUrl,
		CredentialSubject: credentialSubject,
		Type:              schema.SchemaType,
	}

	claim, err := iss.IssueClaim(user.Did, credentialReq)
	if err != nil {
		Log(r).Error("failed to issue claim:", err.Error())
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if claim == nil {
		ape.RenderErr(w, problems.InternalError())
		return
	}

	org.Status = resources.OrganizationStatus_Verified.Int16()
	org.UpdatedAt = time.Now().UTC()

	err = Storage(r).OrganizationQ().UpdateCtx(r.Context(), org)
	if err != nil {
		panic(errors.Wrap(err, "failed to update organization", logan.F{
			"org_id": req.ID,
		}))
	}

	jsonClaim, err := json.Marshal(claim)
	if err != nil {
		panic(errors.Wrap(err, "failed to marshal claim", logan.F{
			"org_id": req.ID,
		}))
	}

	var rawMsg json.RawMessage
	err = json.Unmarshal(jsonClaim, &rawMsg)

	inc := resources.Included{}
	inc.Add(resources.ClaimOffer{
		Id:   claim.Data.ID,
		Type: string(claim.Data.Type),
	})

	ape.Render(w, resources.OrganizationResponse{
		Data:     populateOrg(*org),
		Included: inc,
	})
}

func verifyCodeInTxtRecords(domain string, verifyCode string) bool {
	txtrecords, err := net.LookupTXT(domain)
	if err != nil {
		return false
	}

	if len(txtrecords) == 0 {
		return false
	}

	for _, txtrecord := range txtrecords {
		if txtrecord == verifyCode {
			return true
		}
	}

	return false
}
