package error

import (
	"fyne.io/fyne/v2/widget"

	contactsDomain "contacts/internal/domain/contacts"
	"contacts/internal/model"
	"contacts/ui/dto"
)

var allowedLinks = contactsDomain.AllowedLinks()

func Show(fieldMsgs map[model.Field]string, contactInfoWidget *dto.ContactInfoWidget, errorLabel *widget.Label) {
	var messageToShow *string

	// Базовые поля
	for field, message := range fieldMsgs {
		messageToShow = &message

		switch field {
		case model.FieldName:
			contactInfoWidget.AssignedByLabel["Name"].Label.Importance = widget.DangerImportance
		case model.FieldSurname:
			contactInfoWidget.AssignedByLabel["Surname"].Label.Importance = widget.DangerImportance
		case model.FieldBirthday:
			contactInfoWidget.AssignedByLabel["Birthday"].Label.Importance = widget.DangerImportance
		case model.FieldEmail:
			contactInfoWidget.AssignedByLabel["Email"].Label.Importance = widget.DangerImportance
		case model.FieldPhone:
			contactInfoWidget.AssignedByLabel["Phone"].Label.Importance = widget.DangerImportance
		}
	}

	// Ссылки
	for _, link := range allowedLinks {
		contactWidgetRow, ok := contactInfoWidget.AssignedByLabel[string(link)]
		if !ok {
			continue
		}

		msg, ok := fieldMsgs[model.Field(link)]
		if !ok {
			continue
		}

		messageToShow = &msg
		contactWidgetRow.Label.Importance = widget.DangerImportance
	}

	// Обновим стейт всех лейблов
	for _, contactInfoWidgetRow := range contactInfoWidget.AssignedByLabel {
		contactInfoWidgetRow.Label.Refresh()
	}

	if messageToShow == nil {
		return
	}

	errorLabel.SetText(*messageToShow)
	errorLabel.Show()
}
