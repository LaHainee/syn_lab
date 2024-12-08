package update_contact

import (
	"context"
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"contacts/internal/model"
	"contacts/ui/dto"
	"contacts/ui/presenter/phone"
	wigetContactInfo "contacts/ui/widget/contact_info"
	contactsList "contacts/ui/widget/contacts_list"
	errorWidget "contacts/ui/widget/error"
	"contacts/util/pointer"
)

type Builder struct {
	// Объект приложения
	app fyne.App

	// Компонент, который отвечает за список контактов. Нужен для обновления списка после создания
	contactListBuilder *contactsList.Builder

	handler handler
	storage storage
}

func NewBuilder(
	app fyne.App,
	contactsListBuilder *contactsList.Builder,
	handler handler,
	storage storage,
) *Builder {
	return &Builder{
		app:                app,
		contactListBuilder: contactsListBuilder,
		handler:            handler,
		storage:            storage,
	}
}

func (b *Builder) Build(contactUuid string) fyne.Window {
	contact, err := b.storage.FetchByUUID(contactUuid)
	if err != nil {
		panic(err)
	}

	contactInfoWidgetRowsData := []dto.ContactInfoWidgetRowData{
		{
			Label: "Surname",
			Entry: dto.ContactInfoWidgetRowEntry{
				Type:  dto.ContactWidgetRowTypeText,
				Value: &contact.Surname,
			},
		},
		{
			Label: "Name",
			Entry: dto.ContactInfoWidgetRowEntry{
				Type:  dto.ContactWidgetRowTypeText,
				Value: &contact.Name,
			},
		},
		{
			Label: "Birthday",
			Entry: dto.ContactInfoWidgetRowEntry{
				Type:  dto.ContactWidgetRowTypeDatePicker,
				Value: pointer.To(contact.Birthday.Format("02.01.2006")),
			},
		},
		{
			Label: "Phone",
			Entry: dto.ContactInfoWidgetRowEntry{
				Type:  dto.ContactWidgetRowTypeText,
				Value: pointer.To(phone.Present(contact.Phone.Number())),
			},
		},
		{
			Label: "Email",
			Entry: dto.ContactInfoWidgetRowEntry{
				Type:  dto.ContactWidgetRowTypeText,
				Value: &contact.Email,
			},
		},
	}

	for link, value := range contact.Links {
		contactInfoWidgetRowsData = append(contactInfoWidgetRowsData, dto.ContactInfoWidgetRowData{
			Label: string(link),
			Entry: dto.ContactInfoWidgetRowEntry{
				Type:  dto.ContactWidgetRowTypeText,
				Value: &value,
			},
		})
	}

	contactInfoWidgetBuilder := wigetContactInfo.NewBuilder(
		dto.Position{
			X: -100,
			Y: 25,
		},
		50,
	)

	contactInfoWidget := contactInfoWidgetBuilder.Build(contactInfoWidgetRowsData)

	window := b.app.NewWindow("Изменить контакт")
	window.Resize(fyne.NewSize(contactInfoWidget.Size.Width-75, contactInfoWidget.Size.Height+100))
	window.CenterOnScreen()
	window.SetFixedSize(true)

	// Форма для отображения текста об ошибке
	errorLabel := widget.NewLabel("")
	errorLabel.Resize(fyne.NewSize(contactInfoWidget.Size.Width-50, 50))
	errorLabel.Move(fyne.NewPos(20, contactInfoWidget.Size.Height-10))
	errorLabel.Hide()

	closeButton := widget.NewButton("Cancel", func() {
		window.Close()
	})
	closeButton.Resize(fyne.NewSize(70, 30))
	closeButton.Move(fyne.NewPos(contactInfoWidget.Size.Width-100-closeButton.Size().Width, contactInfoWidget.Size.Height+50))

	confirmButton := widget.NewButton("OK", func() {
		// Очистим предыдущий стейт:
		// 1. Скроем сообщения об ошибке
		// 2. Перекрасим лейблы в черный цвет
		errorLabel.Hide()
		for _, contactInfoWidgetRow := range contactInfoWidget.AssignedByLabel {
			contactInfoWidgetRow.Label.Importance = widget.MediumImportance
			contactInfoWidgetRow.Label.Refresh()
		}

		links := make(map[model.ContactLink]string)
		for link := range contact.Links {
			contactWidgetRow, ok := contactInfoWidget.AssignedByLabel[string(link)]
			if !ok {
				continue
			}
			links[link] = contactWidgetRow.Entry.Text
		}

		fieldMsgs, err := b.handler.Update(context.Background(), model.ContactForCreate{
			UUID:     &contact.UUID,
			Surname:  contactInfoWidget.AssignedByLabel["Surname"].Entry.Text,
			Name:     contactInfoWidget.AssignedByLabel["Name"].Entry.Text,
			Birthday: contactInfoWidget.AssignedByLabel["Birthday"].Entry.Text,
			Phone:    contactInfoWidget.AssignedByLabel["Phone"].Entry.Text,
			Email:    contactInfoWidget.AssignedByLabel["Email"].Entry.Text,
			Links:    links,
		})
		if err != nil {
			if errors.Is(err, model.ErrValidation) {
				errorWidget.Show(fieldMsgs, &contactInfoWidget, errorLabel)
				return
			}

			panic(err)
		}

		contacts, err := b.storage.Fetch()
		if err != nil {
			panic(err)
		}

		b.contactListBuilder.Refresh(contacts)

		window.Close()
	})
	confirmButton.Resize(fyne.NewSize(70, 30))
	confirmButton.Move(fyne.NewPos(closeButton.Position().X-25-confirmButton.Size().Width, contactInfoWidget.Size.Height+50))

	box := container.NewWithoutLayout()
	box.Add(contactInfoWidget.Box)
	box.Add(closeButton)
	box.Add(confirmButton)
	box.Add(errorLabel)

	window.SetContent(box)

	return window
}
