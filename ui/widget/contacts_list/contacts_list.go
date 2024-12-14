package contacts_list

import (
	"context"
	"sort"

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
	fetchHandler  fetchHandler
	searchHandler searchHandler
	appBox        appBox

	// Для хранения стейта
	contactInfoBox  *fyne.Container
	contactsListBox *container.Scroll
	searchInputBox  *fyne.Container
	searchLabelBox  *fyne.Container
	selectedContact *model.Contact
}

func NewBuilder(fetchHandler fetchHandler, searchHandler searchHandler, appBox appBox) *Builder {
	return &Builder{
		appBox:        appBox,
		fetchHandler:  fetchHandler,
		searchHandler: searchHandler,
	}
}

func (b *Builder) ContactListBoxSize() fyne.Size {
	return contactListSize
}

func (b *Builder) ContactListBoxPos() fyne.Position {
	return contactListPos
}

func (b *Builder) Build() {
	filtered, err := b.fetchHandler.Fetch(context.Background())
	if err != nil {
		panic(err)
	}

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

	// Обработка ввода в поисковой строке
	searchInput.OnChanged = func(text string) {
		filtered, err = b.searchHandler.Search(context.Background(), model.SearchRequest{
			Query: text,
		})
		if err != nil {
			panic(err)
		}

		sort.Sort(BySurname(filtered))

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

func (b *Builder) SelectedContactUUID() *string {
	if b.selectedContact == nil {
		return nil
	}

	return &b.selectedContact.UUID
}

func (b *Builder) Refresh() {
	defer b.appBox.Refresh()

	// Удаляем предыдущее наполнение
	b.appBox.Remove(b.contactsListBox)
	b.appBox.Remove(b.contactInfoBox)

	b.Build()
}
