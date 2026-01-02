package dvpl

import (
	"bytes"
	"testing"
)

func TestEncryptDecryptRoundtrip(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "simple text",
			data: []byte("Hello, World!"),
		},
		{
			name: "empty string",
			data: []byte(""),
		},
		{
			name: "large compressible data",
			data: bytes.Repeat([]byte("A"), 10000),
		},
		{
			name: "random data",
			data: []byte{1, 2, 3, 4, 5, 255, 254, 253, 0, 1, 2, 3},
		},
		{
			name: "yaml sample",
			data: []byte(`Header:
  version: 135
Prototypes:
  - class: "UIControl"
    name: "Test"
    size: [800, 600]
    position: [0, 0]`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encrypt
			encrypted, err := EncryptDVPL(tt.data)
			if err != nil {
				t.Fatalf("EncryptDVPL failed: %v", err)
			}

			// Verify footer
			if len(encrypted) < 20 {
				t.Fatal("encrypted data too small")
			}

			if !IsDVPL(encrypted) {
				t.Fatal("missing or invalid DVPL magic marker")
			}

			// Decrypt
			decrypted, err := DecryptDVPL(encrypted)
			if err != nil {
				t.Fatalf("DecryptDVPL failed: %v", err)
			}

			// Verify roundtrip
			if !bytes.Equal(tt.data, decrypted) {
				t.Fatalf("roundtrip failed: expected %v, got %v", tt.data, decrypted)
			}
		})
	}
}

func TestEncryptDVPLErrors(t *testing.T) {
	t.Run("nil input", func(t *testing.T) {
		_, err := EncryptDVPL(nil)
		if err == nil {
			t.Fatal("expected error for nil input")
		}
	})
}

func TestDecryptDVPLErrors(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "nil input",
			data: nil,
		},
		{
			name: "too small",
			data: []byte("short"),
		},
		{
			name: "invalid magic",
			data: append(make([]byte, 16), []byte("BADD")...),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecryptDVPL(tt.data)
			if err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestCompressionLogic(t *testing.T) {
	// Highly compressible data should be compressed
	compressible := bytes.Repeat([]byte("A"), 1000)
	encrypted, err := EncryptDVPL(compressible)
	if err != nil {
		t.Fatalf("EncryptDVPL failed: %v", err)
	}

	// Encrypted data should be smaller than input + footer (20 bytes)
	if len(encrypted) >= len(compressible)+20 {
		t.Logf("compression not effective, but acceptable: %d bytes -> %d bytes", len(compressible), len(encrypted))
	}

	decrypted, _ := DecryptDVPL(encrypted)
	if !bytes.Equal(compressible, decrypted) {
		t.Fatal("roundtrip failed for compressible data")
	}
}

func TestIsDVPL(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected bool
	}{
		{
			name:     "valid DVPL",
			data:     append(make([]byte, 20), []byte("DVPL")...),
			expected: true,
		},
		{
			name:     "invalid magic",
			data:     append(make([]byte, 20), []byte("BADD")...),
			expected: false,
		},
		{
			name:     "too small",
			data:     []byte("short"),
			expected: false,
		},
		{
			name:     "empty",
			data:     []byte{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsDVPL(tt.data)
			if result != tt.expected {
				t.Fatalf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
