package handlers

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/rarimo/rarime-orgs-svc/internal/services/core/issuer/schemas"
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

	if OrgsConfig(r).VerifyDomain {
		if !verifyCodeInTxtRecords("rarimo."+org.Domain, org.VerificationCode.String) {
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"status": errors.New("domain verification failed"),
			})...)
			return
		}
	}

	user, err := Storage(r).UserQ().UserByIDCtx(r.Context(), org.Owner, true)
	if err != nil {
		Log(r).WithError(err).Error("Failed to get user by ID")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if user == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	schema, err := Storage(r).ClaimSchemaQ().SchemaByActionTypeCtx(r.Context(), schemas.ActionDomainOwnership)
	if err != nil {
		Log(r).WithError(err).Error("Failed to get schema by action type: %s", schemas.ActionDomainOwnership)
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if schema == nil {
		Log(r).Errorf("Schema not found by action type: %s", schemas.ActionDomainOwnership)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	credentialReq := schemas.CreateCredentialRequest{
		CredentialSchema: schema.SchemaUrl,
		CredentialSubject: schemas.DomainOwnership{
			IdentityID: user.Did,
			Domain:     org.Domain,
		},
		Type: schema.SchemaType,
	}

	claim, err := Issuer(r).IssueClaim(user.Did, credentialReq)
	if err != nil {
		Log(r).WithError(err).Error("Failed to issue claim")
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
