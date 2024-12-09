package search

import "contacts/internal/model"

type storage interface {
	Search(request model.SearchRequest) ([]model.Contact, error)
}
