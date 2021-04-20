package db

import (
	"fmt"

	"github.com/Oleham/simplevp/xtrf"
)

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
	fmt.Println(current)
	newJobs := xtrf.JobsInProgress(current.URL, current.Email, current.Password)

	for _, item := range *newJobs {

		var entry Job

		entry.VendorID = string(item.Id)
		entry.Navn = item.Main.ProjectName
		entry.Quantity = item.Main.JobQuantities.Weighted[0].Value
		entry.Unit = item.Main.JobQuantities.Weighted[0].Unit
		entry.Deadline = item.Main.Deadline.Unix()
		entry.ProjectManager = item.Main.ProjectManager.FirstName +
			item.Main.ProjectManager.LastName
		entry.SourceLang = item.Main.SourceLanguage.Name
		entry.TargetLang = item.Main.Targets[0].Name
		entry.SourceFiles = true
		entry.Delievered = false

		sVPDB.Create(&entry)
	}
}
