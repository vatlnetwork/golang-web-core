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
	fmt.Printf("   # of Routes: %v\n", len(server.Routes))
	fmt.Println("")
}
