package main

import (
	"log"

	gateway_controller "github.com/matzapata/ipfs-compute/provider/internal/controllers/gateway"
)

func main() {
	controller, err := gateway_controller.NewApiHandler()
	if err != nil {
		log.Fatal(err)
	}

	controller.Handle(":4000")
}
