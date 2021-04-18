package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/pojntfx/bofied/pkg/config"
	"github.com/pojntfx/bofied/pkg/servers"
)

func main() {
	// Get default working dir
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("could not get home directory", err)
	}
	workingDirDefault := filepath.Join(home, ".local", "share", "bofied", "var", "lib", "bofied")

	// Parse flags
	workingDir := flag.String("workingDir", workingDirDefault, "Working directory")
	configFileName := flag.String("configFileName", "config.go", "Name of the config file (in the working directory)")
	advertisedIP := flag.String("advertisedIP", "100.64.154.246", "IP to advertise for DHCP clients")
	dhcpListenAddress := flag.String("dhcpListenAddress", ":67", "Listen address for DHCP server")
	proxyDHCPListenAddress := flag.String("proxyDHCPListenAddress", ":4011", "Listen address for proxyDHCP server")
	tftpListenAddress := flag.String("tftpListenAddress", ":69", "Listen address for TFTP server")
	webDAVListenAddress := flag.String("webDAVListenAddress", ":15256", "Listen address for WebDAV server")
	httpListenAddress := flag.String("httpListenAddress", ":15257", "Listen address for HTTP server")

	flag.Parse()

	// Initialize the working directory
	if err := config.CreateConfigIfNotExists(filepath.Join(*workingDir, *configFileName)); err != nil {
		log.Fatal(err)
	}

	// Create servers
	dhcpServer := servers.NewDHCPServer(*dhcpListenAddress, *advertisedIP)
	proxyDHCPServer := servers.NewProxyDHCPServer(
		*proxyDHCPListenAddress,
		*advertisedIP,
		filepath.Join(*workingDir, *configFileName),
	)
	tftpServer := servers.NewTFTPServer(*workingDir, *tftpListenAddress)
	webDAVServer := servers.NewWebDAVServer(*workingDir, *webDAVListenAddress)
	httpServer := servers.NewHTTPServer(*workingDir, *httpListenAddress)

	// Start servers
	go func() {
		log.Fatal(dhcpServer.ListenAndServe())
	}()

	go func() {
		log.Fatal(proxyDHCPServer.ListenAndServe())
	}()

	go func() {
		log.Fatal(tftpServer.ListenAndServe())
	}()

	go func() {
		log.Fatal(webDAVServer.ListenAndServe())
	}()

	log.Fatal(httpServer.ListenAndServe())
}
