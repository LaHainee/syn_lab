package update_contact

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

type updateHandler interface {
	Update(ctx context.Context, contactForCreate model.ContactForCreate) (map[model.Field]string, error)
}

type fetchHandler interface {
	Fetch(ctx context.Context) ([]model.Contact, error)
	FetchByUuid(ctx context.Context, uuid string) (model.Contact, error)
}
