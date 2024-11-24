package contacts_list

import (
	"contacts/internal/model"
	"contacts/internal/ui/dto"
	widgetContactInfo "contacts/internal/ui/widget/contact_info"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"strconv"
	"strings"
)

type Builder struct {
	contacts []model.Contact

	// Ссылка на контейнер всего приложения
	appBox *fyne.Container

	contactInfoBox  *fyne.Container
	contactsListBox *container.Scroll
	searchInputBox  *fyne.Container
	searchLabelBox  *fyne.Container
}

func NewBuilder(contacts []model.Contact, appBox *fyne.Container) *Builder {
	return &Builder{
		contacts: contacts,
		appBox:   appBox,
	}
}

func (b *Builder) Build() {
	filtered := b.contacts

	// Список контактов
	contactsList := widget.NewList(
		func() int {
			return len(filtered) // Количество строк в списке
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("") // Создание элемента списка
		},
		func(id int, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(filtered[id].Surname) // Установка текста для элемента
		},
	)

	contactsList.OnSelected = func(id int) {
		contact := filtered[id]

		contactsWidgetRowsData := []dto.ContactInfoWidgetRowData{
			{
				Label: "Surname",
				Value: contact.Surname,
			},
			{
				Label: "Name",
				Value: contact.Name,
			},
			{
				Label: "Birthday",
				Value: contact.Birthday.Format("02.01.2006"),
			},
			{
				Label: "Phone",
				Value: presentPhone(contact.Phone.Number()),
			},
			{
				Label: "Email",
				Value: contact.Email,
			},
		}
		for link, value := range contact.Links {
			contactsWidgetRowsData = append(contactsWidgetRowsData, dto.ContactInfoWidgetRowData{
				Label: string(link),
				Value: value,
			})
		}

		contactInfoWidgetBuilder := widgetContactInfo.NewBuilder(
			dto.Position{
				X: 300,
				Y: 50,
			},
			50,
		)

		contactInfoWidget := contactInfoWidgetBuilder.Build(contactsWidgetRowsData)

		b.appBox.Remove(b.contactInfoBox)
		b.contactInfoBox = contactInfoWidget.Box
		b.appBox.Add(b.contactInfoBox)
	}

	b.contactsListBox = container.NewVScroll(contactsList)
	b.contactsListBox.Resize(fyne.NewSize(350, 500))
	b.contactsListBox.Move(fyne.NewPos(50, 100))

	b.appBox.Add(b.contactsListBox)

	if b.searchInputBox != nil {
		return
	}

	// Поисковая строка
	searchInput := widget.NewEntry()

	searchInput.OnChanged = func(text string) {
		filtered = []model.Contact{}

		for _, contact := range b.contacts {
			if strings.Contains(strings.ToLower(contact.Surname), strings.ToLower(text)) {
				filtered = append(filtered, contact)
			}
		}

		contactsList.Refresh()
	}

	searchInputBox := container.NewVBox(searchInput)
	searchInputBox.Resize(fyne.NewSize(300, 40))
	searchInputBox.Move(fyne.NewPos(100, 50))

	// Текст для поисковой строки
	searchLabel := widget.NewLabel("Find:")
	searchLabel.Alignment = fyne.TextAlignLeading
	searchLabelBox := container.NewVBox(searchLabel)
	searchLabelBox.Resize(fyne.NewSize(50, 40))
	searchLabelBox.Move(fyne.NewPos(50, 50))

	// Добавляем компоненты в приложение
	b.appBox.Add(searchInputBox)
	b.appBox.Add(searchLabelBox)

	return
}

func (b *Builder) Refresh(contacts []model.Contact) {
	b.contacts = contacts

	// Удаляем предыдущее наполнение
	b.appBox.Remove(b.contactsListBox)
	b.appBox.Remove(b.contactInfoBox)

	b.Build()

	b.appBox.Refresh()
}

func presentPhone(phone int64) string {
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
