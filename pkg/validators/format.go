package validators

import "go/format"

func FormatGoSrc(src string) (string, error) {
	s, err := format.Source([]byte(src))
	if err != nil {
		return src, err
	}

	return string(s), err
}
