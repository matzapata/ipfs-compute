package services

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func (cs *ComputeService) Compute(deploymentPath string, execEnv []string) (*ComputeResponse, error) {
	// Prepare the docker run command
	args := []string{"run", "--rm", "-v", fmt.Sprintf("%s:/app", deploymentPath)}
	for _, env := range execEnv {
		args = append(args, "-e", env)
	}
	args = append(args, "binary_runner")

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
