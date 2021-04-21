package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/Oleham/simplevp/db"
)

func jobPage() *container.Scroll {
	// Displays the

	jobs := db.Jobs()

	// Refresh button
	refresh := widget.NewButton("Refresh", func() {
		db.UpdateJobs()
		mainWindow.SetContent(jobPage())
		confirmation := dialog.NewInformation("Updated", "Updated jobs", mainWindow)
		confirmation.Show()
	})

	title := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), widget.NewLabel("Viewing jobs"), layout.NewSpacer(), refresh)

	bilde := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), title)

	for _, job := range *jobs {

		bilde.Add(widget.NewLabel(fmt.Sprintf("%s -- %s -- %s", job.Name, job.DeadlineString(), job.ProjectManager)))
	}

	return container.NewScroll(bilde)
}
