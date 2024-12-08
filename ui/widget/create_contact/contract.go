package create_contact

import (
	"context"

	"contacts/internal/model"
)

type handler interface {
	Create(_ context.Context, contact model.ContactForCreate) (map[model.Field]string, error)
}

type storage interface {
	Fetch() ([]model.Contact, error)
}
