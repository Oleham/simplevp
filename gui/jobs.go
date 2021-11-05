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
	"github.com/Oleham/simplevp/xtrf"
)

var checkedItems = make(map[string]binding.Bool)

// Function for download button
func downloadFiles() {

	var count int

	for k, v := range checkedItems {
		value, err := v.Get()
		if err != nil {
			fmt.Println("ops")
		}
		// if checked, list all work files. (to be downloaded)
		if value {
			files := db.FilesByJob(k)

			for _, f := range files {

				if f.MetaCategory == "WORKFILE" || f.MetaCategory == "WORK_FILE" {
					job := db.JobById(f.JobID)
					set := db.SettingById(job.SettingID)
					xtrf.Download(set.URL, set.Email, set.Password, set.DownloadPath, f.Name, f.JobID, f.ID, job.Smart)
					count += 1
				}
			}
		}
	}
	fmt.Printf("\nFinished downloading %d files.\n", count)
}

// Show description of job
func showDesc() {

	var confirmation dialog.Dialog

	for k, v := range checkedItems {
		value, err := v.Get()
		if err != nil {
			fmt.Println("ops")
		}

		if value {
			job := db.JobById(k)
			confirmation = dialog.NewInformation("Description", job.Communication, mainWindow)
			break
		} else {
			confirmation = dialog.NewInformation("Attention", "Select 1 single job to see description", mainWindow)
		}
	}
	confirmation.Show()
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

	checklist := fyne.NewContainerWithLayout(layout.NewVBoxLayout())

	for _, job := range *jobs {

		if job.Status == "IN_PROGRESS" {
			checkedItems[job.ID] = binding.NewBool()

			setting := db.SettingById(job.SettingID)

			var printName string

			if job.Name == "" {
				printName = job.IdNumber
			} else {
				printName = job.Name
			}

			checklist.Add(widget.NewCheckWithData(fmt.Sprintf("%s | %s -- %s (%s)", job.DeadlineString(), printName, job.Type, setting.Name), checkedItems[job.ID]))
		}
	}

	bilde := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), title, checklist, widget.NewButton("Description", showDesc), widget.NewButton("Download", downloadFiles))

	return container.NewScroll(bilde)
}
