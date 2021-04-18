package servers

import (
	"log"
	"net"
)

type UDPServer struct {
	listenAddress string
	handlePacket  func(conn *net.UDPConn, raddr *net.UDPAddr, braddr *net.UDPAddr, rawIncomingUDPPacket []byte) (int, error)
}

func (s *UDPServer) ListenAndServe() error {
	// Parse the addresses
	laddr, err := net.ResolveUDPAddr("udp", s.listenAddress)
	if err != nil {
		return err
	}
	braddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:68")
	if err != nil {
		log.Fatal(err)
	}

	// Listen
	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		return err
	}

	// Loop over packets
	for {
		// Read packet into buffer
		buf := make([]byte, DHCPServerReadBufSize)
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
