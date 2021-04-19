package db

import "github.com/Oleham/simplevp/xtrf"

type Job struct {
	PK             uint   `gorm:"primaryKey"`
	VendorID       string `gorm:"unique"`
	Navn           string
	Quantity       float64
	Unit           string
	Deadline       int64
	ProjectManager string
	SourceLang     string
	TargetLang     string
	SourceFiles    bool
	StartTime      int64 `gorm:"autoCreateTime"`
	EndTime        int64
	Delievered     bool
}

func UpdateJobs() {

	current := ShowSettings()
	newJobs := xtrf.JobsInProgress(current.URL, current.Email, current.Password)

	for _, item := range *newJobs {

		entry := new(Job)

		entry.VendorID = item.Id
		entry.Navn = item.Overview.ProjectName
		entry.Quantity = item.Overview.JobQuantities.WeightedQuantities[0].Value
		entry.Unit = item.Overview.JobQuantities.WeightedQuantities[0].Unit
		entry.Deadline = item.Overview.Deadline.Unix()
		entry.ProjectManager = item.Overview.ProjectManager.FirstName +
			item.Overview.ProjectManager.LastName
		entry.SourceLang = item.Overview.SourceLanguage.Name
		entry.TargetLang = item.Overview.TargetLanguages[0].Name
		entry.SourceFiles = true
		entry.Delievered = false

		sVPDB.Create(entry)
	}
}
