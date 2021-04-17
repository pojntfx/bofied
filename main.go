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

const (
	DHCPOptUUIDGUIDClientIdentifier = 97 // Option: (97) UUID/GUID-based Client Identifier
)

func main() {
	// Parse flags
	webDAVListenAddress := flag.String("webDAVListenAddress", "localhost:15256", "Listen address for the WebDAV server")
	httpListenAddress := flag.String("httpListenAddress", ":15257", "Listen address for the HTTP server")
	dhcpListenAddress := flag.String("dhcpListenAddress", ":67", "Listen address for the DHCP server")
	proxyDHCPListenAddress := flag.String("proxyDHCPListenAddress", ":4011", "Listen address for the proxyDHCP server")
	tftpListenAddress := flag.String("tftpListenAddress", ":69", "Listen address for the TFTP server")
	workingDir := flag.String("workingDir", ".", "Directory to store data in")
	advertisedIPFlag := flag.String("advertisedIP", "100.64.154.246", "IP address to advertise in DHCP and proxyDHCP")

	flag.Parse()

	// Process flags
	advertisedIP := net.ParseIP(*advertisedIPFlag).To4()
	broadcastAddress, err := net.ResolveUDPAddr("udp", "255.255.255.255:68")
	if err != nil {
		log.Fatal(err)
	}

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
				log.Fatal(err)
			}

			go func(r *net.UDPAddr, rawPacket []byte) {
				// Decode incoming packet
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

				// Process DHCP options
				isPXE := false
				clientIdentifierOpt := layers.DHCPOption{}
				if dhcpPacket.Operation == layers.DHCPOpRequest {
					for _, option := range dhcpPacket.Options {
						switch option.Type {
						case layers.DHCPOptClassID:
							isPXE, _, _, err = parsePXEClassIdentifier(string(option.Data))
							if err != nil {
								log.Fatal(err)
							}
						case DHCPOptUUIDGUIDClientIdentifier:
							clientIdentifierOpt = option
						}
					}
				}

				// Ignore non-PXE packets
				if !isPXE {
					return
				}

				// Create DHCP suboptions
				bootMenuDescription := "Boot iPXE (BIOS)"
				subOptions := []layers.DHCPOption{
					layers.NewDHCPOption(
						6,                        // Option 43 Suboption: (6) PXE discovery control
						[]byte{byte(0x00000003)}, // discovery control: 0x03, Disable Broadcast, Disable Multicast
					),
					layers.NewDHCPOption(
						10, // Option 43 Suboption: (10) PXE menu prompt
						append( // menu prompt: 00505845
							[]byte{
								0x00, // Timeout: 0
							},
							[]byte("PXE")..., // Prompt: PXE
						),
					),
					layers.NewDHCPOption(
						8, // Option 43 Suboption: (8) PXE boot servers
						append( // boot servers: 80000164409af2
							[]byte{
								0x80, 0x00, // Type: Unknown (32768)
								0x01, // IP count: 1
							},
							advertisedIP..., // IP: 100.64.154.246
						),
					),
					layers.NewDHCPOption(
						9, // Option 43 Suboption: (9) PXE boot menu
						append(
							[]byte{
								0x80, 0x00, // Type: Unknown (32768)
								byte(len(bootMenuDescription)), // Length: 16
							},
							[]byte(bootMenuDescription)..., // Description: Boot iPXE (BIOS)
						),
					),
				}

				// Serialize suboptions
				serializedSubOptions := []byte{}
				for _, subOption := range subOptions {
					serializedSubOptions = append(
						serializedSubOptions,
						append(
							[]byte{
								byte(subOption.Type),
								subOption.Length,
							},
							subOption.Data...,
						)...,
					)
				}

				// Serialize the outgoing packet
				outBuf := gopacket.NewSerializeBuffer()
				gopacket.SerializeLayers(
					outBuf,
					gopacket.SerializeOptions{
						FixLengths: true,
					},
					&layers.DHCPv4{
						Operation:    layers.DHCPOpReply,
						HardwareType: layers.LinkTypeEthernet,
						HardwareLen:  uint8(len(dhcpPacket.ClientHWAddr)),
						Xid:          dhcpPacket.Xid,
						ClientIP:     r.IP,
						YourClientIP: r.IP,
						NextServerIP: r.IP,
						RelayAgentIP: r.IP,
						ClientHWAddr: dhcpPacket.ClientHWAddr,
						Options: layers.DHCPOptions{
							layers.NewDHCPOption(
								layers.DHCPOptMessageType,
								[]byte{byte(layers.DHCPMsgTypeOffer)},
							),
							layers.NewDHCPOption(
								layers.DHCPOptServerID,
								advertisedIP,
							),
							layers.NewDHCPOption(
								layers.DHCPOptClassID,
								[]byte("PXEClient"),
							),
							clientIdentifierOpt,
							layers.NewDHCPOption(
								layers.DHCPOptVendorOption,
								serializedSubOptions,
							),
						},
					},
				)

				// Broadcast the packet
				n, err := conn.WriteToUDP(outBuf.Bytes(), broadcastAddress)
				if err != nil {
					log.Fatal(err)
				}

				log.Printf(`sent %v bytes of DHCP packets to client "%v"`, n, broadcastAddress)
			}(raddr, buf[:length])
		}
	}()

	go func() {
		// Create the server
		laddr, err := net.ResolveUDPAddr("udp", *proxyDHCPListenAddress)
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
				log.Fatal(err)
			}

			go func(r *net.UDPAddr, rawPacket []byte) {
				// Decode incoming packet
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

				// Process DHCP options
				isPXE := false
				clientIdentifierOpt := layers.DHCPOption{}
				if dhcpPacket.Operation == layers.DHCPOpRequest {
					for _, option := range dhcpPacket.Options {
						switch option.Type {
						case layers.DHCPOptClassID:
							isPXE, _, _, err = parsePXEClassIdentifier(string(option.Data))
							if err != nil {
								log.Fatal(err)
							}
						case DHCPOptUUIDGUIDClientIdentifier:
							clientIdentifierOpt = option
						}
					}
				}

				// Ignore non-PXE packets
				if !isPXE {
					return
				}

				// Create DHCP suboptions
				subOption :=
					layers.NewDHCPOption(
						71, // Option 43 Suboption: (71) PXE boot item
						[]byte{ // boot item: 80000000
							0x80, 0x00, // Type: 32768
							0x00, 0x00, // Layer: 0000
						},
					)

				// Serialize suboptions
				serializedSubOptions :=
					append(
						append(
							[]byte{
								byte(subOption.Type),
								subOption.Length,
							},
							subOption.Data...,
						),
						byte(0xff),
					)

				// Serialize the outgoing packet
				outBuf := gopacket.NewSerializeBuffer()
				bootFileName := "ipxe.kpxe"
				gopacket.SerializeLayers(
					outBuf,
					gopacket.SerializeOptions{
						FixLengths: true,
					},
					&layers.DHCPv4{
						Operation:    layers.DHCPOpReply,
						HardwareType: layers.LinkTypeEthernet,
						HardwareLen:  uint8(len(dhcpPacket.ClientHWAddr)),
						Xid:          dhcpPacket.Xid,
						ClientIP:     net.ParseIP("0.0.0.0").To4(),
						YourClientIP: r.IP,
						NextServerIP: advertisedIP,
						RelayAgentIP: net.ParseIP("0.0.0.0").To4(),
						ClientHWAddr: dhcpPacket.ClientHWAddr,
						File:         []byte(bootFileName),
						Options: layers.DHCPOptions{
							layers.NewDHCPOption(
								layers.DHCPOptMessageType,
								[]byte{byte(layers.DHCPMsgTypeAck)},
							),
							layers.NewDHCPOption(
								layers.DHCPOptServerID,
								advertisedIP,
							),
							layers.NewDHCPOption(
								layers.DHCPOptClassID,
								[]byte("PXEClient"),
							),
							clientIdentifierOpt,
							layers.NewDHCPOption(
								layers.DHCPOptVendorOption,
								serializedSubOptions,
							),
						},
					},
				)

				// Send the packet to the client
				n, err := conn.WriteToUDP(outBuf.Bytes(), r)
				if err != nil {
					log.Fatal(err)
				}

				log.Printf(`sent %v bytes of proxyDHCP packets to client "%v"`, n, r)
			}(raddr, buf[:length])
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
