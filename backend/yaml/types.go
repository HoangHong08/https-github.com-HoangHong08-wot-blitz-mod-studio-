package yaml

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// Vector2 represents a 2D vector or size
type Vector2 struct {
	X float64
	Y float64
}

// UnmarshalYAML implements custom YAML unmarshaling for [x, y] format
func (v *Vector2) UnmarshalYAML(value *yaml.Node) error {
	var arr [2]float64
	if err := value.Decode(&arr); err != nil {
		return fmt.Errorf("failed to decode Vector2: %w", err)
	}
	v.X = arr[0]
	v.Y = arr[1]
	return nil
}

// MarshalYAML implements custom YAML marshaling for [x, y] format
func (v Vector2) MarshalYAML() (interface{}, error) {
	return [2]float64{v.X, v.Y}, nil
}

// Color represents RGBA color
type Color struct {
	R float64
	G float64
	B float64
	A float64
}

// UnmarshalYAML implements custom YAML unmarshaling for [r, g, b, a] format
func (c *Color) UnmarshalYAML(value *yaml.Node) error {
	var arr [4]float64
	if err := value.Decode(&arr); err != nil {
		return fmt.Errorf("failed to decode Color: %w", err)
	}
	c.R = arr[0]
	c.G = arr[1]
	c.B = arr[2]
	c.A = arr[3]
	return nil
}

// MarshalYAML implements custom YAML marshaling for [r, g, b, a] format
func (c Color) MarshalYAML() (interface{}, error) {
	return [4]float64{c.R, c.G, c.B, c.A}, nil
}

// Component represents a generic component with dynamic fields
type Component map[string]interface{}

// UIControl represents a UI control element in Dava UI system
type UIControl struct {
	Class       string                 `yaml:"class,omitempty"`
	CustomClass string                 `yaml:"customClass,omitempty"`
	Name        string                 `yaml:"name,omitempty"`
	Position    *Vector2               `yaml:"position,omitempty"`
	Size        *Vector2               `yaml:"size,omitempty"`
	Pivot       *Vector2               `yaml:"pivot,omitempty"`
	Visible     *bool                  `yaml:"visible,omitempty"`
	Input       *bool                  `yaml:"input,omitempty"`
	Classes     string                 `yaml:"classes,omitempty"`
	Prototype   string                 `yaml:"prototype,omitempty"`
	Components  map[string]interface{} `yaml:"components,omitempty"`
	Children    []*UIControl           `yaml:"children,omitempty"`

	// Dynamic properties not explicitly typed
	Properties map[string]interface{} `yaml:",inline"`
}

// UIPackage represents a complete YAML UI package
type UIPackage struct {
	Header           Header            `yaml:"Header,omitempty"`
	ImportedPackages []string          `yaml:"ImportedPackages,omitempty"`
	Prototypes       []*UIControl      `yaml:"Prototypes,omitempty"`
	ExternalPackages map[string]string `yaml:"ExternalPackages,omitempty"`
}

// Header represents YAML header information
type Header struct {
	Version int `yaml:"version,omitempty"`
}

// FileData represents the complete parsed UI file
type FileData struct {
	Path    string     `json:"path"`
	Content string     `json:"content"`
	Package *UIPackage `json:"package"`
	Assets  []string   `json:"assets"`
}
