package create

import "contacts/internal/model"

type storage interface {
	Create(contact model.Contact) error
}

type validator interface {
	Validate(contact model.ContactForCreate) map[model.Field]string
}
