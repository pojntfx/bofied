package servers

import (
	"net"

	"github.com/pojntfx/bofied/pkg/config"
	"github.com/pojntfx/bofied/pkg/eventing"
	"github.com/pojntfx/bofied/pkg/transcoding"
)

type ProxyDHCPServer struct {
	UDPServer
}

func NewProxyDHCPServer(listenAddress string, advertisedIP string, configFileLocation string, eventHandler *eventing.EventHandler, pureConfig bool) *ProxyDHCPServer {
	return &ProxyDHCPServer{
		UDPServer: UDPServer{
			listenAddress: listenAddress,
			advertisedIP:  advertisedIP,
			handlePacket: func(conn *net.UDPConn, raddr *net.UDPAddr, braddr *net.UDPAddr, rawIncomingUDPPacket []byte) (int, error) {
				return handleProxyDHCPPacket(conn, raddr, braddr, rawIncomingUDPPacket, net.ParseIP(advertisedIP).To4(), configFileLocation, eventHandler.Emit, pureConfig)
			},
		},
	}
}

func handleProxyDHCPPacket(conn *net.UDPConn, raddr *net.UDPAddr, _ *net.UDPAddr, rawIncomingUDPPacket []byte, advertisedIP net.IP, configFileLocation string, emit func(f string, v ...interface{}), pureConfig bool) (int, error) {
	// Decode packet
	incomingDHCPPacket, err := transcoding.DecodeDHCPPacket(rawIncomingUDPPacket)
	if err != nil {
		return 0, err
	}

	// Ignore non-PXE packets
	if !incomingDHCPPacket.IsPXE {
		return 0, nil
	}

	emit(
		`handling proxyDHCP for client with IP %v, MAC %v, architecture %v and architecture ID %v`,
		raddr.IP.String(),
		incomingDHCPPacket.ClientHWAddr.String(),
		config.GetNameForArchId(incomingDHCPPacket.Arch),
		incomingDHCPPacket.Arch,
	)

	// Get the boot file name
	bootFileName, err := config.GetFileName(
		configFileLocation,
		raddr.IP.String(),
		incomingDHCPPacket.ClientHWAddr.String(),
		config.GetNameForArchId(incomingDHCPPacket.Arch),
		incomingDHCPPacket.Arch,
		pureConfig,
		func(s string) {
			if s != "" {
				emit("from config: %v", s)
			}
		},
	)
	if err != nil {
		emit("could not process config: %v", err)

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

	emit(`sending %v bytes of proxyDHCP packets to client "%v"`, len(outgoingDHCPPacket), raddr)

	// Send the packet to the client
	return conn.WriteToUDP(outgoingDHCPPacket, raddr)
}
