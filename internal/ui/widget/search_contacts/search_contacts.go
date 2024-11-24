package search_contacts

import (
	"contacts/internal/model"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type Builder struct {
	contacts []model.Contact
}

func NewBuilder(contacts []model.Contact) *Builder {
	return &Builder{}
}

func (b *Builder) Build() *fyne.Container {
	entry := widget.NewEntry()

	entry.OnChanged = func(text string) {

	}
}
