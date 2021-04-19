// json.go contains structs for XTRF json
// and functions for parsing
package xtrf

import (
	"encoding/json"
	"log"
	"time"
)

type XTRFJob struct {
	Id       string `json:"id"`
	Overview `json:"overview"`
}

type Overview struct {
	IdNumber        string    `json:"idNumber"`
	ProjectName     string    `json:"projectName"`
	Type            string    `json:"type"`
	Deadline        time.Time `json:"deadline"`
	ProjectManager  `json:"projectManager"`
	JobQuantities   `json:"jobQuantities"`
	SourceLanguage  `json:"sourceLanguage"`
	TargetLanguages []TargetLanguages `json:"targetLanguages"`
}

type ProjectManager struct {
	FirstName, LastName string
}

type SourceLanguage struct {
	Name string
}

type TargetLanguages struct {
	Name string
}

type JobQuantities struct {
	WeightedQuantities []weightedQuantities
}

type weightedQuantities struct {
	Value float64
	Unit  string
}

func unpack(js *[]byte) *[]XTRFJob {
	//Unpacks the JSON

	var jobber *[]XTRFJob

	err := json.Unmarshal(*js, jobber)
	if err != nil {
		log.Fatal(err)
	}

	return jobber

}
