package create_contact

import (
	"context"
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	contactsDomain "contacts/internal/domain/contacts"
	"contacts/internal/model"
	"contacts/ui/dto"
	wigetContactInfo "contacts/ui/widget/contact_info"
	contactsList "contacts/ui/widget/contacts_list"
	errorWidget "contacts/ui/widget/error"
	"contacts/util/pointer"
)

var allowedLinks = contactsDomain.AllowedLinks()

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

func (b *Builder) Build() fyne.Window {
	contactInfoWidgetRowsData := []dto.ContactInfoWidgetRowData{
		{
			Label: "Surname",
			Entry: dto.ContactInfoWidgetRowEntry{
				Type:        dto.ContactWidgetRowTypeText,
				Placeholder: pointer.To("Ершов"),
			},
		},
		{
			Label: "Name",
			Entry: dto.ContactInfoWidgetRowEntry{
				Type:        dto.ContactWidgetRowTypeText,
				Placeholder: pointer.To("Виталий"),
			},
		},
		{
			Label: "Birthday",
			Entry: dto.ContactInfoWidgetRowEntry{
				Type:        dto.ContactWidgetRowTypeDatePicker,
				Placeholder: pointer.To("10.01.2001"),
			},
		},
		{
			Label: "Phone",
			Entry: dto.ContactInfoWidgetRowEntry{
				Type:        dto.ContactWidgetRowTypeText,
				Placeholder: pointer.To("+7 (915) 159-67-81"),
			},
		},
		{
			Label: "Email",
			Entry: dto.ContactInfoWidgetRowEntry{
				Type:        dto.ContactWidgetRowTypeText,
				Placeholder: pointer.To("vaershov@avito.ru"),
			},
		},
	}

	for _, allowedLink := range allowedLinks {
		contactInfoWidgetRowsData = append(contactInfoWidgetRowsData, dto.ContactInfoWidgetRowData{
			Label: string(allowedLink),
			Entry: dto.ContactInfoWidgetRowEntry{
				Type:        dto.ContactWidgetRowTypeText,
				Placeholder: pointer.To("https://ya.ru"),
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

	window := b.app.NewWindow("Добавить контакт")
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
		for _, link := range allowedLinks {
			contactWidgetRow, ok := contactInfoWidget.AssignedByLabel[string(link)]
			if !ok {
				continue
			}
			links[link] = contactWidgetRow.Entry.Text
		}

		fieldMsgs, err := b.handler.Create(context.Background(), model.ContactForCreate{
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
