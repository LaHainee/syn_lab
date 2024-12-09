package dto

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type ContactWidgetRowType string

const (
	ContactWidgetRowTypeDatePicker = "date_picker"
	ContactWidgetRowTypeText       = "text"
)

type ContactInfoWidgetRowData struct {
	Label string
	Entry ContactInfoWidgetRowEntry
}

type ContactInfoWidgetRowEntry struct {
	Value       *string
	Placeholder *string
	Type        ContactWidgetRowType
	DisableEdit bool
}

type ContactInfoWidget struct {
	AssignedByLabel map[string]ContactWidgetRow
	Box             *fyne.Container
	Size            *fyne.Size
}

type ContactWidgetRow struct {
	Label *widget.Label
	Entry *widget.Entry
}
