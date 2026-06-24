package main

import (
	"fmt"
	"os"
	"strings"

	"teaTui/geo"
	"teaTui/render"
	"teaTui/tiles"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/log"

	tea "github.com/charmbracelet/bubbletea"
)

// model holds all application state.
// Rule: only Update() modifies the model.
type model struct {
	db     *tiles.DB
	lat    float64
	lon    float64
	zoom   int
	width  int // terminal columns (used for map canvas and border)
	height int // terminal rows   (used for map canvas and border)
	frame  string
	status string
}

type frameReadyMsg string
type statusMsg string

func initialModel(db *tiles.DB) model {
	return model{
		db:     db,
		lat:    28.6139, // Delhi center
		lon:    77.2090,
		zoom:   12,
		status: "Waiting for terminal size...",
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.status = "Rendering..."
		return m, m.renderCmd()

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		// Pan
		case "up", "k", "w":
			m.lat += geo.PanAmount(m.zoom)
			m.status = fmt.Sprintf("lat=%.4f lon=%.4f z=%d", m.lat, m.lon, m.zoom)
			return m, m.renderCmd()
		case "down", "j", "s":
			m.lat -= geo.PanAmount(m.zoom)
			m.status = fmt.Sprintf("lat=%.4f lon=%.4f z=%d", m.lat, m.lon, m.zoom)
			return m, m.renderCmd()
		case "left", "h", "a":
			m.lon -= geo.PanAmount(m.zoom)
			m.status = fmt.Sprintf("lat=%.4f lon=%.4f z=%d", m.lat, m.lon, m.zoom)
			return m, m.renderCmd()
		case "right", "l", "d":
			m.lon += geo.PanAmount(m.zoom)
			m.status = fmt.Sprintf("lat=%.4f lon=%.4f z=%d", m.lat, m.lon, m.zoom)
			return m, m.renderCmd()

		// Zoom
		case "+", "=":
			if m.zoom < 15 {
				m.zoom++
				m.status = fmt.Sprintf("lat=%.4f lon=%.4f z=%d", m.lat, m.lon, m.zoom)
				return m, m.renderCmd()
			}
		case "-", "_":
			if m.zoom > 5 {
				m.zoom--
				m.status = fmt.Sprintf("lat=%.4f lon=%.4f z=%d", m.lat, m.lon, m.zoom)
				return m, m.renderCmd()
			}
		}

	// Mouse wheel → zoom in/out
	case tea.MouseMsg:
		switch msg.Button {
		case tea.MouseButtonWheelUp:
			if m.zoom < 15 {
				m.zoom++
				m.status = fmt.Sprintf("lat=%.4f lon=%.4f z=%d", m.lat, m.lon, m.zoom)
				return m, m.renderCmd()
			}
		case tea.MouseButtonWheelDown:
			if m.zoom > 5 {
				m.zoom--
				m.status = fmt.Sprintf("lat=%.4f lon=%.4f z=%d", m.lat, m.lon, m.zoom)
				return m, m.renderCmd()
			}
		}

	case frameReadyMsg:
		m.frame = string(msg)
		return m, nil

	case statusMsg:
		m.status = string(msg)
		return m, nil
	}

	return m, nil
}

// View renders the full TUI:
//
//	╭──────── border top ────────╮
//	│  braille map canvas        │  ← (m.height - 2) rows
//	╰─ z:12  N↑  28.61°N … ─────╯
func (m model) View() string {
	frame := m.frame
	if frame == "" {
		// Placeholder while waiting for the first render — fill with blank lines.
		frame = strings.Repeat("\n", max(m.height-2, 1))
	}

	bdr := lipgloss.NewStyle().Foreground(lipgloss.Color("201")) // magenta

	// ── Top border ────────────────────────────────────────────────────────────
	innerW := m.width - 2 // space between ╭ and ╮
	if innerW < 0 {
		innerW = 0
	}
	top := bdr.Render("╭" + strings.Repeat("─", innerW) + "╮")

	// ── Bottom border with HUD content ────────────────────────────────────────
	hudText := m.hudText()

	// Available space for HUD text inside "╰─ … ─╯"
	// "╰─ " = 3 chars + " " before dashes + "─╯" = 2 chars → 6 total fixed chars
	available := m.width - 6
	if available < 0 {
		available = 0
	}
	runes := []rune(hudText)
	if len(runes) > available {
		runes = runes[:available]
		hudText = string(runes)
	}
	padLen := available - len([]rune(hudText))
	if padLen < 0 {
		padLen = 0
	}

	var hudStyled string
	if m.status == "Rendering..." {
		hudStyled = mutedStyle.Foreground(lipgloss.Color("240")).Render(hudText)
	} else {
		hudStyled = mutedStyle.Render(hudText)
	}

	bottom := bdr.Render("╰─ ") + hudStyled + bdr.Render(" "+strings.Repeat("─", padLen)+"╯")

	// frame ends with '\n' for every row (from buf.Render()), so no extra '\n' needed.
	return top + "\n" + frame + bottom
}

// hudText returns the plain HUD string (no lipgloss styling applied).
func (m model) hudText() string {
	zoom := fmt.Sprintf("z:%d", m.zoom)
	compass := "N↑"
	coords := fmt.Sprintf("%.4f°N  %.4f°E", m.lat, m.lon)
	scale := zoomToScale(m.zoom)
	loading := ""
	if m.status == "Rendering..." {
		loading = "⠿ rendering"
	}
	parts := []string{zoom, compass, coords, scale}
	if loading != "" {
		parts = append(parts, loading)
	}
	return strings.Join(parts, "  │  ")
}

// renderCmd triggers an async map render.
func (m model) renderCmd() tea.Cmd {
	db := m.db
	lat, lon := m.lat, m.lon
	zoom := m.zoom
	pixelW := m.width * 2
	// -2 rows: 1 for top border bar + 1 for bottom HUD bar
	pixelH := (m.height - 2) * 4

	return func() tea.Msg {
		frame := render.Render(render.RenderRequest{
			DB:     db,
			Lat:    lat,
			Lon:    lon,
			Zoom:   zoom,
			PixelW: pixelW,
			PixelH: pixelH,
		})
		return frameReadyMsg(frame)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: mapscii <path-to.mbtiles>")
	}

	db, err := tiles.Open(os.Args[1])
	if err != nil {
		log.Fatalf("Failed to open MBTiles: %v", err)
	}
	defer db.Close()

	f, err := os.OpenFile("trip.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		panic(err)
	}
	log.SetOutput(f)
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)

	p := tea.NewProgram(
		initialModel(db),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

var mutedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))

func zoomToScale(zoom int) string {
	scales := map[int]string{
		5: "~500km", 6: "~250km", 7: "~125km",
		8: "~60km", 9: "~30km", 10: "~15km",
		11: "~7km", 12: "~3.5km", 13: "~1.8km",
		14: "~900m", 15: "~450m", 16: "~225m",
		17: "~110m",
	}
	if s, ok := scales[zoom]; ok {
		return s
	}
	return ""
}
