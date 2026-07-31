package units

import (
	"fmt"
	"strings"

	"github.com/c2h5oh/datasize"
)

const Megabyte = int(datasize.MB)

// ParseByteSize parses a human-readable byte size string into bytes.
// It accepts units supported by datasize, such as "256m", "256MB", and "1g".
func ParseByteSize(size string) (int, error) {
	if strings.TrimSpace(size) == "" {
		return 0, fmt.Errorf("byte size is required")
	}

	var parsed datasize.ByteSize
	if err := parsed.UnmarshalText([]byte(size)); err != nil {
		return 0, fmt.Errorf("invalid byte size %q: %w", size, err)
	}
	if parsed > datasize.ByteSize(^uint(0)>>1) {
		return 0, fmt.Errorf("byte size %q is too large", size)
	}

	return int(parsed), nil
}
