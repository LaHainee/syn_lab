//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package fetch

import "contacts/internal/model"

type storage interface {
	Fetch() ([]model.Contact, error)
	FetchByUuid(uuid string) (model.Contact, error)
}
