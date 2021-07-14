package xtrf

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func file(baseURL, jobID, fileID string, smart bool) *http.Request {

	var format string

	// Smart
	if smart {
		format = "%s/vendors/jobs/smart/%s/files/%s"
		// Classic
	} else {
		format = "%s/vendors/jobs/classic/%s/files/%s"
	}

	request, err := http.NewRequest("GET", fmt.Sprintf(format, baseURL, jobID, fileID), nil)
	if err != nil {
		log.Fatal(err)
	}

	return request

}

func Download(baseURL, email, pw, downloadPath, filename, jobID, fileID string, smart bool) {

	if filename == "" {
		filename = fileID + ".txt"
	}

	fullDownloadPath := downloadPath + filename

	response := requestJSON(file(baseURL, jobID, fileID, smart), Login(baseURL, email, pw))

	err := os.WriteFile(fullDownloadPath, *response, 0666)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s was downloaded!", fullDownloadPath)

}
