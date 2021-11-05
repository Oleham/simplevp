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
	mainWindow.Resize(fyne.NewSize(1000, 750))
	mainWindow.ShowAndRun()
}

func frontPage() *fyne.Container {
	// Dette skal være forsiden i appen min.

	text := `Velkommen til den fine appen min.
    
Her vil du ha muligheten til å laste ned masse spennende oppdrag fra ulike XTRF-installasjoner.
    
Prøv det nå!`

	textObj := widget.NewLabel(text)

	page := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), textObj)

	return page
}

func menu(window fyne.Window) *fyne.MainMenu {

	settings := fyne.NewMenuItem("Settings", func() { window.SetContent(settingsPage()) })

	jobber := fyne.NewMenuItem("Jobs", func() { window.SetContent(jobPage()) })

	start := fyne.NewMenuItem("Home", func() { window.SetContent(frontPage()) })

	home := fyne.NewMenu("Home", start, jobber, settings)

	return fyne.NewMainMenu(home)
}
