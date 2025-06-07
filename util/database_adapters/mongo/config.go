package mongo

import (
	"net/url"
)

type Config struct {
	Hostname string `json:"hostname"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c Config) ConnectionString() string {
	uri := url.URL{
		Scheme: "mongodb",
		Host:   c.Hostname,
		Path:   c.Database,
	}

	if c.UsingAuth() {
		uri.User = url.UserPassword(c.Username, c.Password)
	}

	return uri.String()
}

func (c Config) IsEnabled() bool {
	return c.Hostname != "" && c.Database != ""
}

func (c Config) UsingAuth() bool {
	return c.Username != "" && c.Password != ""
}
