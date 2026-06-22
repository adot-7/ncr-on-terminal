package render

import (
	"teaTui/braille"
	"teaTui/geo"
	"teaTui/style"
	"teaTui/tiles"

	"github.com/charmbracelet/log"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/mvt"
	"github.com/paulmach/orb/simplify"
)

// RenderRequest bundles everything needed for one frame.
type RenderRequest struct {
	DB       *tiles.DB
	Lat, Lon float64
	Zoom     int
	PixelW   int // braille pixel width (= terminal cols * 2)
	PixelH   int // braille pixel height (= terminal rows * 4)
}

// checks each layer in mvt.Layer to find if a layer with name exists. returns nil otherwise
func findLayer(layers mvt.Layers, name string) *mvt.Layer {
	for _, layer := range layers {
		if layer != nil && layer.Name == name {
			return layer
		}
	}
	return nil
}

// Render builds a full frame string from the given request.
// This is called inside a Cmd (goroutine), so it can block.
func Render(req RenderRequest) string {
	buf := braille.New(req.PixelW/2, req.PixelH/4)
	buf.Clear()

	// Step 1: Determine which tiles we need
	vp := geo.Viewport{
		Lat: req.Lat, Lon: req.Lon, Zoom: req.Zoom,
		PixelW: req.PixelW, PixelH: req.PixelH,
	}
	tileRequests := vp.ComputeTiles()

	// Step 2: Define the draw order for layers
	// We iterate layers in this order so that buildings appear on top of roads,
	// roads appear on top of land cover, etc.
	layerOrder := []string{
		"landcover", "landuse", "water", "waterway",
		"boundary", "transportation", "building", "poi",
	}

	// Step 3: Load and draw each tile
	isFirstTile := true
	for _, req2 := range tileRequests {

		data, err := req.DB.ReadTile(req2.Z, req2.X, req2.Y)
		if err != nil || data == nil {
			continue // tile missing or read error — just skip it
		}

		layers, err := mvt.Unmarshal(data)
		if err != nil {
			continue
		}
		if isFirstTile {
			for _, l := range layers {
				if l != nil {
					log.Debugf("Layer and features:%s (%d)", l.Name, len(l.Features))
				}
			}
		}

		// Step 4: Draw layers in order
		for _, layerName := range layerOrder {
			// layer, ok := layers[layerName]
			layer := findLayer(layers, layerName)
			if layer == nil {
				continue
			}

			// Simplify geometry to reduce points we don't need at this resolution
			// tolerance: roughly 0.5 screen pixels worth of tile units
			tolerance := 4096.0 / float64(256) * 0.5 // ≈ 8 tile units per half-pixel
			simplifier := simplify.DouglasPeucker(tolerance)

			for _, feature := range layer.Features {
				// Get class for more specific style lookup
				class, _ := feature.Properties["class"].(string)

				st, ok := style.StyleFor(layerName, class, req.Zoom)
				if !ok {
					continue // don't draw this feature
				}

				// Simplify before transforming
				simplified := simplifier.Simplify(feature.Geometry)

				// Draw it
				drawGeometry(buf, simplified, req2, st)
			}
		}
	}

	return buf.Render()
}

// drawGeometry dispatches to the appropriate draw method based on geometry type.
func drawGeometry(buf *braille.Buffer, g orb.Geometry, req geo.TileRequest, st style.LayerStyle) {
	switch geom := g.(type) {
	case orb.LineString:
		if st.DrawLine {
			drawLineString(buf, geom, req, st.LineColor)
		}
	case orb.MultiLineString:
		if st.DrawLine {
			for _, ls := range geom {
				drawLineString(buf, ls, req, st.LineColor)
			}
		}
	case orb.Polygon:
		if st.DrawFill {
			drawPolygon(buf, geom, req, st.FillColor)
		}
		if st.DrawLine {
			drawLineString(buf, orb.LineString(geom[0]), req, st.LineColor)
		}
	case orb.MultiPolygon:
		for _, poly := range geom {
			if st.DrawFill {
				drawPolygon(buf, poly, req, st.FillColor)
			}
			if st.DrawLine {
				drawLineString(buf, orb.LineString(poly[0]), req, st.LineColor)
			}
		}
	case orb.Point:
		px, py := tileToPixel(geom[0], geom[1], req)
		buf.SetPixel(px, py, st.LineColor)
	}
}

// tileToPixel converts a tile-space coordinate (x, y in [0, 4096]) to a braille pixel
// position on screen, given the tile's pixel offset and scale.
//
// The key transform is:
//
//	screen_pixel = tile_offset_pixel + tile_coord * scale
//
// Where:
//
//	tile_offset_pixel = where this tile's (0,0) lands on the braille pixel grid
//	tile_coord = the coordinate within the tile (0 to 4096)
//	scale = braille pixels per tile-space unit
func tileToPixel(tileX, tileY float64, req geo.TileRequest) (px, py int) {
	px = req.PixelOffsetX + int(tileX*req.Scale)
	py = req.PixelOffsetY + int(tileY*req.Scale)
	return
}

func drawLineString(buf *braille.Buffer, ls orb.LineString, req geo.TileRequest, color int) {
	if len(ls) < 2 {
		return
	}
	xs := make([]int, len(ls))
	ys := make([]int, len(ls))
	for i, pt := range ls {
		xs[i], ys[i] = tileToPixel(pt[0], pt[1], req)
	}
	buf.DrawPolyline(xs, ys, color)
}

func drawPolygon(buf *braille.Buffer, poly orb.Polygon, req geo.TileRequest, color int) {
	if len(poly) == 0 {
		return
	}
	// Draw the outer ring filled
	ring := poly[0]
	xs := make([]int, len(ring))
	ys := make([]int, len(ring))
	for i, pt := range ring {
		xs[i], ys[i] = tileToPixel(pt[0], pt[1], req)
	}
	buf.FillPolygon(xs, ys, color)

	// Erase holes with background color (0 = terminal default)
	// This is the simple approach: draw over the hole with background
	for _, hole := range poly[1:] {
		hxs := make([]int, len(hole))
		hys := make([]int, len(hole))
		for i, pt := range hole {
			hxs[i], hys[i] = tileToPixel(pt[0], pt[1], req)
		}
		buf.FillPolygon(hxs, hys, 0) // 0 = no color = clears the fill
	}
}
