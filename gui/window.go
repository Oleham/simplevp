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

	text := `Welcome to simpleVP.
    
Use the dropdown in the top left corner to navigate.

Create a client in the "Settings"section. 
Add the base URL of the XTRF web site, your email and password. 
Also specify a download folder.
URL and download folder should be WITHOUT trailing slash, i.e. "C:\Users\user\myfolder" and "https://xtrf.myclient.com". 
    
After adding a client, navigate to the "Jobs" section and click "Update" to see all current jobs (deadline not yet reached).
Download all work files or read the description of the task.`

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
