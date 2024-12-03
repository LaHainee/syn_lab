package update_contact

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"contacts/internal/model"
	"contacts/internal/storage"
	"contacts/ui/dto"
	"contacts/ui/presenter/phone"
	wigetContactInfo "contacts/ui/widget/contact_info"
	contactsList "contacts/ui/widget/contacts_list"
	"contacts/util/pointer"
)

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

		for link := range contact.Links {
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

		rawPhone, err := model.NewPhone(contactInfoWidget.AssignedByLabel["Phone"].Entry.Text)
		if err != nil {
			panic(err)
		}

		err = b.storage.Update(model.Contact{
			UUID:     contact.UUID,
			Surname:  contactInfoWidget.AssignedByLabel["Surname"].Entry.Text,
			Name:     contactInfoWidget.AssignedByLabel["Name"].Entry.Text,
			Birthday: birthday,
			Phone:    rawPhone,
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
