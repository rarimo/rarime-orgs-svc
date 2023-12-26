package handlers

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/rarimo/rarime-orgs-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
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

	// TODO: Add verification logic here

	org.Status = resources.OrganizationStatus_Verified.Int16()
	org.UpdatedAt = time.Now().UTC()

	err = Storage(r).OrganizationQ().UpdateCtx(r.Context(), org)
	if err != nil {
		panic(errors.Wrap(err, "failed to update organization", logan.F{
			"org_id": req.ID,
		}))
	}

	ape.Render(w, resources.OrganizationResponse{
		Data: populateOrg(*org),
	})
}
