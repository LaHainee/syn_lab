package birthday

import (
	"image/color"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"

	"contacts/internal/model"
)

var rectangleSize = fyne.NewSize(500, 100)

type Builder struct {
	windowSize fyne.Size
}

func NewBuilder(windowSize fyne.Size) *Builder {
	return &Builder{
		windowSize: windowSize,
	}
}

func (b *Builder) Build(contacts []model.Contact) *fyne.Container {
	birthdayBoysSurnames := make([]string, 0)
	for _, contact := range contacts {
		now := time.Now()

		if contact.Birthday.Month() != now.Month() {
			continue
		}

		if contact.Birthday.Day() != now.Day() {
			continue
		}

		birthdayBoysSurnames = append(birthdayBoysSurnames, contact.Surname)
	}

	if len(birthdayBoysSurnames) == 0 {
		return container.NewWithoutLayout()
	}

	rectangle := canvas.NewRectangle(color.RGBA{R: 186, G: 209, B: 236, A: 255})
	rectangle.CornerRadius = 10
	rectangle.Resize(rectangleSize)

	warningIcon, err := fyne.LoadResourceFromPath("./ui/icons/warning.png")
	if err != nil {
		panic(err)
	}
	warningImage := canvas.NewImageFromResource(warningIcon)
	warningImage.Resize(fyne.NewSize(rectangleSize.Height-20, rectangleSize.Height-20))
	warningImage.Move(fyne.NewPos(20, rectangleSize.Height/2-warningImage.Size().Height/2))

	infoText := canvas.NewText("Сегодня день рождения:", color.Black)
	infoText.TextSize = 16

	birthdayBoysText := canvas.NewText(strings.Join(birthdayBoysSurnames, ", "), color.Black)
	birthdayBoysText.TextSize = 16

	lines := container.NewVBox(infoText, birthdayBoysText)
	lines.Move(fyne.NewPos(warningImage.Position().X+warningImage.Size().Width+10, rectangleSize.Height/2-lines.MinSize().Height/2))

	box := container.NewWithoutLayout(
		rectangle,
		lines,
		warningImage,
	)

	return box
}
