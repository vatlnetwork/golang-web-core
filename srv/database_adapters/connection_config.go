package databaseadapters

type ConnectionConfig struct {
	Hostname string
	Database string
	Username string
	Password string
}

func (c ConnectionConfig) UsingAuth() bool {
	return c.Username != "" && c.Password != ""
}
