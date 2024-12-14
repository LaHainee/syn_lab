package fetch

import (
	"context"

	"contacts/internal/model"
)

type Handler struct {
	storage storage
}

func NewHandler(s storage) *Handler {
	return &Handler{
		storage: s,
	}
}

func (h *Handler) Fetch(_ context.Context) ([]model.Contact, error) {
	return h.storage.Fetch()
}

func (h *Handler) FetchByUuid(_ context.Context, uuid string) (model.Contact, error) {
	return h.storage.FetchByUuid(uuid)
}
