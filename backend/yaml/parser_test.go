package yaml

import (
	"testing"
)

func TestColorUnmarshal(t *testing.T) {
	yamlContent := `Header:
  version: 1
Prototypes:
  - class: "UIControl"
    name: "TestControl"
    components:
      Background:
        color: [1.0, 0.5, 0.25, 0.8]`

	parser := NewParser()
	pkg, err := parser.Parse([]byte(yamlContent))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if pkg == nil || len(pkg.Prototypes) == 0 {
		t.Fatal("expected prototypes in parsed package")
	}

	ctrl := pkg.Prototypes[0]
	if ctrl.Name != "TestControl" {
		t.Fatalf("expected control name TestControl, got %s", ctrl.Name)
	}
}

func TestParseSimpleYAML(t *testing.T) {
	yamlContent := `Header:
  version: 135
ImportedPackages:
  - "~res:/UI/Screens/Battle/Shell.yaml"
Prototypes:
  - class: "UIControl"
    name: "MainScreen"
    size: [1920, 1080]
    position: [0, 0]
    visible: true`

	parser := NewParser()
	pkg, err := parser.Parse([]byte(yamlContent))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if pkg.Header.Version != 135 {
		t.Fatalf("expected version 135, got %d", pkg.Header.Version)
	}

	if len(pkg.ImportedPackages) != 1 {
		t.Fatalf("expected 1 imported package, got %d", len(pkg.ImportedPackages))
	}

	if len(pkg.Prototypes) != 1 {
		t.Fatalf("expected 1 prototype, got %d", len(pkg.Prototypes))
	}

	proto := pkg.Prototypes[0]
	if proto.Name != "MainScreen" {
		t.Fatalf("expected name MainScreen, got %s", proto.Name)
	}

	if proto.Class != "UIControl" {
		t.Fatalf("expected class UIControl, got %s", proto.Class)
	}

	if proto.Size == nil || proto.Size.X != 1920 || proto.Size.Y != 1080 {
		t.Fatalf("expected size [1920, 1080], got [%f, %f]", proto.Size.X, proto.Size.Y)
	}

	if *proto.Visible != true {
		t.Fatal("expected visible to be true")
	}
}

func TestParseNestedControls(t *testing.T) {
	yamlContent := `Header:
  version: 1
Prototypes:
  - class: "UIControl"
    name: "Parent"
    size: [800, 600]
    children:
      - class: "UIControl"
        name: "Child1"
        size: [400, 300]
        position: [0, 0]
      - class: "UIControl"
        name: "Child2"
        size: [400, 300]
        position: [400, 0]
        children:
          - class: "UIControl"
            name: "Grandchild"
            size: [200, 150]`

	parser := NewParser()
	pkg, err := parser.Parse([]byte(yamlContent))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	parent := pkg.Prototypes[0]
	if parent.Name != "Parent" {
		t.Fatalf("expected parent name, got %s", parent.Name)
	}

	if len(parent.Children) != 2 {
		t.Fatalf("expected 2 children, got %d", len(parent.Children))
	}

	if parent.Children[0].Name != "Child1" {
		t.Fatalf("expected Child1, got %s", parent.Children[0].Name)
	}

	if len(parent.Children[1].Children) != 1 {
		t.Fatalf("expected 1 grandchild, got %d", len(parent.Children[1].Children))
	}

	grandchild := parent.Children[1].Children[0]
	if grandchild.Name != "Grandchild" {
		t.Fatalf("expected Grandchild, got %s", grandchild.Name)
	}
}

func TestGenerateYAML(t *testing.T) {
	originalYAML := `Header:
  version: 135
Prototypes:
  - class: "UIControl"
    name: "Test"
    size: [800, 600]`

	parser := NewParser()

	// Parse
	pkg, err := parser.Parse([]byte(originalYAML))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Generate
	generated, err := parser.Generate(pkg)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Re-parse generated YAML
	pkg2, err := parser.Parse(generated)
	if err != nil {
		t.Fatalf("Re-parse failed: %v", err)
	}

	if pkg2.Header.Version != pkg.Header.Version {
		t.Fatal("version mismatch after roundtrip")
	}

	if len(pkg2.Prototypes) != len(pkg.Prototypes) {
		t.Fatal("prototype count mismatch after roundtrip")
	}
}

func TestFindControlByName(t *testing.T) {
	yamlContent := `Header:
  version: 1
Prototypes:
  - class: "UIControl"
    name: "Parent"
    children:
      - class: "UIControl"
        name: "Child1"
      - class: "UIControl"
        name: "Child2"
        children:
          - class: "UIControl"
            name: "Grandchild"`

	parser := NewParser()
	pkg, err := parser.Parse([]byte(yamlContent))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	tests := []struct {
		name        string
		searchName  string
		expectFound bool
	}{
		{"Parent", "Parent", true},
		{"Child1", "Child1", true},
		{"Grandchild", "Grandchild", true},
		{"NotFound", "NotFound", false},
	}

	for _, tt := range tests {
		ctrl := parser.FindControlByName(pkg, tt.searchName)
		if tt.expectFound {
			if ctrl == nil {
				t.Fatalf("expected to find control %s", tt.searchName)
			}
			if ctrl.Name != tt.searchName {
				t.Fatalf("found wrong control: expected %s, got %s", tt.searchName, ctrl.Name)
			}
		} else {
			if ctrl != nil {
				t.Fatalf("expected not to find control %s", tt.searchName)
			}
		}
	}
}

func TestExtractAssets(t *testing.T) {
	yamlContent := `Header:
  version: 1
ImportedPackages:
  - "~res:/UI/Screens/Battle/Shell.yaml"
  - "~res:/UI/Common.yaml"
Prototypes:
  - class: "UIControl"
    name: "Parent"
    prototype: "SharedUI/Header"
    components:
      Background:
        sprite: "~res:/Gfx/UI/bg.webp"
      StyleSheet:
        styles: "~res:/UI/main.style.yaml"
    children:
      - class: "UIControl"
        name: "Child"
        components:
          Background:
            sprite: "~res:/Gfx/UI/child.webp"`

	parser := NewParser()
	pkg, err := parser.Parse([]byte(yamlContent))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	assets := parser.ExtractAssets(pkg)
	if len(assets) == 0 {
		t.Fatal("expected assets to be extracted")
	}

	// Check that we extracted some assets
	hasImport := false
	hasSprite := false
	for _, asset := range assets {
		if asset == "~res:/UI/Screens/Battle/Shell.yaml" {
			hasImport = true
		}
		if asset == "~res:/Gfx/UI/bg.webp" {
			hasSprite = true
		}
	}

	if !hasImport {
		t.Fatal("expected to find imported package in assets")
	}
	if !hasSprite {
		t.Fatal("expected to find sprite in assets")
	}
}
