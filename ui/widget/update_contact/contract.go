package update_contact

import (
	"context"

	"contacts/internal/model"
)

type handler interface {
	Update(_ context.Context, contactForCreate model.ContactForCreate) (map[model.Field]string, error)
}

type storage interface {
	Fetch() ([]model.Contact, error)
	FetchByUUID(uuid string) (model.Contact, error)
}
