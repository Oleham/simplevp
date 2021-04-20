// settings.go contain functions for updating and retrieving user settings.
package db

type Setting struct {
	ID                                             uint
	URL, Email, Password, DownloadPath, UploadPath string
}

func Settings() *Setting {

	var current Setting

	sVPDB.FirstOrCreate(&current, 1)

	return &current
}

func UpdateSettings(entry *Setting) {
	// Updates the settings
	sVPDB.Save(&entry)

}
