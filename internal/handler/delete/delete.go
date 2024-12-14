package delete

import "context"

type Handler struct {
	storage storage
}

func NewHandler(s storage) *Handler {
	return &Handler{
		storage: s,
	}
}

func (h *Handler) Delete(_ context.Context, uuid string) error {
	return h.storage.Delete(uuid)
}
