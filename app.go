package main

import (
	"context"
	"fmt"
	"os"

	"wot-blitz-mod-studio/backend/dvpl"
	"wot-blitz-mod-studio/backend/yaml"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx          context.Context
	parser       *yaml.Parser
	gameDataPath string
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		parser: yaml.NewParser(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	// Try to detect game data path from common locations
	a.detectGameDataPath()
}

// detectGameDataPath tries to find the game data directory
func (a *App) detectGameDataPath() {
	// Common paths to check
	paths := []string{
		"/opt/wot-blitz",
		"/opt/games/wot-blitz",
		"/home/wot-blitz",
		"/usr/local/games/wot-blitz",
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			a.gameDataPath = path
			return
		}
	}
}

// SetGameDataPath sets the path to game data directory
func (a *App) SetGameDataPath(path string) error {
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}
	a.gameDataPath = path
	return nil
}

// OpenFile opens and parses a UI file (DVPL or YAML)
func (a *App) OpenFile(filePath string) (*yaml.FileData, error) {
	if filePath == "" {
		return nil, fmt.Errorf("file path is empty")
	}

	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Detect format and decrypt if needed
	var yamlContent []byte
	if dvpl.IsDVPL(data) {
		// Decrypt DVPL
		decrypted, err := dvpl.DecryptDVPL(data)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt DVPL: %w", err)
		}
		yamlContent = decrypted
	} else {
		// Assume YAML
		yamlContent = data
	}

	// Parse YAML
	pkg, err := a.parser.Parse(yamlContent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Extract assets
	assets := a.parser.ExtractAssets(pkg)

	return &yaml.FileData{
		Path:    filePath,
		Content: string(yamlContent),
		Package: pkg,
		Assets:  assets,
	}, nil
}

// SaveFile saves the modified content back to file
func (a *App) SaveFile(filePath string, content string, wasDVPL bool) error {
	if filePath == "" {
		return fmt.Errorf("file path is empty")
	}

	contentBytes := []byte(content)

	// If original was DVPL, encrypt before saving
	if wasDVPL {
		encrypted, err := dvpl.EncryptDVPL(contentBytes)
		if err != nil {
			return fmt.Errorf("failed to encrypt DVPL: %w", err)
		}
		contentBytes = encrypted
	}

	// Write file
	if err := os.WriteFile(filePath, contentBytes, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// ParseYAML parses YAML content and returns the structured data
func (a *App) ParseYAML(content string) (*yaml.UIPackage, error) {
	if content == "" {
		return nil, fmt.Errorf("content is empty")
	}

	pkg, err := a.parser.Parse([]byte(content))
	if err != nil {
		return nil, fmt.Errorf("parse failed: %w", err)
	}

	return pkg, nil
}

// GenerateYAML generates YAML content from parsed data
func (a *App) GenerateYAML(pkg *yaml.UIPackage) (string, error) {
	if pkg == nil {
		return "", fmt.Errorf("package is nil")
	}

	content, err := a.parser.Generate(pkg)
	if err != nil {
		return "", fmt.Errorf("generate failed: %w", err)
	}

	return string(content), nil
}

// FindControl finds a control by name in the current package
func (a *App) FindControl(pkg *yaml.UIPackage, name string) *yaml.UIControl {
	return a.parser.FindControlByName(pkg, name)
}

// SelectFile opens a file picker dialog
func (a *App) SelectFile() (string, error) {
	options := runtime.OpenDialogOptions{
		DefaultDirectory: a.gameDataPath,
		Title:            "Select YAML or DVPL file",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "UI Files (*.yaml, *.sc2, *.dvpl)",
				Pattern:     "*.yaml;*.sc2;*.dvpl",
			},
			{
				DisplayName: "YAML files (*.yaml)",
				Pattern:     "*.yaml",
			},
			{
				DisplayName: "DVPL files (*.sc2, *.dvpl)",
				Pattern:     "*.sc2;*.dvpl",
			},
			{
				DisplayName: "All files (*.*)",
				Pattern:     "*.*",
			},
		},
	}

	filePath, err := runtime.OpenFileDialog(a.ctx, options)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

// SelectFolder opens a folder picker dialog
func (a *App) SelectFolder() (string, error) {
	options := runtime.OpenDialogOptions{
		DefaultDirectory: a.gameDataPath,
		Title:            "Select WoT Blitz Game Data Folder",
	}

	folderPath, err := runtime.OpenFileDialog(a.ctx, options)
	if err != nil {
		return "", err
	}

	a.gameDataPath = folderPath
	return folderPath, nil
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
