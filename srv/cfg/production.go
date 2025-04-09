package cfg

import "golang-web-core/srv/database_adapters/mongo"

func Production() Config {
	config := Default()

	config.Environment = Prod

	config.Database.Connection = mongo.DefaultConfig()
	config.Database.Adapter = mongo.NewMongoAdapter(config.Database.Connection, false)

	return config
}
