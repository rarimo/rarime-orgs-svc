package handlers

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/rarimo/rarime-orgs-svc/internal/services/api/issuer"
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
		txtrecords, err := net.LookupTXT(org.Domain)
		if err != nil {
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"domain": errors.New("domain not verified"),
			})...)
			return
		}

		if len(txtrecords) == 0 {
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"domain": errors.New("domain not verified"),
			})...)
			return
		}

		var found bool

		for _, txtrecord := range txtrecords {
			if txtrecord == org.VerificationCode.String {
				found = true
				break
			}
		}

		if !found {
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"domain": errors.New("domain not verified"),
			})...)
			return
		}
	}

	credentialSubject := issuer.IdentityProvidersCredentialSubject{
		IdentityID: org.Did.String,
	}

	iss := issuer.New(Log(r), &cfgIssuer)

	sigProof := true
	exp := time.Now().UTC().Add(time.Hour * 24 * 365 * 10)

	credentialReq := issuer.CreateCredentialRequest{
		CredentialSchema:  iss.SchemaURL(),
		CredentialSubject: &credentialSubject,
		Type:              iss.SchemaType(),
		Expiration:        &exp,
		SignatureProof:    &sigProof,
		MtProof:           &sigProof,
	}

	claim, err := iss.IssueClaim(credentialReq)
	if err != nil {
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
