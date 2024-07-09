package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/matzapata/ipfs-compute/provider/internal/config"
	"github.com/matzapata/ipfs-compute/provider/internal/domain"
)

// ComputeExecutor handles executing computations inside Docker containers.
type ComputeExecutor struct {
	DockerClient *client.Client
}

// NewComputeExecutor creates a new ComputeExecutor with the given Docker client.
func NewComputeExecutor(
	cfg *config.Config,
) *ComputeExecutor {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	return &ComputeExecutor{
		DockerClient: cli,
	}
}

// Execute executes a computation inside a Docker container with a timeout of 15 seconds.
func (e *ComputeExecutor) Execute(deploymentPath string, execEnv []string, execArgs string) (*domain.ComputeResponse, error) {
	// Prepare container options
	ctx := context.Background()

	// Define volume mount
	volumeMount := fmt.Sprintf("%s:/app", deploymentPath)

	// Create container configuration
	resp, err := e.DockerClient.ContainerCreate(
		ctx,
		&container.Config{
			Image: "alpine", // Name of the Docker image
			Cmd:   []string{"./app/main", execArgs},
			Env:   execEnv,
		},
		&container.HostConfig{
			Binds: []string{volumeMount},
		},
		nil,
		nil,
		"",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create container: %v", err)
	}

	// Start the container
	if err := e.DockerClient.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return nil, fmt.Errorf("failed to start container: %v", err)
	}

	// Create a context with a timeout of 15 seconds
	execCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// Wait for the container to finish within the timeout
	statusCh, errCh := e.DockerClient.ContainerWait(execCtx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return nil, fmt.Errorf("container wait error: %v", err)
		}
	case <-statusCh:
	}

	// Retrieve container logs
	reader, err := e.DockerClient.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve container logs: %v", err)
	}
	defer reader.Close()

	// Read the first 8 bytes to ignore the HEADER part from docker container logs
	p := make([]byte, 8)
	reader.Read(p)
	content, _ := io.ReadAll(reader)

	// Parse response from container output
	var response domain.ComputeResponse
	if err := json.Unmarshal(content, &response); err != nil {
		return nil, fmt.Errorf("failed to parse output: %v", err)
	}

	return &response, nil
}
