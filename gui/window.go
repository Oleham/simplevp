package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var SIMPLEVP fyne.App = app.New()
var MAINWINDOW fyne.Window = SIMPLEVP.NewWindow("simpleVP 3.0")

func StartGUI() {
	mainMenu := menu(MAINWINDOW)
	MAINWINDOW.SetMainMenu(mainMenu)
	MAINWINDOW.SetContent(frontPage())
	MAINWINDOW.Resize(fyne.NewSize(300, 300))
	MAINWINDOW.ShowAndRun()
}

func frontPage() *fyne.Container {
	// Dette skal være forsiden i appen min.

	text := `Velkommen til den fine appen min.
    
Her vil du ha muligheten til å laste ned masse spennende oppdrag fra ulike XTRF-installasjoner.
    
Prøv det nå!`

	textObj := widget.NewLabel(text)

	button := widget.NewButton("Trykk meg!", func() { SIMPLEVP.SendNotification(&fyne.Notification{Title: "Varsel", Content: "Dette er en test"}) })

	page := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), textObj, button)

	return page
}

func jobber() *container.Scroll {
	labels := make([]fyne.CanvasObject, 20, 30)

	for i := 0; i < 20; i++ {
		labels[i] = widget.NewLabel(fmt.Sprintf("Jobb nr. %d", i))
	}

	bilde := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), labels...)

	return container.NewScroll(bilde)
}

func menu(window fyne.Window) *fyne.MainMenu {

	settings := fyne.NewMenuItem("Settings", func() { window.SetContent(settingsForm()) })

	jobber := fyne.NewMenuItem("Jobber", func() { window.SetContent(jobber()) })

	start := fyne.NewMenuItem("Home", func() { window.SetContent(frontPage()) })

	home := fyne.NewMenu("Home", start, jobber, settings)

	return fyne.NewMainMenu(home)
}
