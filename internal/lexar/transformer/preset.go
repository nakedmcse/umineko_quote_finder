package transformer

import (
	"strconv"
	"strings"

	"umineko_quote/internal/lexar/ast"
)

// PresetContext holds preset information collected from a script.
type PresetContext struct {
	SemanticPresets map[string]string // e.g. "1" -> "red-truth"
	DynamicColours  map[string]string // e.g. "41" -> "#FFAA00"
}

// DefaultSemanticPresets returns the built-in semantic preset mappings.
func DefaultSemanticPresets() map[string]string {
	return map[string]string{
		"1": "red-truth",
		"2": "blue-truth",
	}
}

// DefaultDynamicColours returns the built-in dynamic colour mappings.
// These can be overridden by preset_define lines in a script.
func DefaultDynamicColours() map[string]string {
	return map[string]string{
		"41": "#FFAA00", // Gold text
		"42": "#AA71FF", // Purple text
	}
}

// NewPresetContext creates a PresetContext with default semantic presets and dynamic colours.
func NewPresetContext() *PresetContext {
	return &PresetContext{
		SemanticPresets: DefaultSemanticPresets(),
		DynamicColours:  DefaultDynamicColours(),
	}
}

// CollectFromScript extracts dynamic colours from preset_define lines in a script.
// Colours from the script override the defaults.
func (p *PresetContext) CollectFromScript(script *ast.Script) {
	p.DynamicColours = DefaultDynamicColours()

	for _, line := range script.Lines {
		if preset, ok := line.(*ast.PresetDefineLine); ok {
			presetID := strconv.Itoa(preset.ID)

			// Skip semantic presets (they use classes, not inline colours)
			if _, isSemantic := p.SemanticPresets[presetID]; isSemantic {
				continue
			}

			colour := strings.ToUpper(preset.Colour)
			if colour == "#FFFFFF" || colour == "" {
				continue
			}

			p.DynamicColours[presetID] = colour
		}
	}
}

// GetSemanticClass returns the CSS class for a semantic preset, or empty string if not found.
func (p *PresetContext) GetSemanticClass(presetID string) string {
	return p.SemanticPresets[presetID]
}

// GetDynamicColour returns the colour for a dynamic preset, or empty string if not found.
func (p *PresetContext) GetDynamicColour(presetID string) string {
	return p.DynamicColours[presetID]
}
