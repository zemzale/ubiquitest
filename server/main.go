package main

import (
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/zemzale/ubiquitest/router"
	"github.com/zemzale/ubiquitest/storage"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("Failed to run the server: %s\n", err.Error())
		os.Exit(1)
	}
}

func run() error {
	db, err := storage.NewDB()
	if err != nil {
		return err
	}

	if err := storage.CreateDB(db); err != nil {
		return err
	}

	if err := router.Run(db); err != nil {
		return err
	}

	return nil
}
