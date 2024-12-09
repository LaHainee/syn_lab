package menu

import (
	"fyne.io/fyne/v2"
)

type Builder struct {
	app                 app
	contactList         contactList
	createContactWindow createContactWindow
	updateContactWindow updateContactWindow
	deleteContactWindow deleteContactWindow
	aboutWindow         aboutWindow
}

func NewBuilder(
	app app,
	contactList contactList,
	createContactWindow createContactWindow,
	updateContactWindow updateContactWindow,
	deleteContactWindow deleteContactWindow,
	aboutWindow aboutWindow,
) *Builder {
	return &Builder{
		app:                 app,
		contactList:         contactList,
		createContactWindow: createContactWindow,
		updateContactWindow: updateContactWindow,
		deleteContactWindow: deleteContactWindow,
		aboutWindow:         aboutWindow,
	}
}

func (b *Builder) Build() *fyne.MainMenu {
	// Закрываем приложение
	exit := fyne.NewMenuItem("Exit", func() {
		b.app.Quit()
	})

	file := fyne.NewMenu("File", exit)

	// Создание контакта
	createContact := fyne.NewMenuItem("Add contact", func() {
		window := b.createContactWindow.Build()
		window.Show()
	})
	// Изменение контакта
	updateContact := fyne.NewMenuItem("Edit contact", func() {
		uuid := b.contactList.SelectedContactUUID()

		// Если контакт не выбран, то ничего не делаем
		if uuid == nil {
			return
		}

		window := b.updateContactWindow.Build(*uuid)
		window.Show()
	})
	// Удаление контакта
	deleteContact := fyne.NewMenuItem("Remove contact", func() {
		uuid := b.contactList.SelectedContactUUID()

		// Если контакт не выбран, то ничего не делаем
		if uuid == nil {
			return
		}

		window := b.deleteContactWindow.Build(*uuid)
		window.Show()
	})

	edit := fyne.NewMenu("Edit", createContact, updateContact, deleteContact)

	about := fyne.NewMenuItem("About app", func() {
		window := b.aboutWindow.Build()
		window.Show()
	})

	info := fyne.NewMenu("Info", about)

	return fyne.NewMainMenu(file, edit, info)
}
