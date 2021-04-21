// settings.go contain functions for updating and retrieving user settings.
package db

import "log"

type Setting struct {
	ID                                             uint   `gorm:"primaryKey"`
	Name                                           string `gorm:"unique"`
	URL, Email, Password, DownloadPath, UploadPath string
}

func Settings() *[]Setting {

	var current []Setting

	result := sVPDB.Find(&current)
	if result.Error != nil {
		log.Fatalf(result.Error.Error())
	}

	return &current
}

func SettingByName(name string) *Setting {
	var cur Setting
	sVPDB.Where("name = ?", name).FirstOrCreate(&cur)
	return &cur
}

func UpdateSetting(entry *Setting) {
	// Updates the settings
	sVPDB.Save(entry)
}

func DeleteSetting(entry *Setting) {
	// Delete setting
	sVPDB.Delete(entry)
}
