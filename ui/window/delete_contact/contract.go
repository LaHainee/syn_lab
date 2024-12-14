package delete_contact

import (
	"context"

	"fyne.io/fyne/v2"

	"contacts/internal/model"
)

type app interface {
	NewWindow(title string) fyne.Window
}

type contactList interface {
	Refresh()
}

type deleteHandler interface {
	Delete(ctx context.Context, uuid string) error
}

type fetchHandler interface {
	FetchByUuid(ctx context.Context, uuid string) (model.Contact, error)
}
