package db

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Oleham/simplevp/xtrf"
)

type Job struct {
	ID             string `gorm:"primaryKey"`
	SettingID      uint
	Setting        Setting
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

// Type Jobslice implements sort
type JobSlice []Job

func (js JobSlice) Len() int           { return len(js) }
func (js JobSlice) Less(i, j int) bool { return js[i].Deadline < js[j].Deadline }
func (js JobSlice) Swap(i, j int)      { js[i], js[j] = js[j], js[i] }

type File struct {
	ID                 string `gorm:"primaryKey"`
	Name, MetaCategory string
	JobID              string
}

func deadlineToday(jobs *[]Job) []Job {

	var newList []Job

	for _, job := range *jobs {

		if job.Deadline > time.Now().Unix() {
			newList = append(newList, job)
		}
	}
	return newList
}

func Jobs() *[]Job {

	var jobs []Job

	result := sVPDB.Find(&jobs)
	if result.Error != nil {
		log.Fatalf(result.Error.Error())
	}

	currentJobs := deadlineToday(&jobs)

	return &currentJobs
}

func UpdateJobs() {

	currentSet := *Settings()

	for _, set := range currentSet {

		credentials := xtrf.Login(set.URL, set.Email, set.Password)

		newJobs, err := xtrf.Jobs(set.URL, credentials)
		if err != nil {
			log.Fatal(err)
		}

		for _, item := range *newJobs {

			var entry Job

			entry.ID = item.Id.String
			entry.Smart = item.Id.Smart
			entry.Name = item.Main.ProjectName
			entry.Type = item.Main.Typus
			entry.Status = item.Main.Status
			entry.SettingID = set.ID

			if weight := item.Main.JobQuantities.Weighted; len(weight) > 0 {
				entry.Quantity = weight[0].Value
				entry.Unit = weight[0].Unit
			}

			entry.Deadline = item.Main.Deadline.Integer
			entry.DelieveryDate = item.Main.DeliveryDate.Integer
			entry.ProjectManager = item.Main.ProjectManager.FirstName + " " +
				item.Main.ProjectManager.LastName
			entry.SourceLang = item.Main.SourceLanguage.Name

			if target := item.Main.Targets; len(target) > 0 {
				entry.TargetLang = target[0].Name
			}

			sVPDB.Create(&entry)

			//After creating the entry, update the files table if job is "in progress"
			//(Reducing uncessary API requests.)
			if item.Main.Status == "IN_PROGRESS" {
				UpdateFiles(set.URL, item.Id.String, credentials)
			}
		}
	}
}

func UpdateFiles(url, id string, cookies []*http.Cookie) {

	job := JobById(id)

	fileview, err := xtrf.File(url, job.ID, job.Smart, cookies)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range fileview.SourceFiles {
		var entry File
		entry.ID = file.ID.String
		entry.Name = file.Name
		if job.Smart {
			entry.MetaCategory = file.SmartCategory
		} else {
			entry.MetaCategory = file.Category
		}
		entry.JobID = job.ID

		sVPDB.Create(&entry)
	}
	if job.Smart {
		job.Communication = fmt.Sprintf("Instructions for all:\n%s\n\nInstructions for job:\n%s\n", fileview.InstructionsForAllJobs, fileview.InstructionsForJob)
	} else {
		job.Communication = fileview.Instructions
	}

	sVPDB.Save(job)

}

// jobById takes a job id and returns a job struct from the DB
func JobById(id string) *Job {
	var cur Job
	sVPDB.Where("id = ?", id).First(&cur)
	return &cur
}

// FilesOfJob returns all files associated with the job
func FilesByJob(id string) []File {

	var result []File

	sVPDB.
		Table("jobs").
		Select("files.id, files.name, files.meta_category", "files.job_id").
		Joins("join files on jobs.id = files.job_id").
		Where("jobs.id = ?", id).
		Scan(&result)

	return result
}
