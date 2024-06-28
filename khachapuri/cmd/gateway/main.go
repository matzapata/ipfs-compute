package main

import gateway_controller "github.com/matzapata/ipfs-compute/provider/internal/controllers/gateway"

func main() {
	controller := gateway_controller.NewApiHandler()
	controller.Handle()
}
