// models.go will contain the models for the database.
package db

type Setting struct {
	ID                                             uint
	URL, Email, Password, DownloadPath, UploadPath string
}
