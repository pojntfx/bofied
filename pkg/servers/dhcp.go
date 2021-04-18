package servers

import (
	"log"
	"net"

	"github.com/pojntfx/bofied/transcoding"
)

const (
	DHCPServerBootMenuPromptBIOS      = "PXE"
	DHCPServerBootMenuDescriptionBIOS = "Boot from bofied (BIOS)"
)

type DHCPServer struct {
	UDPServer
}

func NewDHCPServer(listenAddress string, advertisedIP string) *DHCPServer {
	return &DHCPServer{
		UDPServer: UDPServer{
			listenAddress: listenAddress,
			handlePacket: func(conn *net.UDPConn, _ *net.UDPAddr, braddr *net.UDPAddr, rawIncomingUDPPacket []byte) (int, error) {
				return handleDHCPPacket(conn, braddr, rawIncomingUDPPacket, net.ParseIP(advertisedIP).To4())
			},
		},
	}
}

func handleDHCPPacket(conn *net.UDPConn, braddr *net.UDPAddr, rawIncomingUDPPacket []byte, advertisedIP net.IP) (int, error) {
	// Decode packet
	incomingDHCPPacket, err := transcoding.DecodeDHCPPacket(rawIncomingUDPPacket)
	if err != nil {
		return 0, err
	}

	// Ignore non-PXE packets
	if !incomingDHCPPacket.IsPXE {
		return 0, nil
	}

	// Encode packet
	outgoingDHCPPacket := transcoding.EncodeDHCPPacket(
		incomingDHCPPacket.ClientHWAddr,
		incomingDHCPPacket.Xid,
		advertisedIP,
		incomingDHCPPacket.ClientIdentifierOpt,
		incomingDHCPPacket.Arch,
		DHCPServerBootMenuPromptBIOS,
		DHCPServerBootMenuDescriptionBIOS,
	)

	log.Printf(`sending %v bytes of DHCP packets to client "%v"`, len(outgoingDHCPPacket), braddr)

	// Broadcast the packet
	return conn.WriteToUDP(outgoingDHCPPacket, braddr)
}
