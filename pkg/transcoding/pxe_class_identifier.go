package transcoding

import (
	"strconv"
	"strings"
)

func ParsePXEClassIdentifier(classID string) (isPXE bool, arch int, err error) {
	parts := strings.Split(classID, ":")

	for i, part := range parts {
		switch part {
		case "PXEClient":
			isPXE = true
		case "Arch":
			if len(parts) > i {
				arch, err = strconv.Atoi(parts[i+1])
				if err != nil {
					return
				}
			}
		}
	}

	return
}
