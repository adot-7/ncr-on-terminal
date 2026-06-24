package style

import "teaTui/braille"

type LayerStyle struct {
	DrawFill    bool
	DrawLine    bool
	FillColor   int    // xterm-256 fill color
	LineColor   int    // xterm-256 line color
	DrawLabel   bool   // emit a label for this feature
	LabelColor  int    // xterm-256 label text color
	LabelSymbol string // if non-empty, render this single char instead of the feature name
	MinZoom     int
	MaxZoom     int
	Priority    int // draw order (higher = drawn on top)
}

// StyleFor returns the rendering style for a given layer name and class.
// Call with class = "parentClass/subclass" to get POI subclass-specific styles.
func StyleFor(layerName, class string, zoom int) (LayerStyle, bool) {
	key := layerName
	if class != "" {
		key = layerName + "/" + class
	}

	styles := map[string]LayerStyle{
		// ── Water ────────────────────────────────────────────────────────────────
		"water": {
			DrawFill:  true,
			FillColor: braille.RGBToXterm256(20, 100, 220), // vivid deep blue
			MinZoom:   0, MaxZoom: 22, Priority: 10,
		},
		"waterway": {
			DrawLine:  true,
			LineColor: braille.RGBToXterm256(30, 120, 230),
			MinZoom:   8, MaxZoom: 22, Priority: 11,
		},

		// ── Land cover ───────────────────────────────────────────────────────────
		"landcover/wood": {
			DrawFill:  true,
			FillColor: braille.RGBToXterm256(40, 160, 60), // vivid forest green
			MinZoom:   7, MaxZoom: 22, Priority: 5,
		},
		"landcover/grass": {
			DrawFill:  true,
			FillColor: braille.RGBToXterm256(130, 210, 100), // fresh grass green
			MinZoom:   7, MaxZoom: 22, Priority: 5,
		},

		// ── Roads — geometry only; names are in transportation_name ──────────────
		"transportation/motorway": {
			DrawLine:  true,
			LineColor: braille.RGBToXterm256(255, 80, 20), // strong orange-red
			MinZoom:   5, MaxZoom: 22, Priority: 30,
		},
		"transportation/trunk": {
			DrawLine:  true,
			LineColor: braille.RGBToXterm256(255, 160, 0), // bright amber
			MinZoom:   6, MaxZoom: 22, Priority: 29,
		},
		"transportation/primary": {
			DrawLine:  true,
			LineColor: braille.RGBToXterm256(255, 210, 0), // vivid yellow
			MinZoom:   8, MaxZoom: 22, Priority: 28,
		},
		"transportation/secondary": {
			DrawLine:  true,
			LineColor: braille.RGBToXterm256(220, 210, 100), // warm yellow-white
			MinZoom:   10, MaxZoom: 22, Priority: 27,
		},
		"transportation/residential": {
			DrawLine:  true,
			LineColor: braille.RGBToXterm256(180, 180, 180),
			MinZoom:   12, MaxZoom: 22, Priority: 26,
		},
		"transportation/service": {
			DrawLine:  true,
			LineColor: braille.RGBToXterm256(130, 130, 130),
			MinZoom:   13, MaxZoom: 22, Priority: 25,
		},
		"transportation/rail": {
			DrawLine:  true,
			LineColor: braille.RGBToXterm256(90, 90, 90),
			MinZoom:   8, MaxZoom: 22, Priority: 20,
		},

		// ── Buildings ────────────────────────────────────────────────────────────
		"building": {
			DrawFill: true, DrawLine: true,
			FillColor: braille.RGBToXterm256(210, 155, 100), // warm clay
			LineColor: braille.RGBToXterm256(170, 115, 70),
			MinZoom:   13, MaxZoom: 22, Priority: 40,
		},

		// ── Road name labels (separate layer in this tile set) ───────────────────
		// "transportation_name/motorway": {
		// 	DrawLabel:  true,
		// 	LabelColor: braille.RGBToXterm256(255, 80, 20),
		// 	MinZoom:    7, MaxZoom: 22, Priority: 30,
		// },
		// "transportation_name/trunk": {
		// 	DrawLabel:  true,
		// 	LabelColor: braille.RGBToXterm256(255, 160, 0),
		// 	MinZoom:    8, MaxZoom: 22, Priority: 29,
		// },
		// "transportation_name/primary": {
		// 	DrawLabel:  true,
		// 	LabelColor: braille.RGBToXterm256(255, 210, 0),
		// 	MinZoom:    10, MaxZoom: 22, Priority: 28,
		// },
		// "transportation_name/secondary": {
		// 	DrawLabel:  true,
		// 	LabelColor: braille.RGBToXterm256(220, 210, 100),
		// 	MinZoom:    12, MaxZoom: 22, Priority: 27,
		// },
		// "transportation_name": {
		// 	DrawLabel:  true,
		// 	LabelColor: braille.RGBToXterm256(170, 170, 100),
		// 	MinZoom:    13, MaxZoom: 22, Priority: 25,
		// },

		// ── Place names ──────────────────────────────────────────────────────────
		"place": {
			DrawLabel:  true,
			LabelColor: braille.RGBToXterm256(255, 235, 50), // bright gold
			MinZoom:    8, MaxZoom: 22, Priority: 60,
		},

		// ── POI — allowlist only; anything not listed here is silently ignored ───
		//
		// Metro / subway stations
		"poi/railway/subway": {
			DrawLabel:   true,
			LabelSymbol: "M",
			LabelColor:  braille.RGBToXterm256(0, 210, 255), // bright cyan
			MinZoom:     12, MaxZoom: 22, Priority: 72,
		},
		// Train / mainline stations
		"poi/railway/station": {
			DrawLabel:   true,
			LabelSymbol: "T",
			LabelColor:  braille.RGBToXterm256(120, 180, 255), // light blue
			MinZoom:     12, MaxZoom: 22, Priority: 71,
		},
		// Any other railway poi (tram, monorail, halt …)
		// "poi/railway": {
		// 	DrawLabel:   true,
		// 	LabelSymbol: "TT",
		// 	LabelColor:  braille.RGBToXterm256(100, 160, 240),
		// 	MinZoom:     13, MaxZoom: 22, Priority: 68,
		// },
		// Hospitals
		"poi/health/hospital": {
			DrawLabel:   true,
			LabelSymbol: "🏥",
			LabelColor:  braille.RGBToXterm256(255, 70, 70), // vivid red
			MinZoom:     12, MaxZoom: 22, Priority: 65,
		},
		// Clinics / other health
		"poi/health": {
			DrawLabel:   true,
			LabelSymbol: "🏥",
			LabelColor:  braille.RGBToXterm256(255, 140, 140),
			MinZoom:     14, MaxZoom: 22, Priority: 62,
		},
		// Restaurants / food
		"poi/food": {
			DrawLabel:   true,
			LabelSymbol: "🍴",
			LabelColor:  braille.RGBToXterm256(255, 185, 50), // warm amber
			MinZoom:     15, MaxZoom: 22, Priority: 50,
		},
		// Fuel / gas stations
		"poi/fuel": {
			DrawLabel:   true,
			LabelSymbol: "⛽",
			LabelColor:  braille.RGBToXterm256(80, 220, 80), // vivid green
			MinZoom:     14, MaxZoom: 22, Priority: 52,
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
