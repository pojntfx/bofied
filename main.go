package main

import (
	"flag"
	"log"
	"net/http"

	"golang.org/x/net/webdav"
)

func main() {
	listenAddress := flag.String("listenAddress", "localhost:15256", "Listen address")
	workingDir := flag.String("workingDir", ".", "Directory to store data in")

	flag.Parse()

	srv := &webdav.Handler{
		FileSystem: webdav.Dir(*workingDir),
		LockSystem: webdav.NewMemLS(),
	}

	http.Handle("/", srv)

	if err := http.ListenAndServe(*listenAddress, srv); err != nil {
		log.Fatal(err)
	}
}
