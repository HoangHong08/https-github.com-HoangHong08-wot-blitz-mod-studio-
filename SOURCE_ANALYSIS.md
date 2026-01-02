# WoT Blitz Mod Studio - Source Analysis

## DVPL Format Analysis

### EncryptDVPL Function
```go
func EncryptDVPL(inputBuf []byte) ([]byte, error) {
	inputSize := len(inputBuf)

	compressedBuf := make([]byte, lz4.CompressBlockBound(inputSize))
	actualCompressedSize, err := lz4.CompressBlock(inputBuf, compressedBuf, nil)
	if err != nil {
		return nil, errors.New("failed to compress lz4")
	}

	compressedBuf = compressedBuf[:actualCompressedSize]

	var outputBuf []byte

	// Skip compression when compressed size > input size
	if inputSize < actualCompressedSize {
		outputBuf = binary.LittleEndian.AppendUint32(inputBuf, uint32(inputSize))
		outputBuf = binary.LittleEndian.AppendUint32(outputBuf, uint32(inputSize))

		crc32DataSum := crc32.ChecksumIEEE(inputBuf)
		outputBuf = binary.LittleEndian.AppendUint32(outputBuf, crc32DataSum)

		outputBuf = binary.LittleEndian.AppendUint32(outputBuf, 0)
	} else {
		outputBuf = binary.LittleEndian.AppendUint32(compressedBuf, uint32(inputSize))
		outputBuf = binary.LittleEndian.AppendUint32(outputBuf, uint32(actualCompressedSize))

		crc32DataSum := crc32.ChecksumIEEE(compressedBuf)
		outputBuf = binary.LittleEndian.AppendUint32(outputBuf, crc32DataSum)

		outputBuf = binary.LittleEndian.AppendUint32(outputBuf, 1)
	}
	outputBuf = append(outputBuf, []byte("DVPL")...)
	return outputBuf, nil
}
```

### DecryptDVPL Function
```go
func DecryptDVPL(inputBuf []byte) ([]byte, error) {
	dataBuf := inputBuf[:len(inputBuf)-20]
	footerBuf := inputBuf[len(inputBuf)-20:]

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
	outputBuf := make([]byte, originalSize)
	if compressType == 0 {
		outputBuf = dataBuf
	} else {
		actualOutputSize, err := lz4.UncompressBlock(dataBuf, outputBuf)
		if err != nil {
			return nil, errors.New("failed to uncompressed lz4")
		}
		outputBuf = outputBuf[:actualOutputSize]
	}
	return outputBuf, nil
}
```

## DVPL Format Specification

### File Structure
```
[Encrypted/Compressed Data] [20-byte Footer]
```

### Footer Structure (20 bytes)
- **Bytes 0-3**: Original size (uint32, little-endian)
- **Bytes 4-7**: Compressed size (uint32, little-endian)
- **Bytes 8-11**: CRC32 checksum (uint32, little-endian) - IEEE polynomial
- **Bytes 12-15**: Compression type (uint32, little-endian)
  - `0`: No compression (data stored as-is)
  - `1`: LZ4 block compression
- **Bytes 16-19**: Magic marker `"DVPL"`

### Compression Logic
- Uses LZ4 block-mode compression (`github.com/pierrec/lz4/v4`)
- Only applies compression if it reduces file size
- If compressed size ≥ original size, stores uncompressed with type=0
- CRC32 is calculated on the data BEFORE appending footer

## YAML Format Patterns Found

### Standard YAML Structure
```yaml
Header:
    version: 135

ImportedPackages:
  - "~res:/UI/Screens/Battle/Shell.yaml"

Prototypes:
  - class: "UIControl"
    name: "ControlName"
    position: [x, y]
    size: [width, height]
    visible: true
    input: false
    pivot: [0.5, 0.5]
    classes: "classname another-class"
    components:
        Background:
            drawType: "DRAW_FILL"
            sprite: "~res:/Gfx/UI/path/image.webp"
            color: [1.0, 1.0, 1.0, 1.0]
        Anchor:
            leftAnchorEnabled: true
            rightAnchorEnabled: true
            topAnchorEnabled: true
            bottomAnchorEnabled: true
    children:
      - class: "UIControl"
        name: "ChildControl"
        # ... nested structure
```

