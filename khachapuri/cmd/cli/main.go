package main

import (
	cli_controller "github.com/matzapata/ipfs-compute/provider/internal/controllers/cli"
)

func main() {
	controller := cli_controller.NewCliHandler()
	controller.Handle()
}
