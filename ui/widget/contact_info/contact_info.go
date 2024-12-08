package contact_info

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	expWidget "fyne.io/x/fyne/widget"

	"contacts/ui/dto"
)

var (
	labelBoxSize           = fyne.NewSize(200, 30)
	textEntryBoxSize       = fyne.NewSize(400, 30)
	datePickerEntryBoxSize = fyne.NewSize(90, 30)
	datePickerButtonSize   = fyne.NewSize(30, 30)
	calendarSize           = fyne.NewSize(225, 200)
	calendarBackgroundSize = fyne.NewSize(calendarSize.Width+10, calendarSize.Height+10)
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

	// Для того, чтобы связать созданный label и entry
	assignedByLabel := make(map[string]dto.ContactWidgetRow, len(rowsData))

	// Выношу из цикла, т.к. компонент должен быть поверх всех остальных
	var (
		calendars           []*expWidget.Calendar
		calendarBackgrounds []*canvas.Rectangle
	)

	for _, rowData := range rowsData {
		label := widget.NewLabel(rowData.Label + ":")
		label.Alignment = fyne.TextAlignTrailing

		// Контейнер для текста
		labelBox := container.NewVBox(label)
		labelBox.Resize(labelBoxSize)
		labelBox.Move(fyne.NewPos(w.firstRowPosition.X, currentPosY))

		box.Add(labelBox)

		var entry *widget.Entry

		switch rowData.Entry.Type {
		case dto.ContactWidgetRowTypeText:
			entry = w.buildEntry(rowData.Entry)

			entryBox := container.NewVBox(entry)
			entryBox.Resize(textEntryBoxSize)
			entryBox.Move(fyne.NewPos(w.firstRowPosition.X+labelBoxSize.Width, currentPosY))

			box.Add(entryBox)
		case dto.ContactWidgetRowTypeDatePicker:
			entry = w.buildEntry(rowData.Entry)

			entryBox := container.NewVBox(entry)
			entryBox.Resize(datePickerEntryBoxSize)
			entryBox.Move(fyne.NewPos(w.firstRowPosition.X+labelBoxSize.Width, currentPosY))

			var (
				calendar           *expWidget.Calendar
				calendarBackground *canvas.Rectangle
			)

			currentTime := time.Now()
			if rowData.Entry.Value != nil {
				t, err := time.Parse("02.01.2006", *rowData.Entry.Value)
				if err == nil {
					currentTime = t
				}
			}

			calendar = expWidget.NewCalendar(currentTime, func(t time.Time) {
				entry.SetText(t.Format("02.01.2006"))
				// После выбора даты скрываем календарь и подложку
				calendar.Hide()
				calendarBackground.Hide()
			})
			calendar.Resize(calendarSize)
			calendar.Hide()

			// Достаем иконку календаря
			datePickerIcon, err := fyne.LoadResourceFromPath("./ui/icons/datepicker.png")
			if err != nil {
				panic(err)
			}

			// Создаем кнопку для открытия календаря
			datePickerButton := widget.NewButtonWithIcon("", datePickerIcon, func() {
				if !calendar.Hidden {
					calendar.Hide()
					calendarBackground.Hide()
					return
				}
				calendar.Show()
				calendarBackground.Show()
			})
			datePickerButton.Resize(datePickerButtonSize)
			datePickButtonPos := fyne.NewPos(
				w.firstRowPosition.X+labelBoxSize.Width+datePickerEntryBoxSize.Width+5,
				currentPosY+2,
			)
			datePickerButton.Move(datePickButtonPos)

			// Подложка под календарь
			calendarBackground = canvas.NewRectangle(color.White)
			calendarBackground.CornerRadius = 10
			calendarBackground.Resize(calendarBackgroundSize)
			calendarBackground.Hide() // По умолчанию скрыт

			// Позиционирование календаря
			calendarPos := fyne.NewPos(
				datePickerButton.Position().X+50,
				currentPosY-2*w.spacingBetweenRows+10,
			)
			calendar.Move(fyne.NewPos(calendarPos.X+5, calendarPos.Y))
			calendarBackground.Move(calendarPos)

			box.Add(entryBox)
			box.Add(datePickerButton)

			// Т.к. календарей может быть несколько, то добавляем их в массив
			calendars = append(calendars, calendar)
			calendarBackgrounds = append(calendarBackgrounds, calendarBackground)
		}

		currentPosY += w.spacingBetweenRows

		assignedByLabel[rowData.Label] = dto.ContactWidgetRow{
			Label: label,
			Entry: entry,
		}
	}

	// Вынесено сюда, чтобы компоненты были поверх полей ввода
	for i := 0; i < len(calendars); i++ {
		box.Add(calendarBackgrounds[i])
		box.Add(calendars[i])
	}

	return dto.ContactInfoWidget{
		AssignedByLabel: assignedByLabel,
		Box:             box,
		Size: &fyne.Size{
			Width:  labelBoxSize.Width + textEntryBoxSize.Width,
			Height: currentPosY,
		},
	}
}

func (w *Builder) buildEntry(entryDto dto.ContactInfoWidgetRowEntry) *widget.Entry {
	entry := widget.NewEntry()

	// Подставляем текст в форму
	if entryDto.Value != nil {
		entry.SetText(*entryDto.Value)
	}

	// Шаблон для заполнения
	if entryDto.Placeholder != nil {
		entry.SetPlaceHolder(*entryDto.Placeholder)
	}

	// Для опции DisableEdit запрещаем редактировать форму
	if entryDto.DisableEdit {
		text := entry.Text
		entry.OnChanged = func(_ string) {
			entry.SetText(text)
		}
	}

	return entry
}
