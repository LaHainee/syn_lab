package search

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

func (h *Handler) Search(_ context.Context, request model.SearchRequest) ([]model.Contact, error) {
	return h.storage.Search(request)
}
