package db

import (
	"log"
	"time"

	"github.com/Oleham/simplevp/xtrf"
)

type Job struct {
	ID             string `gorm:"primaryKey"`
	Client         string // <-- Name of client, i.e. XTRF owner, to allow multiple different jobs in the DB. Maybe make this a relationship?
	IdNumber       string
	Smart          bool
	Status         string
	Name           string
	Type           string
	Quantity       float64
	Unit           string
	Value          float64
	Currency       string
	Deadline       int64
	DelieveryDate  int64
	ProjectManager string
	SourceLang     string
	TargetLang     string
	SourceFiles    []File
	Communication  string
	Invoice        string // <-- This will be the foreign key for invoice relation
}

func (j Job) DeadlineString() string {
	return time.Unix(j.Deadline, 0).Format("02/01/2006 15.04.05")
}

type File struct {
	ID                 string `gorm:"primaryKey"`
	Name, MetaCategory string
	JobID              string
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

	current := *Settings()
	newJobs := xtrf.JobsInProgress(current[0].URL, current[0].Email, current[0].Password)

	for _, item := range *newJobs {

		var entry Job

		entry.ID = string(item.Id)
		entry.Name = item.Main.ProjectName
		entry.Type = item.Main.Typus

		if len(item.Main.Name) > 18 {
			entry.Smart = true
		}

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

		sVPDB.Create(&entry)
	}
}
