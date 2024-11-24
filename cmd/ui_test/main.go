package main

import (
	"contacts/internal/model"
	"contacts/internal/storage"
	widgetContactsList "contacts/internal/ui/widget/contacts_list"
	widgetCreateContact "contacts/internal/ui/widget/create_contact"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"strconv"
)

const (
	windowWidth  = 1920
	windowHeight = 1080
)

func main() {
	// Создание нового приложения
	myApp := app.New()
	myApp.Settings().SetTheme(theme.LightTheme())
	myWindow := myApp.NewWindow("Contacts App")
	myWindow.Resize(fyne.NewSize(windowWidth, windowHeight))
	appBox := container.NewWithoutLayout()

	storageInstance := storage.New("internal/database/database.json")

	contacts, err := storageInstance.Fetch()
	if err != nil {
		panic(err)
	}

	contactsListWidgetBuilder := widgetContactsList.NewBuilder(contacts, appBox)
	contactsListWidgetBuilder.Build()

	createContactWindowBuilder := widgetCreateContact.NewBuilder(myApp, storageInstance, contactsListWidgetBuilder)

	createContactButton := widget.NewButton("+", func() {
		createContactWindow := createContactWindowBuilder.Build()
		createContactWindow.Show()
	})
	createContactButton.Resize(fyne.NewSize(30, 30))
	createContactButton.Move(fyne.NewPos(50, 800))

	appBox.Add(createContactButton)

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
