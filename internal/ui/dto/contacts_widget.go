package dto

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type ContactInfoWidgetRowData struct {
	Label string
	Value string
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
