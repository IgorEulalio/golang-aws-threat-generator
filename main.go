package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/IgorEulalio/golang-threat-generator/internal/config"
	"github.com/IgorEulalio/golang-threat-generator/pkg/client"
	httphandlers "github.com/IgorEulalio/golang-threat-generator/pkg/http"
)

func main() {
	// Load configuration (if any)
	cfg := config.LoadConfig()

	err := client.Init()
	if err != nil {
		log.Printf("Failed to initialize AWS client, AWS integration will not be available: %v", err)
	}
	log.Printf("Starting app, listening on port: %v", cfg.Port)

	// Initialize HTTP handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/", httphandlers.HelloHandler)
	mux.HandleFunc("/events", httphandlers.EventsHandler)
	mux.HandleFunc("/files", httphandlers.FilesHandler)
	mux.HandleFunc("/iamRolesEnumeration", httphandlers.IamRoleEnumeratorHandler)
	mux.HandleFunc("/iamUserEnumeration", httphandlers.IamGroupsAndUserEnumeratorHandler)
	mux.HandleFunc("/assumeRole", httphandlers.AssumeRoleHandler)
	addr := fmt.Sprintf(":%d", cfg.Port)
	fmt.Printf("Starting server at port %d\n", cfg.Port)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
