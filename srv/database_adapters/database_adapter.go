package databaseadapters

type DatabaseAdapter interface {
	Name() string
	Connection() ConnectionConfig
	TestConnection() error
}
