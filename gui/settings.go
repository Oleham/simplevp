// settings.go creates the settings interface
package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/Oleham/simplevp/db"
)

func newEntryWText(f func() *widget.Entry, text string) *widget.Entry {
	// Shortcut so I can set text and create Entry in one go
	e := f()
	e.SetText(text)
	return e
}

func settingsPage() *fyne.Container {
	// Displays the settings

	// Radio buttons
	radio := widget.NewRadioGroup([]string{}, func(string) {})

	settings := db.Settings()

	// Create button
	create := widget.NewButton("Create", func() {
		mainWindow.SetContent(settingsForm(&db.Setting{}))
	})

	// Update button
	update := widget.NewButton("Update", func() {
		sel := radio.Selected
		mainWindow.SetContent(settingsForm(db.SettingByName(sel)))
	})

	// Delete button
	delete := widget.NewButton("Delete", func() {

		// Function called when confirmbox is clicked.
		delfunc := func(check bool) {
			if check {
				sel := radio.Selected
				db.DeleteSetting(db.SettingByName(sel))
				dialog.NewInformation("Client deleted", "The client was deleted", mainWindow).Show()
				mainWindow.SetContent(settingsPage())
			}
		}
		// Create and show confirm box
		dialog.NewConfirm("Delete Client", "Are you sure?", delfunc, mainWindow).Show()
	})

	// Title bar with buttons
	title := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), widget.NewLabel("All XTRF Connections"), create, update, delete)

	for _, set := range *settings {
		radio.Append(set.Name)
	}

	// Container
	return fyne.NewContainerWithLayout(layout.NewVBoxLayout(), title, radio)
}

func settingsForm(settings *db.Setting) *widget.Form {
	// creates the form for entering settings.

	keys := []string{"URL", "Email", "Password", "DownloadPath", "UploadPath"}

	fields := map[string]*widget.Entry{
		keys[0]: newEntryWText(widget.NewEntry, settings.URL),
		keys[1]: newEntryWText(widget.NewEntry, settings.Email),
		keys[2]: newEntryWText(widget.NewPasswordEntry, settings.Password),
		keys[3]: newEntryWText(widget.NewEntry, settings.DownloadPath),
		keys[4]: newEntryWText(widget.NewEntry, settings.UploadPath),
	}

	// Check if we are updating or creating
	var creating bool

	if settings.Name == "" {
		creating = true
		fields["Name"] = newEntryWText(widget.NewEntry, "New Client")
	}

	form := new(widget.Form)

	// Save button
	form.SubmitText = "Save"
	form.OnSubmit = func() {

		if creating {
			settings.Name = fields["Name"].Text
		}
		settings.URL = fields["URL"].Text
		settings.Email = fields["Email"].Text
		settings.Password = fields["Password"].Text
		settings.DownloadPath = fields["DownloadPath"].Text
		settings.UploadPath = fields["UploadPath"].Text

		db.UpdateSetting(settings)

		confirmation := dialog.NewInformation("Settings saved", "The settings have been saved", mainWindow)
		confirmation.Show()

		mainWindow.SetContent(settingsPage())
	}

	form.OnCancel = func() {
		mainWindow.SetContent(settingsPage())
	}

	// Append the field names pluss entries.
	if creating {
		form.Append("Name", fields["Name"])
	}
	for _, k := range keys {
		form.Append(k, fields[k])
	}

	return form
}
