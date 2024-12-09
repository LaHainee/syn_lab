package about

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

const (
	versionText = "v. 1.0.0"
)

var (
	windowSize = fyne.NewSize(500, 300)
)

type Builder struct {
	app app
}

func NewBuilder(app app) *Builder {
	return &Builder{
		app: app,
	}
}

func (b *Builder) Build() fyne.Window {
	// Конфигурация нового окна
	window := b.app.NewWindow("About")
	window.Resize(windowSize)
	window.SetFixedSize(true)
	window.CenterOnScreen()

	// Компоненты окна
	title := canvas.NewText("ContactsApp", color.Black)
	title.TextSize = 24
	title.TextStyle = fyne.TextStyle{Bold: true}

	version := canvas.NewText(versionText, color.Black)
	version.TextSize = 14
	version.Move(fyne.NewPos(0, title.Position().Y+title.TextSize+10))

	links := buildLinks()
	links.Move(fyne.NewPos(0, version.Position().Y+version.TextSize+25))

	footer := canvas.NewText("2024 Ershov V.A.", color.Black)
	footer.TextSize = 14
	footer.Move(fyne.NewPos(0, windowSize.Height-30))

	box := container.NewWithoutLayout(
		title,
		version,
		links,
		footer,
	)

	window.SetContent(box)

	return window
}

// buildLinks – формирует строки вида Author: Ershov Vitaliy
//
// y – с какой позиции по y начинать список
// spacing – расстояние между элементами списка
func buildLinks() *fyne.Container {
	type link struct {
		setRedirectStyle bool
		url              string
	}

	type row struct {
		label string
		link  link
	}

	rows := []row{
		{
			label: "Author:",
			link: link{
				url: "Ershov V.A.",
			},
		},
		// Пустая строка
		{},
		{
			label: "email for feedback:",
			link: link{
				setRedirectStyle: true,
				url:              "vaershov@avito.ru",
			},
		},
		{
			label: "github:",
			link: link{
				setRedirectStyle: true,
				url:              "https://github.com/LaHainee/syn_lab",
			},
		},
	}

	vertical := container.NewVBox()

	for _, r := range rows {
		labelText := canvas.NewText(r.label, color.Black)
		labelText.TextSize = 14

		linkText := canvas.NewText(r.link.url, color.Black)
		linkText.TextSize = 14

		linkBox := container.NewWithoutLayout()
		linkBox.Add(linkText)

		// Если ссылку требуется стилизовать
		if r.link.setRedirectStyle {
			linkText.Color = color.RGBA{B: 255, A: 255}

			// В fyne не поддерживается подчеркнутый текст, поэтому костылим свою линию
			line := canvas.NewLine(color.RGBA{B: 255, A: 255})
			lineY := linkText.TextSize + 5
			line.Position1 = fyne.NewPos(0, lineY)
			line.Position2 = fyne.NewPos(linkBox.MinSize().Width, lineY)
			linkBox.Add(line)
		}

		horizontal := container.NewHBox(labelText, linkBox)

		vertical.Add(horizontal)
	}

	return vertical
}
