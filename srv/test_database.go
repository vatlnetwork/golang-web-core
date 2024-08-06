package srv

import (
	"golang-web-core/util"
	"log"
)

func (s *Server) TestDatabase() error {
	log.Println("Testing database connection...")
	err := s.Config.Database.Adapter.TestConnection()
	if err != nil {
		return err
	}
	util.LogColor("lightgreen", "Database connection successful")

	return nil
}
