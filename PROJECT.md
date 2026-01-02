# WoT Blitz Mod Studio

A professional GUI application for modding World of Tanks Blitz UI files.

## Features

- **DVPL Encryption/Decryption**: Full support for encrypted UI file format
- **YAML Editing**: Real-time YAML editing with syntax validation
- **Visual Preview**: Canvas-based preview of UI controls
- **Tree Navigation**: Hierarchical control tree browser
- **In-Memory Processing**: All operations happen in memory with zero temp files
- **Cross-Platform**: Built with Wails v2 for Windows, macOS, and Linux

## Project Structure

```
wot-blitz-mod-studio/
├── main.go                    # Wails entry point
├── app.go                     # Go ↔ Svelte bridge
├── go.mod / go.sum
├── wails.json                 # Wails configuration
│
├── backend/
│   ├── dvpl/
│   │   ├── dvpl.go           # DVPL encryption/decryption
│   │   └── dvpl_test.go      # DVPL tests
│   │
│   └── yaml/
│       ├── types.go          # YAML data structures
│       ├── parser.go         # YAML parsing logic
│       └── parser_test.go    # YAML tests
│
└── frontend/
    ├── package.json
    ├── vite.config.js
    ├── index.html
    └── src/
        ├── main.js
        ├── App.svelte
        ├── style.css
        └── components/
            ├── Toolbar.svelte
            ├── Sidebar.svelte
            ├── Editor.svelte
            └── Preview.svelte
```

## Technology Stack

- **Backend**: Go 1.23 with Wails v2
- **Frontend**: Svelte 4 with Vite
- **Compression**: github.com/pierrec/lz4/v4 (pure Go, ARM64 compatible)
- **YAML**: gopkg.in/yaml.v3 with custom unmarshalers

## Building

### Development Mode

```bash
wails dev
```

This will start the development server with hot reload.

### Production Build (Linux ARM64)

```bash
CGO_ENABLED=0 wails build -platform linux/arm64 -clean
```

### Creating Debian Package

```bash
dpkg-deb --build debian-package wot-blitz-mod-studio_1.0.0_arm64.deb
```

## API Reference

### Backend Methods (Go → Frontend)

- `OpenFile(path string) (*FileData, error)` - Open and parse a UI file
- `SaveFile(path string, content string, wasDVPL bool) error` - Save modified content
- `ParseYAML(content string) (*Package, error)` - Parse YAML content
- `GenerateYAML(pkg *Package) (string, error)` - Generate YAML from package
- `FindControl(pkg *Package, name string) *UIControl` - Search for control by name
- `SetGameDataPath(path string) error` - Set game data directory

## File Format Support

### DVPL Format (.sc2)

- **Header**: 20 bytes footer with metadata
- **Compression**: LZ4 block mode
- **CRC32**: IEEE polynomial validation
- **Structure**: [data][originalSize(4)][compressedSize(4)][crc32(4)][type(4)][DVPL(4)]

### YAML Format

Standard Dava Engine YAML with support for:
- **Vector2**: `[x, y]` arrays
- **Color**: `[r, g, b, a]` arrays (0.0-1.0)
- **Asset Paths**: `~res:/` prefix for resource resolution
- **Nested Controls**: Full hierarchy with unlimited depth

## Testing

### Run All Tests

```bash
go test ./...
```

### DVPL Tests

```bash
go test ./backend/dvpl -v
```

### YAML Parser Tests

```bash
go test ./backend/yaml -v
```

## Performance Targets

- **Binary Size**: < 50MB
- **Memory Usage**: < 200MB during editing
- **Startup Time**: < 2 seconds
- **Zero Temp Files**: All processing in-memory

## Known Limitations

- Preview canvas is 2D representation only
- WebP atlas support (future feature)
- Asset preview (future feature)
- Real-time Monaco editor integration (simplified textarea in this version)

## Contributing

Contributions are welcome. Please ensure:
- All tests pass
- Code follows Go conventions
- Backend changes include unit tests
- No temporary files are created during operations

## License

See LICENSE file

## Progress Status

### Completed ✅
- Source code analysis from reference implementations
- DVPL encryption/decryption backend with full test coverage
- YAML parser with custom unmarshalers for Vector2, Color, UIControl
- Wails project initialization
- 3-column UI layout with Toolbar, Sidebar, Editor, Preview components
- Backend bridge (app.go) with all core methods
- Unit tests for DVPL and YAML parsing

### In Progress 🔄
- Frontend component integration and styling refinement
- File picker dialog implementation

### Not Started ⏳
- Manual integration testing with real RXD-MODPACK files
- ARM64 Debian package creation
- Advanced features (WebP atlas, asset preview, Monaco editor integration)

---

**Last Updated**: January 2, 2026
**Version**: 0.1.0-alpha
