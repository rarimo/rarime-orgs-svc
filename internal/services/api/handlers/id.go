package handlers

import (
	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
)

func orgIDFromRequest(r *http.Request) (uuid.UUID, error) {
	rawID := chi.URLParam(r, "id")
	if rawID == "" {
		return uuid.UUID{}, validation.Errors{
			"id": errors.New("Organization ID is required"),
		}
	}

	id, err := uuid.Parse(rawID)
	if err != nil {
		return uuid.UUID{}, validation.Errors{
			"id": errors.Wrap(err, "failed to parse id", logan.F{
				"id": rawID,
			}),
		}
	}

	return id, nil
}
