package main

import (
	"log"

	executor_controller "github.com/matzapata/ipfs-compute/provider/internal/controllers/executor"
)

func main() {
	controller, err := executor_controller.NewApiHandler()
	if err != nil {
		log.Panic(err)
	}

	controller.Handle()
}
