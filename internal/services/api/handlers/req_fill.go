package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/rarimo/rarime-orgs-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func newReqFillRequest(r *http.Request) (*resources.RequestFillRequest, error) {
	var req resources.RequestFillRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(err, "failed to decode body")
	}

	id, err := reqIDFromRequest(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get request id")
	}

	req.Data.ID = id.String()

	return &req, nil
}

func RequestFill(w http.ResponseWriter, r *http.Request) {
	// TODO: add auth for the user that is associated with the request
	req, err := newReqFillRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	id := uuid.MustParse(req.Data.ID)

	request, err := Storage(r).RequestQ().RequestByIDCtx(r.Context(), id, true)
	if err != nil {
		panic(errors.Wrap(err, "failed to get request", logan.F{
			"req_id": id,
		}))
	}
	if request == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}
	if request.Status != resources.RequestStatus_Accepted.Int16() {
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"data/attributes/status": errors.New("request is not accepted"),
		})...)
		return
	}

	request.Status = resources.RequestStatus_Filled.Int16()

	credentialRequestsJSON, err := json.Marshal(req.Data.Attributes.CredentialRequests)
	if err != nil {
		panic(errors.Wrap(err, "failed to marshal credential requests"))
	}

	request.Metadata = credentialRequestsJSON
	request.UpdatedAt = time.Now().UTC()

	if err := Storage(r).RequestQ().UpdateCtx(r.Context(), request); err != nil {
		panic(errors.Wrap(err, "failed to update request", logan.F{
			"req_id": id,
		}))
	}

	respRequest, err := populateRequest(*request)
	if err != nil {
		Log(r).WithError(err).Error("failed to populate request")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, resources.RequestResponse{
		Data:     respRequest,
		Included: resources.Included{},
	})
}
