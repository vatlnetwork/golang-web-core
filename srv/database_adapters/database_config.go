package databaseadapters

type DatabaseConfig struct {
	Adapter    DatabaseAdapter
	Connection ConnectionConfig
}

func (c DatabaseConfig) UsingDatabase() bool {
	return c.Adapter != nil && c.Connection.Hostname != "" && c.Connection.Database != ""
}

func (c DatabaseConfig) UsingAuth() bool {
	return c.UsingDatabase() && c.Connection.UsingAuth()
}
