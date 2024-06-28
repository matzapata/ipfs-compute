package main

import executor_controller "github.com/matzapata/ipfs-compute/provider/internal/controllers/executor"

func main() {
	controller := executor_controller.NewApiHandler()
	controller.Handle()
}
