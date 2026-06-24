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

// Label holds a text label to be written into the braille buffer's text overlay.
type Label struct {
	Text  string
	ColX  int // terminal column (0-indexed)
	RowY  int // terminal row   (0-indexed)
	Color int // xterm-256 index; 0 = terminal default
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

	// Step 2: Define the draw order for layers.
	// "place" is last so city/town names render on top of everything.
	layerOrder := []string{
		"landcover", "landuse", "water", "waterway",
		"boundary", "transportation", "building", "poi", "place",
	}

	// Step 3: Load and draw each tile, collecting labels along the way.
	var labels []Label
	seenLabels := make(map[string]bool)

	isFirstTile := true
	for _, req2 := range tileRequests {

		data, err := req.DB.ReadTile(req2.Z, req2.X, req2.Y)
		if err != nil || data == nil {
			continue
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
			isFirstTile = false
		}

		for _, layerName := range layerOrder {
			layer := findLayer(layers, layerName)
			if layer == nil {
				continue
			}

			tolerance := 4096.0 / float64(256) * 0.5
			simplifier := simplify.DouglasPeucker(tolerance)

			for _, feature := range layer.Features {
				class, _ := feature.Properties["class"].(string)

				st, ok := style.StyleFor(layerName, class, req.Zoom)
				if !ok {
					continue
				}

				simplified := simplifier.Simplify(feature.Geometry)
				drawGeometry(buf, simplified, req2, st)

				if st.DrawLabel {
					name, _ := feature.Properties["name"].(string)
					if name == "" {
						name, _ = feature.Properties["name_en"].(string)
					}
					if name != "" {
						if layerName == "transportation" && seenLabels[name] {
							continue
						}
						if tx, ty, ok2 := featurePoint(simplified); ok2 {
							px, py := tileToPixel(tx, ty, req2)
							col, row := px/2, py/4
							labels = append(labels, Label{
								Text:  name,
								ColX:  col,
								RowY:  row,
								Color: st.LabelColor,
							})
							if layerName == "transportation" {
								seenLabels[name] = true
							}
						}
					}
				}
			}
		}
	}

	termW := req.PixelW / 2
	termH := req.PixelH / 4

	// --- DIAGNOSTIC ---
	// 1. Log how many labels were collected and where the first few land.
	log.Debugf("LABELS: collected=%d  termW=%d termH=%d", len(labels), termW, termH)
	for i, l := range labels {
		if i >= 8 {
			break
		}
		log.Debugf("  label[%d] %q  col=%d row=%d", i, l.Text, l.ColX, l.RowY)
	}

	// 2. Hardcoded test: paint a bright magenta "*" at cell (3,3).
	//    If you can see it on screen, the text-overlay mechanism is working and
	//    the problem is purely in label data / coordinates.
	//    Remove this block once labels are confirmed working.
	testColor := braille.RGBToXterm256(255, 0, 255) // magenta
	for i, r := range []rune("*LABEL*") {
		buf.SetText(3+i, 3, r, testColor)
	}
	// --- END DIAGNOSTIC ---

	writeLabelsToBuffer(buf, labels, termW, termH)
	return buf.Render()
}

// writeLabelsToBuffer writes collected label text directly into the braille buffer's
// text overlay layer. A per-cell occupancy grid prevents overlapping labels.
func writeLabelsToBuffer(buf *braille.Buffer, labels []Label, termW, termH int) {
	occupied := make(map[[2]int]bool)
	written, skipped := 0, 0

	for _, l := range labels {
		if l.ColX < 0 || l.RowY < 0 || l.RowY >= termH {
			log.Debugf("LABEL skip OOB: %q col=%d row=%d", l.Text, l.ColX, l.RowY)
			skipped++
			continue
		}
		maxLen := termW - l.ColX
		if maxLen <= 0 {
			log.Debugf("LABEL skip maxLen: %q col=%d termW=%d", l.Text, l.ColX, termW)
			skipped++
			continue
		}

		runes := []rune(l.Text)
		if len(runes) > maxLen {
			runes = runes[:maxLen]
		}

		collision := false
		for i := range runes {
			if occupied[[2]int{l.ColX + i, l.RowY}] {
				collision = true
				break
			}
		}
		if collision {
			skipped++
			continue
		}

		for i, r := range runes {
			col := l.ColX + i
			occupied[[2]int{col, l.RowY}] = true
			buf.SetText(col, l.RowY, r, l.Color)
		}
		written++
	}

	log.Debugf("LABELS: written=%d skipped=%d", written, skipped)
}

// featurePoint returns a representative tile-space coordinate for label placement.
func featurePoint(g orb.Geometry) (x, y float64, ok bool) {
	switch geom := g.(type) {
	case orb.Point:
		return geom[0], geom[1], true
	case orb.LineString:
		if len(geom) == 0 {
			return
		}
		mid := geom[len(geom)/2]
		return mid[0], mid[1], true
	case orb.MultiLineString:
		if len(geom) == 0 || len(geom[0]) == 0 {
			return
		}
		mid := geom[0][len(geom[0])/2]
		return mid[0], mid[1], true
	case orb.Polygon:
		if len(geom) == 0 || len(geom[0]) == 0 {
			return
		}
		ring := geom[0]
		var sx, sy float64
		for _, pt := range ring {
			sx += pt[0]
			sy += pt[1]
		}
		n := float64(len(ring))
		return sx / n, sy / n, true
	case orb.MultiPolygon:
		if len(geom) == 0 || len(geom[0]) == 0 || len(geom[0][0]) == 0 {
			return
		}
		ring := geom[0][0]
		var sx, sy float64
		for _, pt := range ring {
			sx += pt[0]
			sy += pt[1]
		}
		n := float64(len(ring))
		return sx / n, sy / n, true
	}
	return
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
		if st.DrawLine {
			px, py := tileToPixel(geom[0], geom[1], req)
			buf.SetPixel(px, py, st.LineColor)
		}
	}
}

// tileToPixel converts a tile-space coordinate (x, y in [0, 4096]) to a braille pixel
// position on screen, given the tile's pixel offset and scale.
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
	ring := poly[0]
	xs := make([]int, len(ring))
	ys := make([]int, len(ring))
	for i, pt := range ring {
		xs[i], ys[i] = tileToPixel(pt[0], pt[1], req)
	}
	buf.FillPolygon(xs, ys, color)

	for _, hole := range poly[1:] {
		hxs := make([]int, len(hole))
		hys := make([]int, len(hole))
		for i, pt := range hole {
			hxs[i], hys[i] = tileToPixel(pt[0], pt[1], req)
		}
		buf.FillPolygon(hxs, hys, 0)
	}
}
