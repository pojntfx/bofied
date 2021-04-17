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
	httpListenAddress := flag.String("httpListenAddress", ":15257", "Listen address for HTTP server")

	flag.Parse()

	// Create servers
	webDAVServer := servers.NewWebDAVServer(*workingDir, *webDAVListenAddress)
	httpServer := servers.NewHTTPServer(*workingDir, *httpListenAddress)

	// Start servers
	go func() {
		log.Fatal(httpServer.ListenAndServe())
	}()

	log.Fatal(webDAVServer.ListenAndServe())
}
