package container

import (
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
	"github.com/zemzale/ubiquitest/config"
)

func Load() {
	do.Provide(nil, func(i *do.Injector) (*config.Config, error) {
		return config.Load(), nil
	})

	do.Provide(nil, func(i *do.Injector) (*sqlx.DB, error) {
		config, err := do.Invoke[*config.Config](i)
		if err != nil {
			return nil, err
		}

		db, err := sqlx.Open(config.DB.Driver, config.DB.DSN)
		if err != nil {
			return nil, err
		}

		return db, nil
	})
}
