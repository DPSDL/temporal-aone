package main

import (
	"context"
	"encoding/json"
	"github.com/google/gops/agent"
	"go.temporal.io/sdk/client"
	"log"
	"net/http"
	"temporal-aone/backend/pkg"
	"time"
)

// 定义请求结构体
type WorkflowRequest struct {
	RepoURL        string `json:"repo_url"`
	Token          string `json:"token"`
	BinaryPath     string `json:"binary_path"`
	ConfigFilePath string `json:"config_file_path"`
	Version        string `json:"version"`
	ECSUploadPath  string `json:"ecs_upload_path"`
	ECSServer      string `json:"ecs_server"`
	ECSUser        string `json:"ecs_user"`
	ECSPassword    string `json:"ecs_password"`
	HealthCheckURL string `json:"health_check_url"`
}

// Handler for starting config workflow
func startConfigWorkflow(w http.ResponseWriter, r *http.Request) {
	var req WorkflowRequest

	// Parse the request body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create Temporal client
	c, err := client.NewClient(client.Options{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer c.Close()

	options := client.StartWorkflowOptions{
		ID:        "config-workflow",
		TaskQueue: "release-task-queue",
	}

	config := pkg.Config{
		RepoURL:        req.RepoURL,
		Token:          req.Token,
		BinaryPath:     req.BinaryPath,
		ConfigFilePath: req.ConfigFilePath,
		Version:        req.Version,
		ECSUploadPath:  req.ECSUploadPath,
		ECSServer:      req.ECSServer,
		ECSUser:        req.ECSUser,
		ECSPassword:    req.ECSPassword,
		HealthCheckURL: req.HealthCheckURL,
	}

	we, err := c.ExecuteWorkflow(context.Background(), options, pkg.ConfigWorkflow, config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"workflow_id": we.GetID(),
		"run_id":      we.GetRunID(),
	})
}

// Handler for starting build and upload workflow
func startBuildUploadWorkflow(w http.ResponseWriter, r *http.Request) {
	var req WorkflowRequest

	// Parse the request body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create Temporal client
	c, err := client.NewClient(client.Options{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer c.Close()

	options := client.StartWorkflowOptions{
		ID:        "build-upload-workflow",
		TaskQueue: "release-task-queue",
	}

	config := pkg.Config{
		RepoURL:        req.RepoURL,
		Token:          req.Token,
		BinaryPath:     req.BinaryPath,
		ConfigFilePath: req.ConfigFilePath,
		Version:        req.Version,
		ECSUploadPath:  req.ECSUploadPath,
		ECSServer:      req.ECSServer,
		ECSUser:        req.ECSUser,
		ECSPassword:    req.ECSPassword,
		HealthCheckURL: req.HealthCheckURL,
	}

	we, err := c.ExecuteWorkflow(context.Background(), options, pkg.BuildUploadWorkflow, config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"workflow_id": we.GetID(),
		"run_id":      we.GetRunID(),
	})
}

// Handler for starting release workflow
func startReleaseWorkflow(w http.ResponseWriter, r *http.Request) {
	var req WorkflowRequest

	// Parse the request body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create Temporal client
	c, err := client.NewClient(client.Options{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer c.Close()

	options := client.StartWorkflowOptions{
		ID:        "release-workflow",
		TaskQueue: "release-task-queue",
	}

	config := pkg.Config{
		RepoURL:        req.RepoURL,
		Token:          req.Token,
		BinaryPath:     req.BinaryPath,
		ConfigFilePath: req.ConfigFilePath,
		Version:        req.Version,
		ECSUploadPath:  req.ECSUploadPath,
		ECSServer:      req.ECSServer,
		ECSUser:        req.ECSUser,
		ECSPassword:    req.ECSPassword,
		HealthCheckURL: req.HealthCheckURL,
	}

	we, err := c.ExecuteWorkflow(context.Background(), options, pkg.ReleaseWorkflow, config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"workflow_id": we.GetID(),
		"run_id":      we.GetRunID(),
	})
}

func main() {
	http.HandleFunc("/api/start-config", startConfigWorkflow)
	http.HandleFunc("/api/start-build-upload", startBuildUploadWorkflow)
	http.HandleFunc("/api/start-release", startReleaseWorkflow)

	if err := agent.Listen(agent.Options{}); err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Hour)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
