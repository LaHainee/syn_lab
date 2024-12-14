package create

import (
	"context"
	"fmt"
	"time"

	"contacts/internal/model"
)

type Handler struct {
	storage   storage
	uuid      uuid
	validator validator
}

func NewHandler(s storage, uuid uuid, v validator) *Handler {
	return &Handler{
		storage:   s,
		uuid:      uuid,
		validator: v,
	}
}

func (h *Handler) Create(_ context.Context, contactForCreate model.ContactForCreate) (map[model.Field]string, error) {
	fieldMsgs := h.validator.Validate(contactForCreate)

	if len(fieldMsgs) > 0 {
		return fieldMsgs, model.ErrValidation
	}

	birthday, err := time.Parse("02.01.2006", contactForCreate.Birthday)
	if err != nil {
		return nil, err
	}

	phone, err := model.NewPhone(contactForCreate.Phone)
	if err != nil {
		return nil, err
	}

	contact := model.Contact{
		UUID:     h.uuid.NewString(),
		Surname:  contactForCreate.Surname,
		Name:     contactForCreate.Name,
		Birthday: birthday,
		Phone:    phone,
		Email:    contactForCreate.Email,
		Links:    contactForCreate.Links,
	}

	err = h.storage.Create(contact)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	return nil, nil
}
