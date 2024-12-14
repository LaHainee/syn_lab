//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package search

import "contacts/internal/model"

type storage interface {
	Search(request model.SearchRequest) ([]model.Contact, error)
}
