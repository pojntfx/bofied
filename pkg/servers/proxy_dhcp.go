package servers

import (
	"log"
	"net"

	"github.com/pojntfx/bofied/pkg/config"
	"github.com/pojntfx/bofied/pkg/transcoding"
)

const (
	ConfigFunctionIdentifier = "config.GetFileName"
)

type ProxyDHCPServer struct {
	UDPServer
}

func NewProxyDHCPServer(listenAddress string, advertisedIP string, configFileLocation string) *ProxyDHCPServer {
	return &ProxyDHCPServer{
		UDPServer: UDPServer{
			listenAddress: listenAddress,
			handlePacket: func(conn *net.UDPConn, raddr *net.UDPAddr, braddr *net.UDPAddr, rawIncomingUDPPacket []byte) (int, error) {
				return handleProxyDHCPPacket(conn, raddr, braddr, rawIncomingUDPPacket, net.ParseIP(advertisedIP).To4(), configFileLocation)
			},
		},
	}
}

func handleProxyDHCPPacket(conn *net.UDPConn, raddr *net.UDPAddr, _ *net.UDPAddr, rawIncomingUDPPacket []byte, advertisedIP net.IP, configFileLocation string) (int, error) {
	// Decode packet
	incomingDHCPPacket, err := transcoding.DecodeDHCPPacket(rawIncomingUDPPacket)
	if err != nil {
		return 0, err
	}

	// Ignore non-PXE packets
	if !incomingDHCPPacket.IsPXE {
		return 0, nil
	}

	// Get the boot file name
	bootFileName, err := config.GetFileName(
		ConfigFunctionIdentifier,
		configFileLocation,
		raddr.IP.String(),
		incomingDHCPPacket.ClientHWAddr.String(),
		incomingDHCPPacket.Arch,
		incomingDHCPPacket.Undi,
	)
	if err != nil {
		return 0, err
	}

	// Encode packet
	outgoingDHCPPacket := transcoding.EncodeProxyDHCPPacket(
		incomingDHCPPacket.ClientHWAddr,
		incomingDHCPPacket.Xid,
		advertisedIP,
		incomingDHCPPacket.ClientIdentifierOpt,
		raddr.IP.To4(),
		bootFileName,
	)

	log.Printf(`sending %v bytes of proxyDHCP packets to client "%v"`, len(outgoingDHCPPacket), raddr)

	// Send the packet to the client
	return conn.WriteToUDP(outgoingDHCPPacket, raddr)
}
