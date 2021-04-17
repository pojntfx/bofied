package main

import (
	"flag"
	"log"

	"github.com/pojntfx/bofied/pkg/servers"
)

func main() {
	// Parse flags
	workingDir := flag.String("workingDir", ".", "Working directory")
	webDAVListenAddress := flag.String("webDAVListenAddress", ":15256", "Listen address for WebDAV server")

	flag.Parse()

	// Create servers
	webDAVServer := servers.NewWebDAVServer(*workingDir, *webDAVListenAddress)

	// Start servers
	log.Fatal(webDAVServer.ListenAndServe())
}
