package srv

import (
	"fmt"
)

func PrintServerConfig(server *Server) {
	c := server.Config

	fmt.Println("Server Config:")
	fmt.Printf("   Port: %v\n", c.Port)
	fmt.Printf("   Using SSL: %v\n", c.IsSSL())
	if c.IsSSL() {
		fmt.Printf("      Cert Path: %v\n", c.SSL.CertPath)
		fmt.Printf("      Key Path: %v\n", c.SSL.KeyPath)
	}
	fmt.Printf("   Public FS Enabled: %v\n", c.PublicFS)
	fmt.Printf("   Environment: %v\n", c.Environment)
	fmt.Printf("   Database Enabled: %v\n", c.Database.UsingDatabase())
	if c.Database.UsingDatabase() {
		fmt.Printf("      Adapter: %v\n", c.Database.Adapter.Name())
		fmt.Printf("      Hostname: %v\n", c.Database.Connection.Hostname)
		fmt.Printf("      Database: %v\n", c.Database.Connection.Database)
	}
	fmt.Printf("   Database Auth: %v\n", c.Database.UsingAuth())
	fmt.Printf("   # of Routes: %v\n", len(server.Routes))
	fmt.Println("")
}
