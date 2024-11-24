package contact_info

import (
	"contacts/internal/ui/dto"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var (
	labelBoxSize = fyne.NewSize(200, 30)
	entryBoxSize = fyne.NewSize(400, 30)
)

type Builder struct {
	firstRowPosition   dto.Position
	spacingBetweenRows float32
}

func NewBuilder(
	firstRowPosition dto.Position,
	spacingBetweenRows float32,
) *Builder {
	return &Builder{
		firstRowPosition:   firstRowPosition,
		spacingBetweenRows: spacingBetweenRows,
	}
}

func (w *Builder) Build(rowsData []dto.ContactInfoWidgetRowData) dto.ContactInfoWidget {
	box := container.NewWithoutLayout()

	currentPosY := w.firstRowPosition.Y

	// Для того, чтобы привязаться созданный label и entry
	assignedByLabel := make(map[string]dto.ContactWidgetRow, len(rowsData))

	for _, rowData := range rowsData {
		label := widget.NewLabel(rowData.Label + ":")
		label.Alignment = fyne.TextAlignTrailing

		entry := widget.NewEntry()
		entry.SetText(rowData.Value)

		labelBox := container.NewVBox(label)
		labelBox.Resize(labelBoxSize)
		labelBox.Move(fyne.NewPos(w.firstRowPosition.X, currentPosY))

		entryBox := container.NewVBox(entry)
		entryBox.Resize(entryBoxSize)
		entryBox.Move(fyne.NewPos(w.firstRowPosition.X+labelBoxSize.Width, currentPosY))

		box.Add(labelBox)
		box.Add(entryBox)

		currentPosY += w.spacingBetweenRows

		assignedByLabel[rowData.Label] = dto.ContactWidgetRow{
			Label: label,
			Entry: entry,
		}
	}

	return dto.ContactInfoWidget{
		AssignedByLabel: assignedByLabel,
		Box:             box,
		Size: &fyne.Size{
			Width:  labelBoxSize.Width + entryBoxSize.Width,
			Height: currentPosY,
		},
	}
}
