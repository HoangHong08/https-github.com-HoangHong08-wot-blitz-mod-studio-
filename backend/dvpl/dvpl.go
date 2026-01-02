package dvpl

import (
	"encoding/binary"
	"errors"
	"hash/crc32"

	"github.com/pierrec/lz4/v4"
)

// EncryptDVPL compresses and encrypts data to DVPL format
// Returns data with DVPL footer: [data][20-byte footer]
// Footer: [originalSize(4)][compressedSize(4)][crc32(4)][type(4)][DVPL(4)]
func EncryptDVPL(inputBuf []byte) ([]byte, error) {
	if inputBuf == nil {
		return nil, errors.New("input buffer is nil")
	}

	inputSize := len(inputBuf)

	compressedBuf := make([]byte, lz4.CompressBlockBound(inputSize))
	actualCompressedSize, err := lz4.CompressBlock(inputBuf, compressedBuf, nil)
	if err != nil {
		return nil, errors.New("failed to compress lz4: " + err.Error())
	}

	compressedBuf = compressedBuf[:actualCompressedSize]

	var outputBuf []byte
	var compressionType uint32

	// Skip compression when compressed size >= input size
	if inputSize <= actualCompressedSize {
		// Store uncompressed (compression type = 0)
		outputBuf = binary.LittleEndian.AppendUint32(inputBuf, uint32(inputSize))
		outputBuf = binary.LittleEndian.AppendUint32(outputBuf, uint32(inputSize))

		crc32DataSum := crc32.ChecksumIEEE(inputBuf)
		outputBuf = binary.LittleEndian.AppendUint32(outputBuf, crc32DataSum)

		compressionType = 0
	} else {
		// Store compressed (compression type = 1)
		outputBuf = binary.LittleEndian.AppendUint32(compressedBuf, uint32(inputSize))
		outputBuf = binary.LittleEndian.AppendUint32(outputBuf, uint32(actualCompressedSize))

		crc32DataSum := crc32.ChecksumIEEE(compressedBuf)
		outputBuf = binary.LittleEndian.AppendUint32(outputBuf, crc32DataSum)

		compressionType = 1
	}

	outputBuf = binary.LittleEndian.AppendUint32(outputBuf, compressionType)
	outputBuf = append(outputBuf, []byte("DVPL")...)

	return outputBuf, nil
}

// DecryptDVPL decrypts and decompresses DVPL format data
// Expects data with DVPL footer: [data][20-byte footer]
func DecryptDVPL(inputBuf []byte) ([]byte, error) {
	if inputBuf == nil {
		return nil, errors.New("input buffer is nil")
	}

	if len(inputBuf) < 20 {
		return nil, errors.New("input buffer too small (< 20 bytes)")
	}

	// Split data and footer
	dataBuf := inputBuf[:len(inputBuf)-20]
	footerBuf := inputBuf[len(inputBuf)-20:]

	// Parse footer
	originalSize := binary.LittleEndian.Uint32(footerBuf[:4])
	compressedSize := binary.LittleEndian.Uint32(footerBuf[4:8])

	if int(compressedSize) != len(dataBuf) {
		return nil, errors.New("invalid compressed data length")
	}

	crc32DataSum := binary.LittleEndian.Uint32(footerBuf[8:12])
	if crc32DataSum != crc32.ChecksumIEEE(dataBuf) {
		return nil, errors.New("invalid crc32 sum")
	}

	compressType := binary.LittleEndian.Uint32(footerBuf[12:16])

	// Verify magic marker
	if string(footerBuf[16:20]) != "DVPL" {
		return nil, errors.New("invalid DVPL magic marker")
	}

	// Decompress if needed
	outputBuf := make([]byte, originalSize)
	if compressType == 0 {
		// No compression - data is stored as-is
		copy(outputBuf, dataBuf)
	} else if compressType == 1 {
		// LZ4 compression
		actualOutputSize, err := lz4.UncompressBlock(dataBuf, outputBuf)
		if err != nil {
			return nil, errors.New("failed to decompress lz4: " + err.Error())
		}
		outputBuf = outputBuf[:actualOutputSize]
	} else {
		return nil, errors.New("unknown compression type")
	}

	return outputBuf, nil
}

// IsDVPL checks if the data has a valid DVPL footer
func IsDVPL(data []byte) bool {
	if len(data) < 20 {
		return false
	}
	return string(data[len(data)-4:]) == "DVPL"
}
