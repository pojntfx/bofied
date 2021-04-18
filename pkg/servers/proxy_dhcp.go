package servers

import (
	"errors"
	"log"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/pojntfx/bofied/pkg/config"
	"github.com/pojntfx/bofied/transcoding"
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
	// Decode and parse packet
	incomingUDPPacket := gopacket.NewPacket(rawIncomingUDPPacket, layers.LayerTypeDHCPv4, gopacket.Default)

	dhcpLayer := incomingUDPPacket.Layer(layers.LayerTypeDHCPv4)
	if dhcpLayer == nil {
		return 0, errors.New("could not parse DHCP layer: not a DHCP layer")
	}

	incomingDHCPPacket, ok := dhcpLayer.(*layers.DHCPv4)
	if !ok {
		return 0, errors.New("could not parse DHCP layer: invalid DHCP layer")
	}

	// Parse DHCP options
	isPXE := false
	arch := 0
	undi := 0
	clientIdentifierOpt := layers.DHCPOption{}
	if incomingDHCPPacket.Operation == layers.DHCPOpRequest {
		for _, option := range incomingDHCPPacket.Options {
			switch option.Type {
			case layers.DHCPOptClassID:
				pxe, a, u, err := transcoding.ParsePXEClassIdentifier(string(option.Data))
				if err != nil {
					return 0, err
				}

				isPXE = pxe
				arch = a
				undi = u
			case DHCPOptUUIDGUIDClientIdentifier:
				clientIdentifierOpt = option
			}
		}
	}

	// Ignore non-PXE packets
	if !isPXE {
		return 0, nil
	}

	// Create DHCP Option 43 suboptions
	subOption :=
		layers.NewDHCPOption(
			71, // Option 43 Suboption: (71) PXE boot item
			[]byte{ // boot item: 80000000
				0x80, 0x00, // Type: 32768
				0x00, 0x00, // Layer: 0000
			},
		)

		// Serialize DHCP Option 43 suboptions
	serializedSubOptions :=
		append(
			append(
				[]byte{
					byte(subOption.Type),
					subOption.Length,
				},
				subOption.Data...,
			),
			byte(0xff), // PXE Client End: 255
		)

		// Get the boot file name
	bootFileName, err := config.GetFileName(
		ConfigFunctionIdentifier,
		configFileLocation,
		raddr.IP.String(),
		incomingDHCPPacket.ClientHWAddr.String(),
		arch,
		undi,
	)
	if err != nil {
		return 0, err
	}

	// Create the outgoing packet
	outgoingDHCPPacket := &layers.DHCPv4{
		Operation:    layers.DHCPOpReply,
		HardwareType: layers.LinkTypeEthernet,
		HardwareLen:  uint8(len(incomingDHCPPacket.ClientHWAddr)),
		Xid:          incomingDHCPPacket.Xid,
		ClientIP:     net.ParseIP("0.0.0.0").To4(),
		YourClientIP: raddr.IP,
		NextServerIP: advertisedIP,
		RelayAgentIP: net.ParseIP("0.0.0.0").To4(),
		ClientHWAddr: incomingDHCPPacket.ClientHWAddr,
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
				[]byte(DHCPBaseClassID),
			),
			clientIdentifierOpt,
			layers.NewDHCPOption(
				layers.DHCPOptVendorOption,
				serializedSubOptions,
			),
		},
	}

	// Serialize the outgoing packet
	buf := gopacket.NewSerializeBuffer()
	gopacket.SerializeLayers(
		buf,
		gopacket.SerializeOptions{
			FixLengths: true,
		},
		outgoingDHCPPacket,
	)

	log.Printf(`sending %v bytes of proxyDHCP packets to client "%v"`, len(buf.Bytes()), raddr)

	// Send the packet to the client
	return conn.WriteToUDP(buf.Bytes(), raddr)
}
