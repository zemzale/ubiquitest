package container

import (
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
	"github.com/zemzale/ubiquitest/config"
	"github.com/zemzale/ubiquitest/domain/tasks"
	"github.com/zemzale/ubiquitest/domain/users"
	"github.com/zemzale/ubiquitest/router"
	"github.com/zemzale/ubiquitest/storage"
	"github.com/zemzale/ubiquitest/ws"
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

		upsertUser, err := do.Invoke[*users.FindOrCreate](i)
		if err != nil {
			return nil, err
		}

		userFindByID, err := do.Invoke[*users.FindByID](i)
		if err != nil {
			return nil, err
		}

		wss, err := do.Invoke[*ws.Server](i)
		if err != nil {
			return nil, err
		}

		taskCalculate, err := do.Invoke[*tasks.CalculateCost](i)
		if err != nil {
			return nil, err
		}

		return router.NewRouter(cfg.HTTP.Port, taskStore, taskList, taskCalculate, upsertUser, userFindByID, wss), nil
	})

	do.Provide(nil, func(i *do.Injector) (*tasks.CalculateCost, error) {
		return tasks.NewCalculateCost(), nil
	})

	do.Provide(nil, func(i *do.Injector) (*tasks.Store, error) {
		db, err := do.Invoke[*sqlx.DB](i)
		if err != nil {
			return nil, err
		}

		taskRepo := storage.NewTaskRepository(db)
		updateParentCost, err := do.Invoke[*tasks.UpdateParentCost](i)
		if err != nil {
			return nil, err
		}

		return tasks.NewStore(updateParentCost, taskRepo, storage.NewUserRepository(db)), nil
	})

	do.Provide(nil, func(i *do.Injector) (*tasks.UpdateParentCost, error) {
		db, err := do.Invoke[*sqlx.DB](i)
		if err != nil {
			return nil, err
		}

		findAllParents, err := do.Invoke[*tasks.FindAllParents](i)
		if err != nil {
			return nil, err
		}

		return tasks.NewUpdateParentCost(findAllParents, storage.NewTaskRepository(db)), nil
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

	do.Provide(nil, func(i *do.Injector) (*users.FindOrCreate, error) {
		db, err := do.Invoke[*sqlx.DB](i)
		if err != nil {
			return nil, err
		}

		return users.NewFindOrCreate(db), nil
	})

	do.Provide(nil, func(i *do.Injector) (*users.FindByID, error) {
		userRepo, err := do.Invoke[*storage.UserRepository](i)
		if err != nil {
			return nil, err
		}

		return users.NewFindById(userRepo), nil
	})

	do.Provide(nil, func(i *do.Injector) (*ws.Server, error) {
		storeTask, err := do.Invoke[*tasks.Store](i)
		if err != nil {
			return nil, err
		}

		updateTask, err := do.Invoke[*tasks.Update](i)
		if err != nil {
			return nil, err
		}

		findUserByUsername, err := do.Invoke[*users.FindByUsername](i)
		if err != nil {
			return nil, err
		}

		taskCalculateCost, err := do.Invoke[*tasks.CalculateCost](i)
		if err != nil {
			return nil, err
		}
		taskFindAllParents, err := do.Invoke[*tasks.FindAllParents](i)
		if err != nil {
			return nil, err
		}

		return ws.NewServer(storeTask, updateTask, taskCalculateCost, taskFindAllParents, findUserByUsername), nil
	})

	do.Provide(nil, func(i *do.Injector) (*tasks.FindAllParents, error) {
		taskRepo, err := do.Invoke[*storage.TaksRepository](i)
		if err != nil {
			return nil, err
		}

		return tasks.NewFindAllParents(taskRepo), nil
	})

	do.Provide(nil, func(i *do.Injector) (*users.FindByUsername, error) {
		userRepo, err := do.Invoke[*storage.UserRepository](i)
		if err != nil {
			return nil, err
		}

		return users.NewFindByUsername(userRepo), nil
	})

	do.Provide(nil, func(i *do.Injector) (*tasks.Update, error) {
		db, err := do.Invoke[*sqlx.DB](i)
		if err != nil {
			return nil, err
		}

		return tasks.NewUpdate(db), nil
	})

	do.Provide(nil, func(i *do.Injector) (*storage.UserRepository, error) {
		db, err := do.Invoke[*sqlx.DB](i)
		if err != nil {
			return nil, err
		}

		return storage.NewUserRepository(db), nil
	})
}
