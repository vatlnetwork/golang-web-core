package cfg

import (
	"golang-web-core/src/interface/database/mongo"
	media_os_int "golang-web-core/src/interface/media/os"
)

type Environment string

const (
	Prod     Environment = "Production"
	MockDev  Environment = "Mock Development"
	IntegDev Environment = "Integration Development"
)

type ServerConfig struct {
	Port  int
	Mongo mongo.Config
	SSL   SSLConfig
	Env   Environment
	Media media_os_int.Config
}

func (s ServerConfig) UsesSSL() bool {
	return s.SSL.CertFilePath != "" && s.SSL.CertKeyFilePath != ""
}

type SSLConfig struct {
	CertFilePath    string
	CertKeyFilePath string
}
