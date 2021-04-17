package config

func GetFileName(
	ip string,
	macAddress string,
	arch int,
	undi int,
) string {
	switch arch {
	case 7:
		return "ipxe.efi"
	default:
		return "ipxe.kpxe"
	}
}
