package cfg

import (
	databaseadapters "golang-web-core/srv/database_adapters"
	"golang-web-core/srv/database_adapters/imdb"
)

func Default() Config {
	dbAdapter := imdb.NewImdbAdapter()

	return Config{
		Port:        3000,
		PublicFS:    true,
		Environment: Dev,
		Database: databaseadapters.DatabaseConfig{
			Adapter:    dbAdapter,
			Connection: dbAdapter.ConnectionConfig,
		},
	}
}
