package gui

import (
	"fmt"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/Oleham/simplevp/db"
)

var checkedItems = make(map[string]binding.Bool)

// Function for download button
func downloadFunc() {
	for k, v := range checkedItems {
		value, err := v.Get()
		if err != nil {
			fmt.Println("ops")
		}
		if value {
			fmt.Println(k)
		}

	}
}

// jobPage returns the scrolling container with all the jobs
// the function creates buttons, accordions and check items (and their bindings)
func jobPage() *container.Scroll {

	jobs := db.Jobs()

	sort.Sort(db.JobSlice(*jobs))

	// Refresh button
	refresh := widget.NewButton("Refresh", func() {
		db.UpdateJobs()
		mainWindow.SetContent(jobPage())
		confirmation := dialog.NewInformation("Updated", "Updated jobs", mainWindow)
		confirmation.Show()
	})

	title := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), widget.NewLabel("Viewing jobs"), layout.NewSpacer(), refresh)

	accordion := widget.NewAccordion()

	for _, job := range *jobs {

		if job.Status == "IN_PROGRESS" {
			//bilde.Add(widget.New(fmt.Sprintf("%s -- %s -- %s", job.Name, job.DeadlineString(), job.SourceFiles), func(bool) {}))
			job.SourceFiles = db.FilesByJob(job.ID)

			checklist := fyne.NewContainerWithLayout(layout.NewVBoxLayout())

			for _, f := range job.SourceFiles {

				checkedItems[f.Name] = binding.NewBool()

				checklist.Add(widget.NewCheckWithData(fmt.Sprintf("%s (%s)", f.Name, f.MetaCategory), checkedItems[f.Name]))
			}
			accordion.Append(widget.NewAccordionItem(fmt.Sprintf("%s | %s", job.DeadlineString(), job.Name), checklist))
		}
	}

	bilde := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), title, accordion, widget.NewButton("Download", downloadFunc))

	return container.NewScroll(bilde)
}
