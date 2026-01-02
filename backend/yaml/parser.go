package yaml

import (
	"fmt"

	yaml "gopkg.in/yaml.v3"
)

// Parser handles YAML parsing and generation
type Parser struct{}

// NewParser creates a new YAML parser
func NewParser() *Parser {
	return &Parser{}
}

// Parse parses YAML content into a UIPackage structure
func (p *Parser) Parse(content []byte) (*UIPackage, error) {
	if len(content) == 0 {
		return nil, fmt.Errorf("empty content")
	}

	var pkg UIPackage
	err := yaml.Unmarshal(content, &pkg)
	if err != nil {
		return nil, fmt.Errorf("yaml unmarshal failed: %w", err)
	}

	return &pkg, nil
}

// Generate generates YAML content from a UIPackage structure
func (p *Parser) Generate(pkg *UIPackage) ([]byte, error) {
	if pkg == nil {
		return nil, fmt.Errorf("package is nil")
	}

	content, err := yaml.Marshal(pkg)
	if err != nil {
		return nil, fmt.Errorf("yaml marshal failed: %w", err)
	}

	return content, nil
}

// Validate checks if the package structure is valid
func (p *Parser) Validate(pkg *UIPackage) error {
	if pkg == nil {
		return fmt.Errorf("package is nil")
	}

	// Basic validation - can be extended
	if pkg.Header.Version == 0 && len(pkg.Prototypes) == 0 && len(pkg.ImportedPackages) == 0 {
		return fmt.Errorf("package appears to be empty")
	}

	return nil
}

// ExtractAssets extracts all asset paths from a package
func (p *Parser) ExtractAssets(pkg *UIPackage) []string {
	assets := make(map[string]bool)

	// Add imported packages as assets
	for _, importedPkg := range pkg.ImportedPackages {
		assets[importedPkg] = true
	}

	// Recursively extract assets from controls
	for _, proto := range pkg.Prototypes {
		p.extractAssetsFromControl(proto, assets)
	}

	// Convert to slice
	result := make([]string, 0, len(assets))
	for asset := range assets {
		result = append(result, asset)
	}

	return result
}

// extractAssetsFromControl recursively extracts assets from a control
func (p *Parser) extractAssetsFromControl(ctrl *UIControl, assets map[string]bool) {
	if ctrl == nil {
		return
	}

	// Extract from components
	if ctrl.Components != nil {
		// Check Background component for sprite
		if bg, ok := ctrl.Components["Background"].(map[string]interface{}); ok {
			if sprite, ok := bg["sprite"].(string); ok && sprite != "" {
				assets[sprite] = true
			}
		}

		// Check StyleSheet component for styles
		if ss, ok := ctrl.Components["StyleSheet"].(map[string]interface{}); ok {
			if styles, ok := ss["styles"].(string); ok && styles != "" {
				assets[styles] = true
			}
		}

		// Check other components for asset references
		if uiAnim, ok := ctrl.Components["UIAnimationComponent"].(map[string]interface{}); ok {
			if animations, ok := uiAnim["animations"].(string); ok && animations != "" {
				// Parse animation references
				assets[animations] = true
			}
		}
	}

	// Extract from prototype reference
	if ctrl.Prototype != "" {
		assets[ctrl.Prototype] = true
	}

	// Recursively extract from children
	for _, child := range ctrl.Children {
		p.extractAssetsFromControl(child, assets)
	}
}

// FindControlByName searches for a control by name in the package
func (p *Parser) FindControlByName(pkg *UIPackage, name string) *UIControl {
	if pkg == nil {
		return nil
	}

	for _, proto := range pkg.Prototypes {
		if result := p.findControlByNameRecursive(proto, name); result != nil {
			return result
		}
	}

	return nil
}

// findControlByNameRecursive recursively searches for a control by name
func (p *Parser) findControlByNameRecursive(ctrl *UIControl, name string) *UIControl {
	if ctrl == nil {
		return nil
	}

	if ctrl.Name == name {
		return ctrl
	}

	for _, child := range ctrl.Children {
		if result := p.findControlByNameRecursive(child, name); result != nil {
			return result
		}
	}

	return nil
}
