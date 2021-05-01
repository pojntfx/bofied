package config

import "strconv"

// See https://www.iana.org/assignments/dhcpv6-parameters/dhcpv6-parameters.xhtml#processor-architecture
func GetNameForArchId(id int) string {
	switch id {
	case 0x00:
		return "x86 BIOS"
	case 0x01:
		return "NEC/PC98 (DEPRECATED)"
	case 0x02:
		return "Itanium"
	case 0x03:
		return "DEC Alpha (DEPRECATED)"
	case 0x04:
		return "Arc x86 (DEPRECATED)"
	case 0x05:
		return "Intel Lean Client (DEPRECATED)"
	case 0x06:
		return "x86 UEFI"
	case 0x07:
		return "x64 UEFI"
	case 0x08:
		return "EFI Xscale (DEPRECATED)"
	case 0x09:
		return "EBC"
	case 0x0a:
		return "ARM 32-bit UEFI"
	case 0x0b:
		return "ARM 64-bit UEFI"
	case 0x0c:
		return "PowerPC Open Firmware"
	case 0x0d:
		return "PowerPC ePAPR"
	case 0x0e:
		return "POWER OPAL v3"
	case 0x0f:
		return "x86 uefi boot from http"
	case 0x10:
		return "x64 uefi boot from http"
	case 0x11:
		return "ebc boot from http"
	case 0x12:
		return "arm uefi 32 boot from http"
	case 0x13:
		return "arm uefi 64 boot from http"
	case 0x14:
		return "pc/at bios boot from http"
	case 0x15:
		return "arm 32 uboot"
	case 0x16:
		return "arm 64 uboot"
	case 0x17:
		return "arm uboot 32 boot from http"
	case 0x18:
		return "arm uboot 64 boot from http"
	case 0x19:
		return "RISC-V 32-bit UEFI"
	case 0x1a:
		return "RISC-V 32-bit UEFI boot from http"
	case 0x1b:
		return "RISC-V 64-bit UEFI"
	case 0x1c:
		return "RISC-V 64-bit UEFI boot from http"
	case 0x1d:
		return "RISC-V 128-bit UEFI"
	case 0x1e:
		return "RISC-V 128-bit UEFI boot from http"
	case 0x1f:
		return "s390 Basic"
	case 0x20:
		return "s390 Extended"
	case 0x21:
		return "MIPS 32-bit UEFI"
	case 0x22:
		return "MIPS 64-bit UEFI"
	case 0x23:
		return "Sunway 32-bit UEFI"
	case 0x24:
		return "Sunway 64-bit UEFI"
	default:
		return strconv.Itoa(id)
	}
}
