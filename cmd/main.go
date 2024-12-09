package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	contactValidator "contacts/internal/domain/validate/contact"
	createContact "contacts/internal/handler/create"
	deleteContact "contacts/internal/handler/delete"
	fetchContact "contacts/internal/handler/fetch"
	searchContact "contacts/internal/handler/search"
	updateContact "contacts/internal/handler/update"
	"contacts/internal/storage"
	"contacts/internal/storage/database"
	"contacts/ui/menu"
	widgetBirthday "contacts/ui/widget/birthday"
	widgetContactsList "contacts/ui/widget/contacts_list"
	windowAbout "contacts/ui/window/about"
	windowCreateContact "contacts/ui/window/create_contact"
	windowDeleteContact "contacts/ui/window/delete_contact"
	windowUpdateContact "contacts/ui/window/update_contact"
	"contacts/util/uuid"
)

var (
	appWindowSize = fyne.NewSize(1920, 1080)
	buttonSize    = fyne.NewSize(30, 30)
)

func main() {
	// Конфигурация приложения
	contactStorage := storage.New(database.New("internal/database/database.json"))

	validator := contactValidator.New()

	uuidGenerator := uuid.NewGenerator()

	createContactHandler := createContact.NewHandler(contactStorage, uuidGenerator, validator)
	updateContactHandler := updateContact.NewHandler(contactStorage, validator)
	deleteContactHandler := deleteContact.NewHandler(contactStorage)
	fetchContactHandler := fetchContact.NewHandler(contactStorage)
	searchContactHandler := searchContact.NewHandler(contactStorage)

	// Создание нового приложения
	myApp := app.New()
	myApp.Quit()
	myApp.Settings().SetTheme(theme.LightTheme())
	myWindow := myApp.NewWindow("Contacts App")
	myWindow.Resize(appWindowSize)
	appBox := container.NewWithoutLayout()

	contacts, err := contactStorage.Fetch()
	if err != nil {
		panic(err)
	}

	// <! Иконки для кнопок
	createContactIcon, err := fyne.LoadResourceFromPath("./ui/icons/plus.png")
	if err != nil {
		panic(err)
	}

	editContactIcon, err := fyne.LoadResourceFromPath("./ui/icons/edit.png")
	if err != nil {
		panic(err)
	}

	deleteContactIcon, err := fyne.LoadResourceFromPath("./ui/icons/minus.png")
	if err != nil {
		panic(err)
	}
	// Иконки для кнопок !>

	contactsListWidgetBuilder := widgetContactsList.NewBuilder(fetchContactHandler, searchContactHandler, appBox)
	contactsListWidgetBuilder.Build()

	contactListPos := contactsListWidgetBuilder.ContactListBoxPos()
	contactListSize := contactsListWidgetBuilder.ContactListBoxSize()

	const horizontalSpacingBetweenButtons = 5

	// Компонент отвечающий за создание контакта
	createContactWindowBuilder := windowCreateContact.NewBuilder(
		myApp,
		contactsListWidgetBuilder,
		createContactHandler,
	)
	createContactButton := widget.NewButtonWithIcon("", createContactIcon, func() {
		createContactWindow := createContactWindowBuilder.Build()
		createContactWindow.Show()
	})
	createContactButton.Resize(buttonSize)
	createContactButton.Move(fyne.NewPos(contactListPos.X, contactListPos.Y+contactListSize.Height+20))

	// Компонент отвечающий за изменение контакта
	updateContactWindowBuilder := windowUpdateContact.NewBuilder(
		myApp,
		contactsListWidgetBuilder,
		updateContactHandler,
		fetchContactHandler,
	)
	updateContactButton := widget.NewButtonWithIcon("", editContactIcon, func() {
		selectedContactUUID := contactsListWidgetBuilder.SelectedContactUUID()
		if selectedContactUUID == nil {
			return
		}

		updateContactWindow := updateContactWindowBuilder.Build(*selectedContactUUID)
		updateContactWindow.Show()
	})
	updateContactButton.Resize(buttonSize)
	updateContactButton.Move(
		fyne.NewPos(
			createContactButton.Position().X+createContactButton.Size().Width+horizontalSpacingBetweenButtons,
			createContactButton.Position().Y,
		))

	// Компонент отвечающий за удаление контакта
	deleteContactWindowBuilder := windowDeleteContact.NewBuilder(
		myApp,
		deleteContactHandler,
		fetchContactHandler,
		contactsListWidgetBuilder,
	)
	deleteContactButton := widget.NewButtonWithIcon("", deleteContactIcon, func() {
		selectedContactUUID := contactsListWidgetBuilder.SelectedContactUUID()
		if selectedContactUUID == nil {
			return
		}

		deleteContactWindow := deleteContactWindowBuilder.Build(*selectedContactUUID)
		deleteContactWindow.Show()
	})
	deleteContactButton.Resize(buttonSize)
	deleteContactButton.Move(
		fyne.NewPos(
			updateContactButton.Position().X+updateContactButton.Size().Width+horizontalSpacingBetweenButtons,
			updateContactButton.Position().Y,
		))

	aboutWindowBuilder := windowAbout.NewBuilder(myApp)

	// Виджет с напоминанием о днях рождения
	birthdayWidgetBuilder := widgetBirthday.NewBuilder(appWindowSize)
	birthdayWidget := birthdayWidgetBuilder.Build(contacts)
	if birthdayWidget != nil {
		birthdayWidget.Move(fyne.NewPos(contactListPos.X+contactListSize.Width+20, contactListPos.Y+contactListSize.Height-50))
		appBox.Add(birthdayWidget)
	}

	appBox.Add(createContactButton)
	appBox.Add(updateContactButton)
	appBox.Add(deleteContactButton)
	myWindow.SetContent(appBox)

	mainMenuBuilder := menu.NewBuilder(
		myApp,
		contactsListWidgetBuilder,
		createContactWindowBuilder,
		updateContactWindowBuilder,
		deleteContactWindowBuilder,
		aboutWindowBuilder,
	)
	myWindow.SetMainMenu(mainMenuBuilder.Build())

	myWindow.ShowAndRun()
}
