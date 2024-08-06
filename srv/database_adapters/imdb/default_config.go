package imdb

import databaseadapters "golang-web-core/srv/database_adapters"

func DefaultConfig() databaseadapters.ConnectionConfig {
	return databaseadapters.ConnectionConfig{
		Hostname: "hostname",
		Database: "database",
	}
}
