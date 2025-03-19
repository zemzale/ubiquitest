package container

import (
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
	"github.com/zemzale/ubiquitest/config"
	"github.com/zemzale/ubiquitest/domain/tasks"
	"github.com/zemzale/ubiquitest/router"
	"github.com/zemzale/ubiquitest/storage"
)

func Load() {
	do.Provide(nil, func(i *do.Injector) (*config.Config, error) {
		return config.Load(), nil
	})

	do.Provide(nil, func(i *do.Injector) (*sqlx.DB, error) {
		cfg, err := do.Invoke[*config.Config](i)
		if err != nil {
			return nil, err
		}

		db, err := sqlx.Open(cfg.DB.Driver, cfg.DB.DSN)
		if err != nil {
			return nil, err
		}

		return db, nil
	})

	do.Provide(nil, func(i *do.Injector) (*router.Router, error) {
		db, err := do.Invoke[*sqlx.DB](i)
		if err != nil {
			return nil, err
		}

		cfg, err := do.Invoke[*config.Config](i)
		if err != nil {
			return nil, err
		}

		taskStore, err := do.Invoke[*tasks.Store](i)
		if err != nil {
			return nil, err
		}

		taskList, err := do.Invoke[*tasks.List](i)
		if err != nil {
			return nil, err
		}

		return router.NewRouter(db, cfg.HTTP.Port, taskStore, taskList), nil
	})

	do.Provide(nil, func(i *do.Injector) (*tasks.Store, error) {
		db, err := do.Invoke[*sqlx.DB](i)
		if err != nil {
			return nil, err
		}

		taskRepo := storage.NewTaskRepository(db)
		return tasks.NewStore(taskRepo, storage.NewUserRepository(db)), nil
	})

	do.Provide(nil, func(i *do.Injector) (*tasks.List, error) {
		db, err := do.Invoke[*sqlx.DB](i)
		if err != nil {
			return nil, err
		}

		taskRepo, err := do.Invoke[*storage.TaksRepository](i)
		if err != nil {
			return nil, err
		}

		return tasks.NewList(db, taskRepo), nil
	})

	do.Provide(nil, func(i *do.Injector) (*storage.TaksRepository, error) {
		db, err := do.Invoke[*sqlx.DB](i)
		if err != nil {
			return nil, err
		}

		return storage.NewTaskRepository(db), nil
	})
}
