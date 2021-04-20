package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/Oleham/simplevp/db"
)

func jobber() *container.Scroll {

	jobs := db.Jobs()

	refresh := widget.NewButton("Refresh", func() { db.UpdateJobs(); mainWindow.SetContent(jobber()) })

	title := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), widget.NewLabel("Jobber"), layout.NewSpacer(), refresh)

	bilde := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), title)

	for _, job := range *jobs {

		bilde.Add(widget.NewLabel(fmt.Sprintf("%s -- %s -- %s", job.Navn, job.DeadlineString(), job.ProjectManager)))
	}

	return container.NewScroll(bilde)
}
