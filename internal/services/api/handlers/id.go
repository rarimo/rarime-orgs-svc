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
	id, err := parseUUID(rawID)
	if err != nil {
		return uuid.UUID{}, errors.Wrap(err, "failed to parse organization id")
	}

	return id, nil
}

func groupIDFromRequest(r *http.Request) (uuid.UUID, error) {
	rawID := chi.URLParam(r, "group_id")
	id, err := parseUUID(rawID)
	if err != nil {
		return uuid.UUID{}, errors.Wrap(err, "failed to parse group id")
	}

	return id, nil
}

func reqIDFromRequest(r *http.Request) (uuid.UUID, error) {
	rawID := chi.URLParam(r, "req_id")
	id, err := parseUUID(rawID)
	if err != nil {
		return uuid.UUID{}, errors.Wrap(err, "failed to parse request id")
	}
	return id, nil
}

func parseUUID(rawID string) (uuid.UUID, error) {
	if rawID == "" {
		return uuid.UUID{}, validation.Errors{
			"id": errors.New("ID is required"),
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
