package handlers

import (
	"database/sql"
	"encoding/base64"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/rarimo/rarime-orgs-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"golang.org/x/crypto/sha3"
	"net/http"
	"time"
)

type orgVerificationCodeRequest struct {
	ID uuid.UUID
}

func newOrgVerificationCodeRequest(r *http.Request) (*orgVerificationCodeRequest, error) {
	var req orgVerificationCodeRequest
	var err error

	req.ID, err = orgIDFromRequest(r)
	if err != nil {
		return nil, err
	}

	return &req, nil
}

func OrgVerificationCode(w http.ResponseWriter, r *http.Request) {
	req, err := newOrgVerificationCodeRequest(r)
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

	code := createVerificationCode(org.ID, org.Domain)

	org.VerificationCode = sql.NullString{
		String: base64.StdEncoding.EncodeToString(code),
		Valid:  true,
	}
	org.UpdatedAt = time.Now().UTC()

	err = Storage(r).OrganizationQ().UpdateCtx(r.Context(), org)
	if err != nil {
		panic(errors.Wrap(err, "failed to update organization", logan.F{
			"org_id": req.ID,
		}))
	}

	ape.Render(w, resources.VerificationCodeResponse{
		Data: resources.VerificationCode{
			Key: resources.Key{
				ID:   uuid.UUID(code[:16]).String(),
				Type: resources.VERIFICATION_CODES,
			},
			Attributes: resources.VerificationCodeAttributes{
				Code: org.VerificationCode.String,
			},
		},
	})
}

func createVerificationCode(id uuid.UUID, domain string) []byte {
	hash := sha3.New256()
	_, err := hash.Write(append(id[:], []byte(domain)...))
	if err != nil {
		panic(errors.Wrap(err, "failed to write to hash", logan.F{
			"org_id": id,
			"domain": domain,
		}))
	}
	return hash.Sum(nil)
}
