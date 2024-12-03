package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Calendar struct {
	selectedDate time.Time
	label        *widget.Label
}

func NewCalendar() *Calendar {
	cal := &Calendar{
		selectedDate: time.Now(),
		label:        widget.NewLabel("Выберите дату:"),
	}

	cal.updateLabel()

	grid := container.NewGridWithColumns(7)
	days := []string{"Вс", "Пн", "Вт", "Ср", "Чт", "Пт", "Сб"}
	for _, day := range days {
		grid.Add(widget.NewLabel(day))
	}

	cal.addDaysToGrid(grid)

	return &Calendar{
		selectedDate: cal.selectedDate,
		label:        cal.label,
	}
}

func (c *Calendar) addDaysToGrid(grid *fyne.Container) *fyne.Container {
	year, month, _ := c.selectedDate.Date()
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	daysInMonth := firstDay.AddDate(0, 1, -1).Day()
	startOffset := int(firstDay.Weekday())

	for i := 0; i < startOffset; i++ {
		grid.Add(widget.NewLabel(""))
	}

	for day := 1; day <= daysInMonth; day++ {
		dayButton := widget.NewButton(fmt.Sprintf("%d", day), func(d int) func() {
			return func() {
				c.selectedDate = time.Date(year, month, d, 0, 0, 0, 0, time.UTC)
				c.updateLabel()
			}
		}(day))
		grid.Add(dayButton)
	}

	return grid
}

func (c *Calendar) updateLabel() {
	c.label.SetText(fmt.Sprintf("Выбранная дата: %s", c.selectedDate.Format("2006-01-02")))
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Date Picker")

	calendar := NewCalendar()

	content := container.NewVBox(calendar.label, calendar.addDaysToGrid(container.NewGridWithColumns(7)))

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(300, 400))
	myWindow.ShowAndRun()
}
