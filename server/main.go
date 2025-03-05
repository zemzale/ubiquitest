package main

import (
	"fmt"
	"os"

	"github.com/zemzale/ubiquitest/router"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("Failed to run the server: %s\n", err.Error())
		os.Exit(1)
	}
}

func run() error {
	if err := router.Run(); err != nil {
		return err
	}

	return nil
}
