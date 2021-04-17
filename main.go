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
	"strconv"
	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
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
			length, _, err := conn.ReadFromUDP(buf)
			if err != nil {
				log.Fatal(err)
			}

			go func(rawPacket []byte) {
				packet := gopacket.NewPacket(rawPacket, layers.LayerTypeDHCPv4, gopacket.Default)

				dhcpLayer := packet.Layer(layers.LayerTypeDHCPv4)
				if dhcpLayer == nil {
					log.Fatal("received a non-DHCP layer")

					return
				}

				dhcpPacket, ok := dhcpLayer.(*layers.DHCPv4)
				if !ok {
					log.Fatal("received invalid UDP layer")
				}

				if dhcpPacket.Operation == layers.DHCPOpRequest {
					for _, option := range dhcpPacket.Options {
						if option.Type == layers.DHCPOptClassID {
							isPXE, arch, undi, err := parsePXEClassIdentifier(string(option.Data))
							if err != nil {
								log.Fatal(err)
							}

							log.Println(isPXE, arch, undi)

							break
						}
					}
				}
			}(buf[:length])
		}
	}()

	if err := tftpSrv.ListenAndServe(*tftpListenAddress); err != nil {
		log.Fatal(err)
	}
}

func parsePXEClassIdentifier(classID string) (isPXE bool, arch int, undi int, err error) {
	parts := strings.Split(classID, ":")

	for i, part := range parts {
		switch part {
		case "PXEClient":
			isPXE = true
		case "Arch":
			if len(parts) > i {
				arch, err = strconv.Atoi(parts[i+1])
				if err != nil {
					return
				}
			}
		case "UNDI":
			if len(parts) > i {
				undi, err = strconv.Atoi(parts[i+1])
				if err != nil {
					return
				}
			}
		}
	}

	return
}
