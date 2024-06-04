package mongo

import "fmt"

type Config struct {
	Host      string
	Port      int
	Username  string
	Password  string
	DbName    string
	CollNames CollectionNames
}

func (c Config) ConnectionString() string {
	authAndHost := c.Host
	if c.UsesAuth() {
		authAndHost = fmt.Sprintf("%v:%v@%v", c.Username, c.Password, c.Host)
	}
	return fmt.Sprintf("mongodb://%v:%v/%v", authAndHost, c.Port, c.DbName)
}

func (c Config) UsesAuth() bool {
	return c.Username != "" && c.Password != ""
}

type CollectionNames struct {
	Users      string
	Sessions   string
	MediaFiles string
}
