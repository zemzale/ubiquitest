package main

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/samber/do"
	"github.com/zemzale/ubiquitest/container"
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
	container.Load()

	db, err := do.Invoke[*sqlx.DB](nil)
	if err != nil {
		return err
	}

	if err := storage.CreateDB(db); err != nil {
		return err
	}

	r, err := do.Invoke[*router.Router](nil)
	if err != nil {
		return err
	}

	if err := r.Run(); err != nil {
		return err
	}

	return nil
}
