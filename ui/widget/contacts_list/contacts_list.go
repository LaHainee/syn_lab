package contacts_list

import (
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"contacts/internal/model"
	"contacts/ui/dto"
	"contacts/ui/presenter/phone"
	widgetContactInfo "contacts/ui/widget/contact_info"
	"contacts/util/pointer"
)

var (
	contactListSize = fyne.NewSize(350, 500)
	contactListPos  = fyne.NewPos(50, 100)
)

type Builder struct {
	contacts []model.Contact

	// Ссылка на контейнер всего приложения
	appBox *fyne.Container

	contactInfoBox  *fyne.Container
	contactsListBox *container.Scroll
	searchInputBox  *fyne.Container
	searchLabelBox  *fyne.Container

	// Выбранный в данный момент контакте
	selectedContact *model.Contact
}

func NewBuilder(contacts []model.Contact, appBox *fyne.Container) *Builder {
	return &Builder{
		contacts: contacts,
		appBox:   appBox,
	}
}

func (b *Builder) ContactListBoxSize() fyne.Size {
	return contactListSize
}

func (b *Builder) ContactListBoxPos() fyne.Position {
	return contactListPos
}

func (b *Builder) Build() {
	filtered := b.contacts

	sort.Sort(BySurname(filtered))

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

		b.selectedContact = &contact

		contactsWidgetRowsData := []dto.ContactInfoWidgetRowData{
			{
				Label: "Surname",
				Entry: dto.ContactInfoWidgetRowEntry{
					Value:       &contact.Surname,
					Type:        dto.ContactWidgetRowTypeText,
					DisableEdit: true,
				},
			},
			{
				Label: "Name",
				Entry: dto.ContactInfoWidgetRowEntry{
					Value:       &contact.Name,
					Type:        dto.ContactWidgetRowTypeText,
					DisableEdit: true,
				},
			},
			{
				Label: "Birthday",
				Entry: dto.ContactInfoWidgetRowEntry{
					Value:       pointer.To(contact.Birthday.Format("02.01.2006")),
					Type:        dto.ContactWidgetRowTypeDatePicker,
					DisableEdit: true,
				},
			},
			{
				Label: "Phone",
				Entry: dto.ContactInfoWidgetRowEntry{
					Value:       pointer.To(phone.Present(contact.Phone.Number())),
					Type:        dto.ContactWidgetRowTypeText,
					DisableEdit: true,
				},
			},
			{
				Label: "Email",
				Entry: dto.ContactInfoWidgetRowEntry{
					Value:       &contact.Email,
					Type:        dto.ContactWidgetRowTypeText,
					DisableEdit: true,
				},
			},
		}
		for link, value := range contact.Links {
			contactsWidgetRowsData = append(contactsWidgetRowsData, dto.ContactInfoWidgetRowData{
				Label: string(link),
				Entry: dto.ContactInfoWidgetRowEntry{
					Value:       &value,
					Type:        dto.ContactWidgetRowTypeText,
					DisableEdit: true,
				},
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
	b.contactsListBox.Resize(contactListSize)
	b.contactsListBox.Move(contactListPos)

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

	sort.Sort(BySurname(filtered))

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

func (b *Builder) SelectedContactUUID() *string {
	if b.selectedContact == nil {
		return nil
	}

	return &b.selectedContact.UUID
}

func (b *Builder) Refresh(contacts []model.Contact) {
	b.contacts = contacts

	// Удаляем предыдущее наполнение
	b.appBox.Remove(b.contactsListBox)
	b.appBox.Remove(b.contactInfoBox)

	b.Build()

	b.appBox.Refresh()
}
