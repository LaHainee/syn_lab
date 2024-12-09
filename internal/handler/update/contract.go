//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package update

import "contacts/internal/model"

type storage interface {
	Update(contact model.Contact) error
}

type validator interface {
	Validate(contact model.ContactForCreate) map[model.Field]string
}
