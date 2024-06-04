package mongo

const TestDbName string = "vatlnetworkWebTest"
const DevDbName string = "vatlnetworkWebDev"
const ProdDbName string = "vatlnetworkWeb"

func DefaultConfig() Config {
	return Config{
		Host:      "localhost",
		Port:      27017,
		DbName:    DevDbName,
		CollNames: DefaultCollNames(),
	}
}

func TestConfig() Config {
	cfg := DefaultConfig()
	cfg.DbName = TestDbName
	return cfg
}

func DefaultCollNames() CollectionNames {
	return CollectionNames{
		Users:      "users",
		Sessions:   "sessions",
		MediaFiles: "mediaFiles",
	}
}
