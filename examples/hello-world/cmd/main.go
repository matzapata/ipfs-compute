package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type ComputeResponse struct {
	Data    string            `json:"data"`
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
}

func main() {
	response := ComputeResponse{
		Data:   fmt.Sprintf("Hello, %v!", os.Getenv("NAME")),
		Status: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	// give response to the client
	out, err := json.Marshal(response)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to marshal response: %v", err)
		return
	}
	fmt.Print(string(out))
}
