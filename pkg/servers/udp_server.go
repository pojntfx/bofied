package servers

import (
	"errors"
	"log"
	"net"

	"github.com/pojntfx/bofied/pkg/utils"
)

const (
	UDPServerReadBufSize = 1024
)

type UDPServer struct {
	listenAddress string
	advertisedIP  string
	handlePacket  func(conn *net.UDPConn, raddr *net.UDPAddr, braddr *net.UDPAddr, rawIncomingUDPPacket []byte) (int, error)
}

func (s *UDPServer) ListenAndServe() error {
	// Parse the addresses
	laddr, err := net.ResolveUDPAddr("udp", s.listenAddress)
	if err != nil {
		return err
	}

	aaddr, err := net.ResolveIPAddr("ip", s.advertisedIP)
	if err != nil {
		return err
	}

	// Get all interfaces
	ifaces, err := net.Interfaces()
	if err != nil {
		return err
	}

	// Find the broadcast address of the interface with the advertised IP
	broadcastIP := ""
ifaceLoop:
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return err
		}

		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}

			if ipnet.IP.Equal(aaddr.IP) {
				broadcastIP, err = utils.GetBroadcastAddress(ipnet)
				if err != nil {
					return err
				}

				break ifaceLoop
			}
		}
	}

	// Return if no interface with the advertised IP could be found
	if broadcastIP == "" {
		return errors.New("could not resolve broadcast IP")
	}

	// Construct the broadcast address
	braddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(broadcastIP, "68"))
	if err != nil {
		return err
	}

	// Listen
	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		return err
	}

	// Loop over packets
	for {
		// Read packet into buffer
		buf := make([]byte, UDPServerReadBufSize)
		length, raddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			return err
		}

		// Handle the read packet
		go func() {
			if _, err := s.handlePacket(conn, raddr, braddr, buf[:length]); err != nil {
				log.Println("could not handle packet:", err)
			}
		}()
	}
}