### Data Type Patterns

#### Vector2 / Size (Arrays with 2 elements)
```yaml
position: [100.0, 200.0]     # [x, y]
size: [800.0, 600.0]         # [width, height]
pivot: [0.5, 0.5]            # [pivotX, pivotY]
```
- Represented as `[x, y]` format
- Floating-point values
- Used for positioning and dimensions

#### Color (Arrays with 4 elements)
```yaml
color: [1.0, 1.0, 1.0, 1.0]  # [R, G, B, A]
```
- Represented as `[r, g, b, a]` format
- Values in range [0.0, 1.0]
- A component is alpha/opacity

#### Asset Paths
```yaml
sprite: "~res:/Gfx/UI/MainScreen/bg.webp"
prototype: "UGN/UgnUvn"
styles: "~res:/UI/Screens/Battle/Aims/Aims.style.yaml"
```
- Resource paths use `~res:/` prefix
- Relative to game data directory
- Can reference YAML, WebP, and other assets

#### Component Types Found
- `Background`: Drawing properties (sprite, color, drawType)
- `Anchor`: Anchoring configuration (leftAnchorEnabled, etc.)
- `SizePolicy`: Size policy (horizontalPolicy, verticalPolicy)
- `UIAnimationComponent`: Animations
- `UIOpacityComponent`: Opacity settings
- `StyleSheet`: Style references

### Key Observations

1. **Nested Hierarchy**: UIControl can have `children` array with unlimited nesting
2. **Components Map**: Each control has a `components` map with various component types
3. **Prototype References**: Can reference prototypes: `prototype: "UGN/UgnUvn"`
4. **Custom Classes**: Support for CSS-like classes: `classes: "classname another-class"`
5. **Boolean Flags**: Common properties like `visible`, `input`, `enabled`
6. **String Keys**: Property names are typically camelCase

## Implementation Requirements

### Required Go Libraries
- `gopkg.in/yaml.v3` - YAML parsing with custom unmarshalers
- `github.com/pierrec/lz4/v4` - LZ4 compression (pure Go, ARM64 compatible)
- `encoding/binary` - Binary data handling
- `hash/crc32` - CRC32 validation
- `golang.org/x/image` - Image decoding (future WebP support)

### Custom YAML Unmarshalers Needed
1. **Vector2**: Unmarshal `[x, y]` array format
2. **Color**: Unmarshal `[r, g, b, a]` array format
3. **UIControl**: Handle nested children hierarchy
4. **Components**: Map-based component structure

### Asset Path Resolution Strategy
```
Input: ~res:/Gfx/UI/MainScreen/bg.webp

Search paths (in order):
1. {gameDataPath}/3d/Gfx/UI/MainScreen/bg.webp
2. {gameDataPath}/Gfx/UI/MainScreen/bg.webp
3. {gameDataPath}/UI/MainScreen/bg.webp
4. {gameDataPath}/resources/Gfx/UI/MainScreen/bg.webp

Caching: Store resolved paths to avoid repeated lookups
```

## Open Questions & Clarifications

### Q: Handle .sc2 vs .yaml file extensions?
**A**: .sc2 files are encrypted DVPL format, .yaml files are plaintext YAML

### Q: DrawType enum values?
**A**: From analysis: `DRAW_FILL`, `STRETCH_BOTH`, `TILED`, `PER_PIXEL_ACCURACY_ENABLED`, etc.

### Q: Size limits for in-memory editing?
**A**: Target < 200MB memory usage, suitable for most UI files (typically < 10MB)

### Q: WebP handling with atlas metadata?
**A**: Future feature - see `packed_webp-1.1` source for format

---

## Next Steps

1. ✅ Analyze sources
2. Create Wails project with Svelte frontend
3. Implement DVPL encryption/decryption backend
4. Implement YAML parser with custom unmarshalers
5. Build frontend 3-column editor layout
6. Integrate with Monaco for YAML editing
7. Add real-time preview capability
8. Package as ARM64 .deb file
