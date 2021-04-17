package main

import (
	"flag"
	"log"

	"github.com/pojntfx/bofied/pkg/servers"
)

func main() {
	// Parse flags
	workingDir := flag.String("workingDir", ".", "Working directory")
	advertisedIP := flag.String("advertisedIP", "100.64.154.246", "IP to advertise for DHCP clients")
	webDAVListenAddress := flag.String("webDAVListenAddress", ":15256", "Listen address for WebDAV server")
	httpListenAddress := flag.String("httpListenAddress", ":15257", "Listen address for HTTP server")
	dhcpListenAddress := flag.String("dhcpListenAddress", ":67", "Listen address for DHCP server")

	flag.Parse()

	// Create servers
	webDAVServer := servers.NewWebDAVServer(*workingDir, *webDAVListenAddress)
	httpServer := servers.NewHTTPServer(*workingDir, *httpListenAddress)
	dhcpServer := servers.NewDHCPServer(*dhcpListenAddress, *advertisedIP)

	// Start servers
	go func() {
		log.Fatal(dhcpServer.ListenAndServe())
	}()

	go func() {
		log.Fatal(httpServer.ListenAndServe())
	}()

	log.Fatal(webDAVServer.ListenAndServe())
}
