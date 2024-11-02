package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/IgorEulalio/golang-threat-generator/pkg/client"
	"github.com/IgorEulalio/golang-threat-generator/pkg/events"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Welcome to Golang event generator!")
}

func FilesHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "Path query parameter is required", http.StatusBadRequest)
		return
	}

	dir, err := os.ReadDir(path)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read directory: %v", err), http.StatusInternalServerError)
		return
	}

	var fileList []string
	for _, file := range dir {
		fileList = append(fileList, file.Name())
	}

	jsonData, err := json.Marshal(fileList)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal JSON: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func EventsHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "Path query parameter is required", http.StatusBadRequest)
		return
	}

	// Generate event
	filePath, err := events.GenerateEvent(path)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate event: %v", err), http.StatusInternalServerError)
		return
	}

	// Return the path of the file
	fmt.Fprintf(w, "File created: %s", filePath)
}

func IamEnumeratorHandler(w http.ResponseWriter, r *http.Request) {

	awsClient := client.GetAWSClient()
	iamEnumerator := events.IAMEnumerator{AWSClient: awsClient}
	err := iamEnumerator.EnumerateRolesThatCanBeAssumed()
	if err != nil {
		return
	}

	// Enumerate roles that can be assumed
	//if err != nil {
	//	http.Error(w, fmt.Sprintf("Failed to enumerate roles: %v", err), http.StatusInternalServerError)
	//	return
	//}
	//
	//fmt.Fprintf(w, "Roles enumerated successfully")
}
