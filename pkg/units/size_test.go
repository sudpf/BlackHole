package units

import "testing"

func TestParseByteSize(t *testing.T) {
	tests := map[string]int{
		"10240": 10240,
		"256m":  268435456,
		"256MB": 268435456,
		"1g":    1073741824,
		"1GB":   1073741824,
	}

	for input, want := range tests {
		got, err := ParseByteSize(input)
		if err != nil {
			t.Fatalf("ParseByteSize(%q) error: %v", input, err)
		}
		if got != want {
			t.Fatalf("ParseByteSize(%q) = %d, want %d", input, got, want)
		}
	}
}

func TestParseByteSizeRejectsInvalidSize(t *testing.T) {
	for _, input := range []string{"", "bad"} {
		if _, err := ParseByteSize(input); err == nil {
			t.Fatalf("ParseByteSize(%q) expected error", input)
		}
	}
}
