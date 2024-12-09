package menu

import "fyne.io/fyne/v2"

type app interface {
	Quit()
}

type contactList interface {
	SelectedContactUUID() *string
}

type createContactWindow interface {
	Build() fyne.Window
}

type deleteContactWindow interface {
	Build(contactUuid string) fyne.Window
}

type updateContactWindow interface {
	Build(contactUuid string) fyne.Window
}

type aboutWindow interface {
	Build() fyne.Window
}
