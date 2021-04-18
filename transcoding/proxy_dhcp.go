package transcoding

import (
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func EncodeProxyDHCPPacket(
	clientHWAddr net.HardwareAddr,
	xid uint32,
	advertisedIP net.IP,
	clientIdentifierOpt layers.DHCPOption,
	yourClientIP net.IP,
	bootFileName string,
) []byte {
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

	// Create the outgoing packet
	outgoingDHCPPacket := &layers.DHCPv4{
		Operation:    layers.DHCPOpReply,
		HardwareType: layers.LinkTypeEthernet,
		HardwareLen:  uint8(len(clientHWAddr)),
		Xid:          xid,
		ClientIP:     net.ParseIP("0.0.0.0").To4(),
		YourClientIP: yourClientIP,
		NextServerIP: advertisedIP,
		RelayAgentIP: net.ParseIP("0.0.0.0").To4(),
		ClientHWAddr: clientHWAddr,
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

	return buf.Bytes()
}
