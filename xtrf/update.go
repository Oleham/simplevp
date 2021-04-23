// Download all current jobs from XTRF
package xtrf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type XTRFJob struct {
	Id   VendorID `json:"id"`
	Main Overview `json:"overview"`
}

type Overview struct {
	IdNumber       string     `json:"idNumber"`
	ProjectName    string     `json:"projectName"`
	Typus          string     `json:"type"`
	Status         string     `json:"status"`
	Deadline       VendorTime `json:"deadline"`
	DeliveryDate   VendorTime `json:"deliveryDate"`
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

// I need to create a custom Unmarshal method to deal with VendorID which might be string or int.
// This also determines whether the job is a smart job or classic job.
type VendorID struct {
	String string
	Smart  bool
}

func (v *VendorID) UnmarshalJSON(s []byte) (err error) {
	err = json.Unmarshal(s, &v.String)
	if err != nil {
		// If the unmarshal fails, try to treat it as an int.
		var i int
		err = json.Unmarshal(s, &i)
		if err != nil {
			return
		}
		v.String = fmt.Sprint(i)
		return
	}
	// If there was no error in the first place, we assume smart job.
	v.Smart = true
	return
}

// Custom VendorTime is for umarshalling a milisecond unix time into normal unix time.
type VendorTime struct {
	Integer int64
}

func (v *VendorTime) UnmarshalJSON(s []byte) (err error) {

	var i int64
	err = json.Unmarshal(s, &i)
	if err != nil {
		return
	}
	v.Integer = i / 1000
	return
}

type XTRFFile struct {
	SourceFiles        []SourceFiles `json:"sourceFiles"`
	SmartCommunication `json:"communication"`
	Instructions       string `json:"instructions"`
}

type SourceFiles struct {
	ID            VendorID `json:"id"`
	Name          string   `json:"name"`
	SmartCategory string   `json:"metaCategory"`
	Category      string   `json:"category"`
}

type SmartCommunication struct {
	InstructionsForJob     string `json:"instructionsForJob"`
	InstructionsForAllJobs string `json:"instructionsForAllJobs"`
}

func Login(baseURL, email, pw string) []*http.Cookie {
	// Function takes email, pw and url to login to XTRF
	// Returns cookies
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

	return cookies
}

func filesRequest(baseURL, jobID string, smart bool) *http.Request {
	// Create request for files (smart or classic)

	var format string
	if smart {
		format = "%s/vendors/jobs/smart/%s"
	} else {
		format = "%s/vendors/jobs/classic/%s"
	}

	request, err := http.NewRequest("GET", fmt.Sprintf(format, baseURL, jobID), nil)
	if err != nil {
		log.Fatal(err)
	}
	return request
}

func jobsRequest(baseURL string) *http.Request {
	// Create request for job list (IN_PROGRESS,PENDING,NOT_INVOICED)

	request, err := http.NewRequest("GET", fmt.Sprintf("%s/vendors/jobs?statuses=IN_PROGRESS,PENDING,NOT_INVOICED", baseURL), nil)
	if err != nil {
		log.Fatal(err)
	}

	return request
}

func requestJSON(request *http.Request, cookies []*http.Cookie) *[]byte {
	// Make a request to the XTRF API.
	// Takes request and session cookie as argument.

	client := &http.Client{}

	// Adding cookies
	for i := 0; i < len(cookies); i++ {
		request.AddCookie(cookies[i])
	}

	resp, err := client.Do(request)
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

func Jobs(url string, cookies []*http.Cookie) (*[]XTRFJob, error) {
	// Download current jobs in progress from URL

	jn := requestJSON(jobsRequest(url), cookies)

	var jobber []XTRFJob

	err := json.Unmarshal(*jn, &jobber)
	if err != nil {
		return &jobber, fmt.Errorf("Klarte ikke å parse jobber JSON fra %s: %v", url, err)
	}
	return &jobber, nil
}

func File(url, job string, smart bool, cookies []*http.Cookie) (*XTRFFile, error) {
	// Update the jobs with Files field and Description fields.

	jn := requestJSON(filesRequest(url, job, smart), cookies)

	var fil XTRFFile

	err := json.Unmarshal(*jn, &fil)
	if err != nil {
		return &fil, fmt.Errorf("Klarte ikke å parse fil JSON.\nJobb: %s\nURL: %s\nError: %v\n\nJSON: %s", job, url, err, jn)
	}
	return &fil, nil
}
