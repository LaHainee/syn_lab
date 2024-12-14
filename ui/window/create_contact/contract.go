package create_contact

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

type createHandler interface {
	Create(ctx context.Context, contact model.ContactForCreate) (map[model.Field]string, error)
}
