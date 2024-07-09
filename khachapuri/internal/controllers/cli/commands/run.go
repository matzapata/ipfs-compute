package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/matzapata/ipfs-compute/provider/internal/config"
	"github.com/matzapata/ipfs-compute/provider/internal/services"
	"github.com/matzapata/ipfs-compute/provider/pkg/archive"
)

func RunCommand(config *config.Config) error {
	unzipper := archive.NewUnzipper()
	computeExecutor := services.NewComputeExecutor(nil)

	// check there's a local build
	if _, err := os.Stat(".khachapuri"); os.IsNotExist(err) {
		return errors.New("no build found. Run 'khachapuri build' first")
	}
	// unzip the build
	unzippedArtifact, err := filepath.Abs("./tmp")
	if err != nil {
		return err
	}
	if err = unzipper.Unzip(".khachapuri/artifact.zip", unzippedArtifact); err != nil {
		return err
	}
	defer os.RemoveAll(unzippedArtifact)

	// load .env file
	envVars, err := godotenv.Read(".env.local")
	if err != nil {
		return errors.New("error loading .env.local file")
	}
	// map envVars to NAME=VALUE format
	var execEnv []string
	for k, v := range envVars {
		execEnv = append(execEnv, k+"="+v)
	}

	// run a remote build
	response, err := computeExecutor.Execute(unzippedArtifact, execEnv, "execArgs")
	if err != nil {
		return err
	}
	fmt.Println("\nStatus:", response.Status)
	fmt.Println("Data:", response.Data)
	fmt.Println("Headers:")
	for k, v := range response.Headers {
		fmt.Println(" - ", k, ":", v)
	}

	return nil
}
