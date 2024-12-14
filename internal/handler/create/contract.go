//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package create

import "contacts/internal/model"

type storage interface {
	Create(contact model.Contact) error
}

type validator interface {
	Validate(contact model.ContactForCreate) map[model.Field]string
}

type uuid interface {
	NewString() string
}
