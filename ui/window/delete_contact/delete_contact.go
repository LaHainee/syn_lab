package delete_contact

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"contacts/internal/storage"
	contactsList "contacts/ui/widget/contacts_list"
)

var (
	windowSize = fyne.NewSize(300, 100)
	buttonSize = fyne.NewSize(70, 30)
)

type Builder struct {
	// Объект приложения
	app fyne.App

	storage *storage.Storage

	contactListBuilder *contactsList.Builder
}

func NewBuilder(app fyne.App, storage *storage.Storage, contactsListBuilder *contactsList.Builder) *Builder {
	return &Builder{
		app:                app,
		storage:            storage,
		contactListBuilder: contactsListBuilder,
	}
}

func (b *Builder) Build(contactUuid string) fyne.Window {
	contact, err := b.storage.FetchByUUID(contactUuid)
	if err != nil {
		panic(err)
	}

	window := b.app.NewWindow("Delete contact")
	window.Resize(windowSize)
	window.SetFixedSize(true)
	window.CenterOnScreen()

	label := widget.NewLabel(fmt.Sprintf("Удалить контакт %s %s?", contact.Name, contact.Surname))

	closeButton := widget.NewButton("Cancel", func() {
		window.Close()
	})
	closeButton.Resize(buttonSize)
	closeButton.Move(
		fyne.NewPos(
			windowSize.Width-buttonSize.Width-20,
			windowSize.Height-buttonSize.Height-20,
		),
	)

	confirmButton := widget.NewButton("OK", func() {
		err = b.storage.Delete(contactUuid)
		if err != nil {
			panic(err)
		}

		contacts, err := b.storage.Fetch()
		if err != nil {
			panic(err)
		}

		b.contactListBuilder.Refresh(contacts)

		window.Close()
	})
	confirmButton.Resize(buttonSize)
	confirmButton.Move(
		fyne.NewPos(
			windowSize.Width-buttonSize.Width*2-20*2,
			windowSize.Height-buttonSize.Height-20,
		),
	)

	box := container.NewWithoutLayout()
	box.Add(label)
	box.Add(closeButton)
	box.Add(confirmButton)

	window.SetContent(box)

	return window
}
