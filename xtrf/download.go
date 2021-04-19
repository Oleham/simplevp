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
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/vendors/jobs?statuses=IN_PROGRESS", baseURL), nil)
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

	return unpack(jn)

}
