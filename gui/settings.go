// settings.go creates the settings interface
package gui

import (
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/Oleham/simplevp/db"
)

func newEntryWText(f func() *widget.Entry, text string) *widget.Entry {
	// Shortcut so I can set text and create Entry in one go
	e := f()
	e.SetText(text)
	return e
}

func settingsForm() *widget.Form {
	// creates the form for entering settings.

	settings := db.ShowSettings()

	// Init the different fields (keeping the references)
	fields := []struct {
		Field string
		Entry *widget.Entry
	}{
		{"URL", newEntryWText(widget.NewEntry, settings.URL)},
		{"Email", newEntryWText(widget.NewEntry, settings.Email)},
		{"Passord", newEntryWText(widget.NewPasswordEntry, settings.Password)},
		{"Download to", newEntryWText(widget.NewEntry, settings.DownloadPath)},
		{"Upload from", newEntryWText(widget.NewEntry, settings.UploadPath)},
	}

	form := new(widget.Form)

	form.SubmitText = "Save"

	form.OnSubmit = func() {
		newSettings := db.Setting{
			ID:           1,
			URL:          fields[0].Entry.Text,
			Email:        fields[1].Entry.Text,
			Password:     fields[2].Entry.Text,
			DownloadPath: fields[3].Entry.Text,
			UploadPath:   fields[4].Entry.Text,
		}

		confirmation := dialog.NewInformation("Settings saved", "The settings have been saved", mainWindow)
		confirmation.Show()

		db.UpdateSettings(&newSettings)
	}

	// Append the field names pluss entries.
	for _, f := range fields {
		form.Append(f.Field, f.Entry)
	}

	return form
}
