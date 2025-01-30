package mongo

import databaseadapters "golang-web-core/srv/database_adapters"

func DefaultConfig() databaseadapters.ConnectionConfig {
	return databaseadapters.ConnectionConfig{
		Hostname: "localhost:27017",
		Database: "vatlnetwork",
	}
}
