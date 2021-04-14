package main

import (
	"errors"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/pin/tftp"
	"golang.org/x/net/webdav"
)

func main() {
	// Parse flags
	webDAVListenAddress := flag.String("webDAVListenAddress", "localhost:15256", "Listen address for the WebDAV server")
	tftpListenAddress := flag.String("tftpListenAddress", ":69", "Listen address for the TFTP server")
	workingDir := flag.String("workingDir", ".", "Directory to store data in")

	flag.Parse()

	// Create servers
	webdavSrv := &webdav.Handler{
		FileSystem: webdav.Dir(*workingDir),
		LockSystem: webdav.NewMemLS(),
	}

	tftpSrv := tftp.NewServer(
		func(filename string, rf io.ReaderFrom) error {
			// Prevent accessing any parent directories
			if strings.Contains(filename, "..") {
				msg := errors.New("blocked request by client to access parent directory" + filename)

				log.Println(msg)

				return msg
			}

			// Open file to send
			fullFilename := filepath.Join(*workingDir, filename)
			file, err := os.Open(fullFilename)
			if err != nil {
				log.Println("could not open file", err)

				return err
			}

			// Send the file to the client
			n, err := rf.ReadFrom(file)
			log.Printf("sent %v (%v bytes)", fullFilename, n)

			return err
		},
		nil,
	)

	// Start servers
	http.Handle("/", webdavSrv)
	go func() {
		if err := http.ListenAndServe(*webDAVListenAddress, webdavSrv); err != nil {
			log.Fatal(err)
		}
	}()

	if err := tftpSrv.ListenAndServe(*tftpListenAddress); err != nil {
		log.Fatal(err)
	}
}
