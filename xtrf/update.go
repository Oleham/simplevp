// Download all current jobs from XTRF
package xtrf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type XTRFJob struct {
	Id   VendorID `json:"id"`
	Main Overview `json:"overview"`
}

type Overview struct {
	IdNumber       string     `json:"idNumber"`
	ProjectName    string     `json:"projectName"`
	Typus          string     `json:"type"`
	Deadline       VendorTime `json:"deadline"`
	ProjectManager `json:"projectManager"`
	JobQuantities  `json:"jobQuantities"`
	SourceLanguage `json:"sourceLanguage"`
	Targets        []TargetLanguages `json:"targetLanguages"`
}

type ProjectManager struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type SourceLanguage struct {
	Name string `json:"name"`
}

type TargetLanguages struct {
	Name string `json:"name"`
}

type JobQuantities struct {
	Weighted []WeightedQuantities `json:"weightedQuantities"`
}

type WeightedQuantities struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

type VendorID string

func (v *VendorID) UnmarshalJSON(s []byte) (err error) {
	// I need to create a custom Unmarshal method to deal with VendorID which might be string or int.
	*v = VendorID(string(s))
	return
}

type VendorTime time.Time

func (v *VendorTime) UnmarshalJSON(s []byte) (err error) {

	q, err := strconv.ParseInt(string(s), 10, 64)
	if err != nil {
		return err
	}
	*v = VendorTime(time.Unix(q/1000, 0))
	return
}

func (v VendorTime) String() string { return time.Time(v).Format("02.01.2006 15.04.05") }
func (v VendorTime) Unix() int64    { return time.Time(v).Unix() }

func unpack(js *[]byte) (*[]XTRFJob, error) {
	//Unpacks the JSON

	var jobber *[]XTRFJob

	err := json.Unmarshal(*js, &jobber)
	if err != nil {
		return nil, fmt.Errorf("Feil under json.Unmarshal: %v", err)
	}

	return jobber, nil
}

func getJobs(baseURL, email, pw string) *[]byte {
	// Function takes email, pw and url to login to XTRF and download not invoiced jobs.
	// Returns []byte containing the JSON string
	client := &http.Client{}

	// Create json body for login
	body, err := json.Marshal(map[string]string{
		"email":    email,
		"password": pw,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Post request to login
	resp, err := client.Post(fmt.Sprintf("%s/vendors/sign-in", baseURL), "application/json;charset=utf-8", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Get the cookies from login
	var cookies []*http.Cookie = resp.Cookies()

	// Read the reponse to close
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Build next request
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/vendors/jobs?statuses=IN_PROGRESS,PENDING,NOT_INVOICED", baseURL), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Add the cookies from the last request.
	// Seems the other headers aren't needed
	for i := 0; i < len(cookies); i++ {
		request.AddCookie(cookies[i])
	}

	resp, err = client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	text, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return &text
}

func JobsInProgress(url, email, pw string) *[]XTRFJob {
	// Download current jobs in progress from URL

	jn := getJobs(url, email, pw)

	jobber, err := unpack(jn)
	if err != nil {
		log.Fatal(err)
	}

	return jobber
}
