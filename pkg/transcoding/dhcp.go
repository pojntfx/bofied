package transcoding

import (
	"errors"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

const (
	DHCPBaseClassID                 = "PXEClient"
	DHCPEFIArch                     = 7  // X86-64_EFI
	DHCPOptUUIDGUIDClientIdentifier = 97 // Option: (97) UUID/GUID-based Client Identifier
)

type DecodedDHCPPacket struct {
	IsPXE               bool
	ClientHWAddr        net.HardwareAddr
	Xid                 uint32
	ClientIdentifierOpt layers.DHCPOption
	Arch                int
	Undi                int
}

func DecodeDHCPPacket(rawIncomingUDPPacket []byte) (DecodedDHCPPacket, error) {
	// Decode and parse packet
	incomingUDPPacket := gopacket.NewPacket(rawIncomingUDPPacket, layers.LayerTypeDHCPv4, gopacket.Default)

	dhcpLayer := incomingUDPPacket.Layer(layers.LayerTypeDHCPv4)
	if dhcpLayer == nil {
		return DecodedDHCPPacket{}, errors.New("could not parse DHCP layer: not a DHCP layer")
	}

	incomingDHCPPacket, ok := dhcpLayer.(*layers.DHCPv4)
	if !ok {
		return DecodedDHCPPacket{}, errors.New("could not parse DHCP layer: invalid DHCP layer")
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
				pxe, a, u, err := ParsePXEClassIdentifier(string(option.Data))
				if err != nil {
					return DecodedDHCPPacket{}, err
				}

				isPXE = pxe
				arch = a
				undi = u
			case DHCPOptUUIDGUIDClientIdentifier:
				clientIdentifierOpt = option
			}
		}
	}

	return DecodedDHCPPacket{
		IsPXE:               isPXE,
		ClientHWAddr:        incomingDHCPPacket.ClientHWAddr,
		Xid:                 incomingDHCPPacket.Xid,
		ClientIdentifierOpt: clientIdentifierOpt,
		Arch:                arch,
		Undi:                undi,
	}, nil
}

func EncodeDHCPPacket(
	clientHWAddr net.HardwareAddr,
	xid uint32,
	advertisedIP net.IP,
	clientIdentifierOpt layers.DHCPOption,
	arch int,
	bootMenuPromptBIOS string,
	bootMenuDescriptionBIOS string,
) []byte {
	// Create the outgoing packet
	outgoingDHCPPacket := &layers.DHCPv4{
		Operation:    layers.DHCPOpReply,
		HardwareType: layers.LinkTypeEthernet,
		HardwareLen:  uint8(len(clientHWAddr)),
		Xid:          xid,
		ClientIP:     net.ParseIP("0.0.0.0").To4(),
		YourClientIP: net.ParseIP("0.0.0.0").To4(),
		NextServerIP: advertisedIP,
		RelayAgentIP: net.ParseIP("0.0.0.0").To4(),
		ClientHWAddr: clientHWAddr,
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
					[]byte(bootMenuPromptBIOS)..., // Prompt: PXE
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
						byte(len(bootMenuDescriptionBIOS)), // Length: 16
					},
					[]byte(bootMenuDescriptionBIOS)..., // Description: Boot iPXE (BIOS)
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

	return buf.Bytes()
}
