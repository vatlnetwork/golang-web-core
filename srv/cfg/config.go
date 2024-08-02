package cfg

import (
	"fmt"
	"os"
)

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

func (s *SSL) SetCertPath(path string) error {
	_, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("the cert path you specified (%v) does not exist", path)
		}
		return err
	}

	s.CertPath = path

	return nil
}

func (s *SSL) SetKeyPath(path string) error {
	_, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("the cert path you specified (%v) does not exist", path)
		}
		return err
	}

	s.KeyPath = path

	return nil
}
