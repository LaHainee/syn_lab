package main

// import fyne
import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

func main() {
	// New App
	a := app.New()
	// New window and title
	w := a.NewWindow("Contact form & Resize widgets")
	// Resize window
	w.Resize(fyne.NewSize(400, 400))
	// title of our form
	title := canvas.NewText("Contact Us", color.Black)
	title.TextSize = 20 // text font size is 20
	title.TextStyle = fyne.TextStyle{Bold: true}
	// I need it bold.. you can make it italic
	title.Resize(fyne.NewSize(300, 35)) // 300 is width & 35 is height
	title.Move(fyne.NewPos(50, 10))
	//position my widget
	//50 px from left and 10 px from top
	// copy / the setting(resize/ postion)
	// Name field
	entry_name := widget.NewEntry()
	entry_name.
		entry_name.SetText("Hello world")
	entry_name.Resize(fyne.NewSize(300, 35))
	entry_name.Move(fyne.NewPos(50, 50))
	// copy / paste for next widget email :)
	entry_email := widget.NewEntry()
	entry_email.SetPlaceHolder("Enter your email..")
	entry_email.Resize(fyne.NewSize(300, 35))
	entry_email.Move(fyne.NewPos(50, 100))
	// copy / paste for address
	entry_addr := widget.NewEntry()
	entry_addr.SetPlaceHolder("Enter your addr..")
	entry_addr.Resize(fyne.NewSize(300, 35))
	entry_addr.Move(fyne.NewPos(50, 150))
	// copy/ paste for last widget
	entry_msg := widget.NewMultiLineEntry()
	entry_msg.SetPlaceHolder("Enter your message..")
	// most important
	entry_msg.MultiLine = true // multiline input / text area for message
	entry_msg.Resize(fyne.NewSize(300, 140))
	entry_msg.Move(fyne.NewPos(50, 200))
	// Submit button
	submit_btn := widget.NewButton("SUBMIT", nil)
	submit_btn.Resize(fyne.NewSize(80, 30))
	// button need to be small as compared to textarea
	submit_btn.Move(fyne.NewPos(50, 355))
	// setup content
	// we are going to use container without layout// my favorite
	w.SetContent(
		container.NewWithoutLayout(
			title,
			entry_name,
			entry_email,
			entry_addr,
			entry_msg,
			submit_btn,
		),
	)
	// Show and run
	w.ShowAndRun()
}
