package cfg

type Config struct {
	Port int
	SSL  SSL
}

func (c Config) IsSSL() bool {
	return c.SSL.CertPath != "" && c.SSL.KeyPath != ""
}

type SSL struct {
	CertPath string
	KeyPath  string
}
