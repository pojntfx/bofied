package servers

import (
	"errors"
	"log"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/pojntfx/bofied/transcoding"
)

const (
	DHCPServerReadBufSize             = 1024
	DHCPOptUUIDGUIDClientIdentifier   = 97 // Option: (97) UUID/GUID-based Client Identifier
	DHCPEFIArch                       = 7  // X86-64_EFI
	DHCPBaseClassID                   = "PXEClient"
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
	clientIdentifierOpt := layers.DHCPOption{}
	if incomingDHCPPacket.Operation == layers.DHCPOpRequest {
		for _, option := range incomingDHCPPacket.Options {
			switch option.Type {
			case layers.DHCPOptClassID:
				pxe, a, _, err := transcoding.ParsePXEClassIdentifier(string(option.Data))
				if err != nil {
					return 0, err
				}

				isPXE = pxe
				arch = a
			case DHCPOptUUIDGUIDClientIdentifier:
				clientIdentifierOpt = option
			}
		}
	}

	// Ignore non-PXE packets
	if !isPXE {
		return 0, nil
	}

	// Create the outgoing packet
	outgoingDHCPPacket := &layers.DHCPv4{
		Operation:    layers.DHCPOpReply,
		HardwareType: layers.LinkTypeEthernet,
		HardwareLen:  uint8(len(incomingDHCPPacket.ClientHWAddr)),
		Xid:          incomingDHCPPacket.Xid,
		ClientIP:     net.ParseIP("0.0.0.0").To4(),
		YourClientIP: net.ParseIP("0.0.0.0").To4(),
		NextServerIP: advertisedIP,
		RelayAgentIP: net.ParseIP("0.0.0.0").To4(),
		ClientHWAddr: incomingDHCPPacket.ClientHWAddr,
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
				[]byte(DHCPBaseClassID),
			),
			clientIdentifierOpt,
		},
	}

	// If the packet is not intended for EFI systems, add additional required options
	if arch != DHCPEFIArch {
		// Create DHCP Option 43 suboptions
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
					[]byte(DHCPServerBootMenuPromptBIOS)..., // Prompt: PXE
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
						byte(len(DHCPServerBootMenuDescriptionBIOS)), // Length: 16
					},
					[]byte(DHCPServerBootMenuDescriptionBIOS)..., // Description: Boot iPXE (BIOS)
				),
			),
		}

		// Serialize DHCP Option 43 suboptions
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

		// Add DHCP Option 43 suboptions and set the next server IP to 0.0.0.0
		outgoingDHCPPacket.Options = append(
			outgoingDHCPPacket.Options,
			layers.NewDHCPOption(
				layers.DHCPOptVendorOption,
				serializedSubOptions,
			),
		)
		outgoingDHCPPacket.NextServerIP = net.ParseIP("0.0.0.0").To4()
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

	log.Printf(`sending %v bytes of DHCP packets to client "%v"`, len(buf.Bytes()), braddr)

	// Broadcast the packet
	return conn.WriteToUDP(buf.Bytes(), braddr)
}
