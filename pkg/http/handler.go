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

func IamRoleEnumeratorHandler(w http.ResponseWriter, r *http.Request) {

	awsClient := client.GetAWSClient()
	iamEnumerator := events.IAMRoleEnumerator{AWSClient: awsClient}

	roles, err := iamEnumerator.EnumerateRolesThatCanBeAssumed()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to enumerate roles that can be assumed: %v", err), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(roles)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal roles to JSON: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func AssumeRoleHandler(w http.ResponseWriter, r *http.Request) {

	awsClient := client.GetAWSClient()
	roleAssumer := events.AssumeRole{AWSClient: awsClient}

	// get from request
	roleArn, err := events.DecodeBodyIntoAssumeRole(r.Body)
	if err != nil {
		return
	}

	output, err := roleAssumer.AssumeByArn(*roleArn)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to assume role: %v", err), http.StatusInternalServerError)
		return
	}

	payload := events.ResponsePayload{
		RoleArn:    *roleArn,
		Expiration: output.Credentials.Expiration.String(),
	}

	jsonResponse, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal assumeRoleOutput to JSON: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func IamGroupsAndUserEnumeratorHandler(w http.ResponseWriter, r *http.Request) {
	enumerator := events.IAMUserEnumerator{}
	err := enumerator.EnumerateUserAndPolicy()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to enumerate IAM users and policies: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("IAM User and Policy enumerated"))
}
