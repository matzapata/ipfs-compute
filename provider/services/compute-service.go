package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
)

type ComputeService struct {
}

type ComputeResponse struct {
	Data    string            `json:"data"`
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
}

func NewComputeService() *ComputeService {
	return &ComputeService{}
}

func (cs *ComputeService) Compute(deploymentPath string, execEnv []string, execArgs string) (*ComputeResponse, error) {
	// Prepare the docker run command
	args := []string{"run", "--rm", "-v", fmt.Sprintf("%s:/app", deploymentPath)}
	for _, env := range execEnv {
		args = append(args, "-e", env)
	}
	args = append(args, "binary_runner", "main") // binary_runner is the name of the docker image
	args = append(args, execArgs)

	// run the binary inside the docker container
	cmd := exec.Command("docker", args...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("execution error: %v, stderr: %v", err, stderr.String())
	}

	var response ComputeResponse
	err = json.Unmarshal(out.Bytes(), &response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse output: %v", err)
	}

	return &response, nil
}

// Creates curl like command to be executed in the gateway
func (cs *ComputeService) ParseRequest(r *http.Request) (string, error) {
	// Prepare the curl command
	args := []string{"-X", r.Method}

	// add headers
	for key, value := range r.Header {
		args = append(args, "-H", fmt.Sprintf("%s: %s", key, value[0]))
	}
	args = append(args, r.URL.String())

	// add data
	if r.Method == "POST" || r.Method == "PUT" {
		body, err := r.GetBody()
		if err != nil {
			return "", fmt.Errorf("failed to get body: %v", err)
		}
		data, err := io.ReadAll(body)
		if err != nil {
			return "", fmt.Errorf("failed to read body: %v", err)
		}
		args = append(args, "-d", string(data))
	}

	// TODO: add more methods

	return fmt.Sprintf("%s", args), nil
}
