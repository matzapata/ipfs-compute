package main

import (
	"log"

	provider_controller "github.com/matzapata/ipfs-compute/provider/internal/controllers/provider"
)

func main() {
	controller, err := provider_controller.NewApiHandler()
	if err != nil {
		log.Panic(err)
	}

	controller.Handle(":3000")
}
