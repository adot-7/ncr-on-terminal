package style

import "teaTui/braille"

type LayerStyle struct {
	DrawFill   bool
	DrawLine   bool
	FillColor  int // xterm-256 color
	LineColor  int
	DrawLabel  bool // draw text labels for features in this layer
	LabelColor int  // xterm-256 color for label text
	MinZoom    int
	MaxZoom    int
	Priority   int // draw order (higher = drawn on top)
}

// StyleFor returns the rendering style for a given OpenMapTiles layer name and feature class.
func StyleFor(layerName, class string, zoom int) (LayerStyle, bool) {
	key := layerName
	if class != "" {
		key = layerName + "/" + class
	}

	styles := map[string]LayerStyle{
		// Water
		"water": {
			DrawFill: true, DrawLine: false,
			FillColor: braille.RGBToXterm256(68, 144, 196),
			MinZoom:   0, MaxZoom: 22, Priority: 10,
		},
		"waterway": {
			DrawFill: false, DrawLine: true,
			LineColor: braille.RGBToXterm256(68, 144, 196),
			MinZoom:   8, MaxZoom: 22, Priority: 11,
		},
		// Land cover
		"landcover/wood": {
			DrawFill:  true,
			FillColor: braille.RGBToXterm256(108, 168, 128),
			MinZoom:   7, MaxZoom: 22, Priority: 5,
		},
		"landcover/grass": {
			DrawFill:  true,
			FillColor: braille.RGBToXterm256(172, 208, 164),
			MinZoom:   7, MaxZoom: 22, Priority: 5,
		},
		// Roads (by class, most important first)
		// Motorway + trunk get labels starting at zoom 10/11 so you can see NH names
		"transportation/motorway": {
			DrawLine: true, DrawLabel: true,
			LineColor:  braille.RGBToXterm256(230, 80, 60),
			LabelColor: braille.RGBToXterm256(230, 80, 60),
			MinZoom:    5, MaxZoom: 22, Priority: 30,
		},
		"transportation/trunk": {
			DrawLine: true, DrawLabel: true,
			LineColor:  braille.RGBToXterm256(230, 150, 60),
			LabelColor: braille.RGBToXterm256(230, 150, 60),
			MinZoom:    6, MaxZoom: 22, Priority: 29,
		},
		"transportation/primary": {
			DrawLine: true, DrawLabel: true,
			LineColor:  braille.RGBToXterm256(230, 200, 60),
			LabelColor: braille.RGBToXterm256(230, 200, 60),
			MinZoom:    8, MaxZoom: 22, Priority: 28,
		},
		"transportation/secondary": {
			DrawLine:  true,
			LineColor: braille.RGBToXterm256(200, 200, 200),
			MinZoom:   10, MaxZoom: 22, Priority: 27,
		},
		"transportation/residential": {
			DrawLine:  true,
			LineColor: braille.RGBToXterm256(160, 160, 160),
			MinZoom:   12, MaxZoom: 22, Priority: 26,
		},
		"transportation/service": {
			DrawLine:  true,
			LineColor: braille.RGBToXterm256(120, 120, 120),
			MinZoom:   13, MaxZoom: 22, Priority: 25,
		},
		// Railways
		"transportation/rail": {
			DrawLine:  true,
			LineColor: braille.RGBToXterm256(80, 80, 80),
			MinZoom:   8, MaxZoom: 22, Priority: 20,
		},
		// Buildings
		"building": {
			DrawFill: true, DrawLine: true,
			FillColor: braille.RGBToXterm256(180, 160, 140),
			LineColor: braille.RGBToXterm256(150, 130, 110),
			MinZoom:   13, MaxZoom: 22, Priority: 40,
		},
		// Points of interest — labels only, no geometry
		"poi": {
			DrawLabel:  true,
			LabelColor: braille.RGBToXterm256(255, 255, 255),
			MinZoom:    14, MaxZoom: 22, Priority: 50,
		},
		// Place names (cities, towns, villages)
		"place": {
			DrawLabel:  true,
			LabelColor: braille.RGBToXterm256(255, 220, 100),
			MinZoom:    8, MaxZoom: 22, Priority: 60,
		},
	}

	if s, ok := styles[key]; ok && zoom >= s.MinZoom && zoom <= s.MaxZoom {
		return s, true
	}
	// Try without the class suffix for fallback
	if class != "" {
		if s, ok := styles[layerName]; ok && zoom >= s.MinZoom && zoom <= s.MaxZoom {
			return s, true
		}
	}
	return LayerStyle{}, false
}
