package main

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	contactValidator "contacts/internal/domain/validate/contact"
	"contacts/internal/handler/create"
	"contacts/internal/handler/update"
	"contacts/internal/model"
	"contacts/internal/storage"
	widgetContactsList "contacts/ui/widget/contacts_list"
	widgetCreateContact "contacts/ui/widget/create_contact"
	widgetDeleteContact "contacts/ui/widget/delete_contact"
	widgetUpdateContact "contacts/ui/widget/update_contact"
)

const (
	windowWidth  = 1920
	windowHeight = 1080
)

func main() {
	// Конфигурация приложения
	contactStorage := storage.New("internal/database/database.json")

	validator := contactValidator.New()

	createContactHandler := create.NewHandler(contactStorage, validator)
	updateContactHandler := update.NewHandler(contactStorage, validator)

	// Создание нового приложения
	myApp := app.New()
	myApp.Settings().SetTheme(theme.LightTheme())
	myWindow := myApp.NewWindow("Contacts App")
	myWindow.Resize(fyne.NewSize(windowWidth, windowHeight))
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

	contactsListWidgetBuilder := widgetContactsList.NewBuilder(contacts, appBox)
	contactsListWidgetBuilder.Build()

	// Компонент отвечающий за создание контакта
	createContactWindowBuilder := widgetCreateContact.NewBuilder(myApp, contactsListWidgetBuilder, createContactHandler, contactStorage)
	createContactButton := widget.NewButtonWithIcon("", createContactIcon, func() {
		createContactWindow := createContactWindowBuilder.Build()
		createContactWindow.Show()
	})
	createContactButton.Resize(fyne.NewSize(30, 30))
	createContactButton.Move(fyne.NewPos(50, 800))

	// Компонент отвечающий за изменение контакта
	updateContactWindowBuilder := widgetUpdateContact.NewBuilder(myApp, contactsListWidgetBuilder, updateContactHandler, contactStorage)
	updateContactButton := widget.NewButtonWithIcon("", editContactIcon, func() {
		selectedContactUUID := contactsListWidgetBuilder.SelectedContactUUID()
		if selectedContactUUID == nil {
			return
		}

		updateContactWindow := updateContactWindowBuilder.Build(*selectedContactUUID)
		updateContactWindow.Show()
	})
	updateContactButton.Resize(fyne.NewSize(30, 30))
	updateContactButton.Move(fyne.NewPos(85, 800))

	// Компонент отвечающий за удаление контакта
	deleteContactWindowBuilder := widgetDeleteContact.NewBuilder(myApp, contactStorage, contactsListWidgetBuilder)
	deleteContactButton := widget.NewButtonWithIcon("", deleteContactIcon, func() {
		selectedContactUUID := contactsListWidgetBuilder.SelectedContactUUID()
		if selectedContactUUID == nil {
			return
		}

		deleteContactWindow := deleteContactWindowBuilder.Build(*selectedContactUUID)
		deleteContactWindow.Show()
	})
	deleteContactButton.Resize(fyne.NewSize(30, 30))
	deleteContactButton.Move(fyne.NewPos(120, 800))

	appBox.Add(createContactButton)
	appBox.Add(updateContactButton)
	appBox.Add(deleteContactButton)

	myWindow.SetContent(appBox)

	// Показать окно и запустить приложение
	myWindow.ShowAndRun()
}

type contactInfo struct {
	Label string
	Value string
}

type createContactField struct {
	Label     *widget.Label
	Entry     *widget.Entry
	ApplyFunc func(contact *model.Contact)
}

func prepareContactsInfo(infos []contactInfo, posByY, offsetByY, labelPosByX, entryPosByX float32) *fyne.Container {
	box := container.NewWithoutLayout()

	for _, info := range infos {
		label := widget.NewLabel(info.Label + ":")
		label.Alignment = fyne.TextAlignTrailing

		labelBox := container.NewVBox(label)
		labelBox.Resize(fyne.NewSize(200, 30))
		labelBox.Move(fyne.NewPos(labelPosByX, posByY))

		entry := widget.NewEntry()
		entry.SetText(info.Value)
		entryBox := container.NewVBox(entry)
		entryBox.Resize(fyne.NewSize(200, 30))
		entryBox.Move(fyne.NewPos(entryPosByX, posByY))

		box.Add(labelBox)
		box.Add(entryBox)

		posByY += offsetByY
	}

	return box
}

func formatPhoneNumber(phone int64) string {
	phoneStr := strconv.FormatInt(phone, 10)

	if len(phoneStr) != 11 {
		return "Неверный номер телефона"
	}

	formatted := fmt.Sprintf("+7 (%s) %s-%s-%s",
		phoneStr[1:4],  // 915
		phoneStr[4:7],  // 159
		phoneStr[7:9],  // 67
		phoneStr[9:11], // 81
	)

	return formatted
}
