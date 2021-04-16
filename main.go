package main

import (
	"errors"
	"flag"
	"io"
	"log"
	"net"
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
	httpListenAddress := flag.String("httpListenAddress", ":15257", "Listen address for the HTTP server")
	dhcpListenAddress := flag.String("dhcpListenAddress", ":67", "Listen address for the DHCP server")
	tftpListenAddress := flag.String("tftpListenAddress", ":69", "Listen address for the TFTP server")
	workingDir := flag.String("workingDir", ".", "Directory to store data in")

	flag.Parse()

	// Create servers
	webdavSrv := &webdav.Handler{
		FileSystem: webdav.Dir(*workingDir),
		LockSystem: webdav.NewMemLS(),
	}

	httpSrv := http.FileServer(http.Dir(*workingDir))

	tftpSrv := tftp.NewServer(
		func(filename string, rf io.ReaderFrom) error {
			// Get remote IP
			raddr := rf.(tftp.OutgoingTransfer).RemoteAddr()

			// Prevent accessing any parent directories
			fullFilename := filepath.Join(*workingDir, filename)
			if strings.Contains(filename, "..") {
				log.Printf(`could not send file: get request to file "%v" by client "%v" blocked because it is located outside the working directory "%v"`, fullFilename, raddr.String(), *workingDir)

				return errors.New("unauthorized: tried to access file outside working directory")
			}

			// Open file to send
			file, err := os.Open(fullFilename)
			if err != nil {
				log.Printf(`could not open file "%v" for client "%v": %v`, fullFilename, raddr.String(), err)

				return err
			}

			// Send the file to the client
			n, err := rf.ReadFrom(file)
			log.Printf(`sent file "%v" (%v bytes) to client "%v"`, fullFilename, n, raddr.String())

			return err
		},
		nil,
	)

	// Start servers
	go func() {
		if err := http.ListenAndServe(*webDAVListenAddress, webdavSrv); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		if err := http.ListenAndServe(*httpListenAddress, httpSrv); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		// Create the server
		laddr, err := net.ResolveUDPAddr("udp", *dhcpListenAddress)
		if err != nil {
			log.Fatal(err)
		}

		conn, err := net.ListenUDP("udp", laddr)
		if err != nil {
			log.Fatal(err)
		}

		// Read UDP datagrams
		for {
			buf := make([]byte, 1024)
			length, raddr, err := conn.ReadFromUDP(buf)
			if err != nil {
				log.Printf(`could not read UDP datagram from client "%v": %v`, raddr.String(), err)

				continue
			}

			go func(b []byte, n int) {
				packet := b[:n]

				log.Println(packet, raddr)
			}(buf, length)
		}
	}()

	if err := tftpSrv.ListenAndServe(*tftpListenAddress); err != nil {
		log.Fatal(err)
	}
}
