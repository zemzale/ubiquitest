package container

import (
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
	"github.com/zemzale/ubiquitest/config"
	"github.com/zemzale/ubiquitest/router"
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

		return router.NewRouter(db, cfg.HTTP.Port), nil
	})
}
