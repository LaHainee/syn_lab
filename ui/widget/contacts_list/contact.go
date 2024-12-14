package contacts_list

import (
	"context"

	"fyne.io/fyne/v2"

	"contacts/internal/model"
)

type appBox interface {
	Remove(rem fyne.CanvasObject)
	Add(add fyne.CanvasObject)
	Refresh()
}

type fetchHandler interface {
	Fetch(ctx context.Context) ([]model.Contact, error)
}

type searchHandler interface {
	Search(ctx context.Context, request model.SearchRequest) ([]model.Contact, error)
}
