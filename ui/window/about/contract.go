package about

import "fyne.io/fyne/v2"

type app interface {
	NewWindow(title string) fyne.Window
}
