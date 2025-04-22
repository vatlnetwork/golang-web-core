package cfg

func Default() Config {
	return Config{
		Port:     3000,
		PublicFS: true,
		Env:      Development,
	}
}
