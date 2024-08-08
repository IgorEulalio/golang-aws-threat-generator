package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Welcome to Golang event generator!")
}

func filesHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "Path query parameter is required", http.StatusBadRequest)
		return
	}

	dir, err := os.ReadDir(path)
	if err != nil {
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

func eventsHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "Path query parameter is required", http.StatusBadRequest)
		return
	}

	// Create a file named with the current timestamp and a .json extension
	timestamp := time.Now().Format("20060102150405")
	fileName := fmt.Sprintf("%s.json", timestamp)
	filePath := filepath.Join(path, fileName)

	// Create the directory if it doesn't exist
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create directory: %v", err), http.StatusInternalServerError)
		return
	}

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create file: %v", err), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create random UUID: %v", err), http.StatusInternalServerError)
		return
	}

	uuidData := map[string]string{"uuid": randomUUID.String()}
	jsonData, err := json.Marshal(uuidData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal JSON: %v", err), http.StatusInternalServerError)
		return
	}

	_, err = file.Write(jsonData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to write to file: %v", err), http.StatusInternalServerError)
		return
	}

	// Return the path of the file
	fmt.Fprintf(w, "File created: %s", filePath)
}

func main() {
	app := http.HandlerFunc(hello)

	http.HandleFunc("/events", eventsHandler)
	http.HandleFunc("/files", filesHandler)
	http.Handle("/", app)
	fmt.Printf("Starting server at port 8080\n")
	http.ListenAndServe(":8080", nil)
}
