package create_contact

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"

	contactsDomain "contacts/internal/domain/contacts"
	"contacts/internal/model"
	"contacts/internal/storage"
	"contacts/ui/dto"
	wigetContactInfo "contacts/ui/widget/contact_info"
	contactsList "contacts/ui/widget/contacts_list"
	"contacts/util/pointer"
)

var allowedLinks = contactsDomain.AllowedLinks()

type Builder struct {
	// Объект приложения
	app fyne.App

	storage *storage.Storage

	contactListBuilder *contactsList.Builder
}

func NewBuilder(app fyne.App, storage *storage.Storage, contactsListBuilder *contactsList.Builder) *Builder {
	return &Builder{
		app:                app,
		storage:            storage,
		contactListBuilder: contactsListBuilder,
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
	window.Resize(fyne.NewSize(contactInfoWidget.Size.Width-75, contactInfoWidget.Size.Height+50))
	window.CenterOnScreen()
	window.SetFixedSize(true)

	closeButton := widget.NewButton("Cancel", func() {
		window.Close()
	})
	closeButton.Resize(fyne.NewSize(70, 30))
	closeButton.Move(fyne.NewPos(contactInfoWidget.Size.Width-100-closeButton.Size().Width, contactInfoWidget.Size.Height))

	confirmButton := widget.NewButton("OK", func() {
		links := make(map[model.ContactLink]string)

		for _, link := range allowedLinks {
			contactWidgetRow, ok := contactInfoWidget.AssignedByLabel[string(link)]
			if !ok {
				continue
			}

			links[link] = contactWidgetRow.Entry.Text
		}

		birthday, err := time.Parse("02.01.2006", contactInfoWidget.AssignedByLabel["Birthday"].Entry.Text)
		if err != nil {
			panic(err)
		}

		phone, err := model.NewPhone(contactInfoWidget.AssignedByLabel["Phone"].Entry.Text)
		if err != nil {
			panic(err)
		}

		err = b.storage.Create(model.Contact{
			UUID:     uuid.NewString(),
			Surname:  contactInfoWidget.AssignedByLabel["Surname"].Entry.Text,
			Name:     contactInfoWidget.AssignedByLabel["Name"].Entry.Text,
			Birthday: birthday,
			Phone:    phone,
			Email:    contactInfoWidget.AssignedByLabel["Email"].Entry.Text,
			Links:    links,
		})
		if err != nil {
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
	confirmButton.Move(fyne.NewPos(closeButton.Position().X-25-confirmButton.Size().Width, contactInfoWidget.Size.Height))

	box := container.NewWithoutLayout()
	box.Add(contactInfoWidget.Box)
	box.Add(closeButton)
	box.Add(confirmButton)

	window.SetContent(box)

	return window
}
