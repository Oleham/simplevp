package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var simpleVP fyne.App = app.New()
var mainWindow fyne.Window = simpleVP.NewWindow("simpleVP 3.0")

func StartGUI() {
	mainMenu := menu(mainWindow)
	mainWindow.SetMainMenu(mainMenu)
	mainWindow.SetContent(frontPage())
	mainWindow.Resize(fyne.NewSize(300, 300))
	mainWindow.ShowAndRun()
}

func frontPage() *fyne.Container {
	// Dette skal være forsiden i appen min.

	text := `Velkommen til den fine appen min.
    
Her vil du ha muligheten til å laste ned masse spennende oppdrag fra ulike XTRF-installasjoner.
    
Prøv det nå!`

	textObj := widget.NewLabel(text)

	button := widget.NewButton("Trykk meg!", func() { simpleVP.SendNotification(&fyne.Notification{Title: "Varsel", Content: "Dette er en test"}) })

	page := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), textObj, button)

	return page
}

func menu(window fyne.Window) *fyne.MainMenu {

	settings := fyne.NewMenuItem("Settings", func() { window.SetContent(settingsForm()) })

	jobber := fyne.NewMenuItem("Jobber", func() { window.SetContent(jobber()) })

	start := fyne.NewMenuItem("Home", func() { window.SetContent(frontPage()) })

	home := fyne.NewMenu("Home", start, jobber, settings)

	return fyne.NewMainMenu(home)
}
