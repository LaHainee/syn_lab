package update

import (
	"context"
	"errors"
	"fmt"
	"time"

	"contacts/internal/model"
)

type Handler struct {
	storage   storage
	validator validator
}

func NewHandler(s storage, v validator) *Handler {
	return &Handler{
		storage:   s,
		validator: v,
	}
}

func (h *Handler) Update(_ context.Context, contactForCreate model.ContactForCreate) (map[model.Field]string, error) {
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

	if contactForCreate.UUID == nil {
		return nil, errors.New("empty uuid")
	}

	contact := model.Contact{
		UUID:     *contactForCreate.UUID,
		Surname:  contactForCreate.Surname,
		Name:     contactForCreate.Name,
		Birthday: birthday,
		Phone:    phone,
		Email:    contactForCreate.Email,
		Links:    contactForCreate.Links,
	}

	err = h.storage.Update(contact)
	if err != nil {
		return nil, fmt.Errorf("update: %w", err)
	}

	return nil, nil
}
