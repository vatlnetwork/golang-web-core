package cfg

import "golang-web-core/srv/database_adapters/mongo"

func Production() Config {
	config := Default()

	config.Environment = Prod

	config.Database.Adapter = mongo.NewMongoAdapter(false)
	config.Database.Connection = config.Database.Adapter.Connection()

	return config
}
