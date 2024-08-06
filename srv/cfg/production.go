package cfg

import "golang-web-core/srv/database_adapters/mongo"

func Production() Config {
	config := Default()

	config.Database.Adapter = mongo.NewMongoAdapter()
	config.Database.Connection = config.Database.Adapter.Connection()

	return config
}
