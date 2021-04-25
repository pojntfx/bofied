package servers

import (
	"net"

	"github.com/pojntfx/bofied/pkg/eventing"
	"github.com/pojntfx/bofied/pkg/transcoding"
)

const (
	DHCPServerBootMenuPromptBIOS      = "PXE"
	DHCPServerBootMenuDescriptionBIOS = "Boot from bofied (BIOS)"
)

type DHCPServer struct {
	UDPServer
}

func NewDHCPServer(listenAddress string, advertisedIP string, eventHandler *eventing.EventHandler) *DHCPServer {
	return &DHCPServer{
		UDPServer: UDPServer{
			listenAddress: listenAddress,
			handlePacket: func(conn *net.UDPConn, _ *net.UDPAddr, braddr *net.UDPAddr, rawIncomingUDPPacket []byte) (int, error) {
				return handleDHCPPacket(conn, braddr, rawIncomingUDPPacket, net.ParseIP(advertisedIP).To4(), eventHandler.Emit)
			},
		},
	}
}

func handleDHCPPacket(conn *net.UDPConn, braddr *net.UDPAddr, rawIncomingUDPPacket []byte, advertisedIP net.IP, emit func(f string, v ...interface{})) (int, error) {
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

	emit(`sending %v bytes of DHCP packets to client "%v"`, len(outgoingDHCPPacket), braddr)

	// Broadcast the packet
	return conn.WriteToUDP(outgoingDHCPPacket, braddr)
}
