package db

import (
	"log"
	"time"

	"github.com/Oleham/simplevp/xtrf"
)

type Job struct {
	PK             uint   `gorm:"primaryKey"`
	VendorID       string `gorm:"unique"`
	Navn           string
	Quantity       float64
	Unit           string
	Deadline       int64 `gorm:"type:datetime"`
	ProjectManager string
	SourceLang     string
	TargetLang     string
	SourceFiles    bool
	StartTime      int64 `gorm:"autoCreateTime"`
	EndTime        int64
	Delievered     bool
}

func (j Job) DeadlineString() string {
	return time.Unix(j.Deadline, 0).Format("02/01/2006 15.04.05")
}

func Jobs() *[]Job {

	var jobs []Job

	result := sVPDB.Find(&jobs)
	if result.Error != nil {
		log.Fatalf(result.Error.Error())
	}
	return &jobs
}

func UpdateJobs() {

	current := Settings()
	newJobs := xtrf.JobsInProgress(current.URL, current.Email, current.Password)

	for _, item := range *newJobs {

		var entry Job

		entry.VendorID = string(item.Id)
		entry.Navn = item.Main.ProjectName

		if weight := item.Main.JobQuantities.Weighted; len(weight) > 0 {
			entry.Quantity = weight[0].Value
			entry.Unit = weight[0].Unit
		}

		entry.Deadline = item.Main.Deadline.Unix()
		entry.ProjectManager = item.Main.ProjectManager.FirstName + " " +
			item.Main.ProjectManager.LastName
		entry.SourceLang = item.Main.SourceLanguage.Name

		if target := item.Main.Targets; len(target) > 0 {
			entry.TargetLang = target[0].Name
		}

		entry.SourceFiles = true
		entry.Delievered = false

		sVPDB.Create(&entry)
	}
}
