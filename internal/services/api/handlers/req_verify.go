package handlers

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
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

func newReqVerifyRequest(r *http.Request) (*resources.RequestVerifyRequest, error) {
	var req resources.RequestVerifyRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(err, "failed to decode body")
	}

	id, err := reqIDFromRequest(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get request id")
	}

	req.Data.ID = id.String()

	if !req.Data.Attributes.Approved {
		return &req, nil
	}
	if req.Data.Attributes.Metadata == nil {
		return nil, validation.Errors{
			"data/attributes/metadata": errors.New("metadata is required for approved request"),
		}
	}
	if req.Data.Attributes.Role == nil {
		return nil, validation.Errors{
			"data/attributes/role": errors.New("role is required for approved request"),
		}
	}
	if valid := json.Valid(*req.Data.Attributes.Metadata); !valid {
		return nil, validation.Errors{
			"data/attributes/metadata": errors.New("invalid request metadata json"),
		}
	}

	return &req, nil
}

func RequestVerify(w http.ResponseWriter, r *http.Request) {
	// TODO: add auth for the group admin or owner or superadmin
	req, err := newReqVerifyRequest(r)
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
	if request.Status != resources.RequestStatus_Filled.Int16() {
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"data/attributes/status": errors.New("request is not filled"),
		})...)
		return
	}

	request.Status = resources.RequestStatus_Rejected.Int16()
	request.UpdatedAt = time.Now().UTC()

	if req.Data.Attributes.Approved {
		request.Status = resources.RequestStatus_Approved.Int16()
		request.Metadata = xo.Jsonb(*req.Data.Attributes.Metadata)
	}

	user, err := Storage(r).UserQ().UserByDidCtx(r.Context(), request.UserDid.String)
	if err != nil {
		panic(errors.Wrap(err, "failed to get user", logan.F{
			"req_id":   id,
			"user_did": request.UserDid.String,
		}))
	}
	if user == nil {
		panic(errors.Wrap(err, "expected that user exist", logan.F{
			"req_id":   id,
			"user_did": request.UserDid.String,
		}))
	}

	err = Storage(r).Transaction(func() error {
		err = Storage(r).RequestQ().UpdateCtx(r.Context(), request)
		if err != nil {
			return errors.Wrap(err, "failed to update request", logan.F{
				"req_id": id,
			})
		}
		if !req.Data.Attributes.Approved {
			return nil
		}

		groupUser := data.GroupUser{
			ID:        uuid.New(),
			GroupID:   request.GroupID,
			UserID:    user.ID,
			Role:      resources.GroupUserRoleFromString(*req.Data.Attributes.Role).Int16(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}

		err = Storage(r).GroupUserQ().InsertCtx(r.Context(), &groupUser)
		if err != nil {
			return errors.Wrap(err, "failed to insert group user", logan.F{
				"group_id": request.GroupID,
				"user_id":  user.ID,
			})
		}

		return nil
	})
	if err != nil {
		Log(r).WithError(err).Error("failed to update the request")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	// TODO: add logic for change request status to submitted if it was approved

	populatedRequest, err := populateRequest(*request)
	if err != nil {
		Log(r).WithError(err).Error("failed to populate request")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, resources.RequestResponse{
		Data:     populatedRequest,
		Included: resources.Included{},
	})
}
