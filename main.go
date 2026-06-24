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
	db       *tiles.DB
	lat      float64
	lon      float64
	zoom     int
	width    int  // terminal columns
	height   int  // terminal rows
	showHelp bool // toggle help overlay with '?'
	frame    string
	status   string
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

func (m model) Init() tea.Cmd { return nil }

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

		// Toggle help overlay
		case "?":
			m.showHelp = !m.showHelp
			return m, nil

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

// View renders the full TUI as a magenta box:
//
//	╭────────────────────────────────────────╮
//	│  braille map canvas  (m.height-2 rows) │
//	╰─ z:12  N↑  28.6139°N  77.2090°E ──────╯
func (m model) View() string {
	bdr := lipgloss.NewStyle().Foreground(lipgloss.Color("201")) // magenta

	innerW := max(m.width-2, 0) // content columns between the │ chars

	// ── Top border ────────────────────────────────────────────────────────────
	top := bdr.Render("╭" + strings.Repeat("─", innerW) + "╮")

	// ── Main content: map frame or help screen ─────────────────────────────────
	var rawContent string
	if m.showHelp {
		rawContent = m.helpContent()
	} else {
		rawContent = m.frame
		if rawContent == "" {
			rawContent = strings.Repeat("\n", max(m.height-2, 1))
		}
	}

	// Wrap every content line with magenta side borders.
	// buf.Render() always ends each row with '\n'; TrimRight removes the trailing one
	// so Split gives exactly (m.height-2) elements with no spurious empty last entry.
	lines := strings.Split(strings.TrimRight(rawContent, "\n"), "\n")
	var framed strings.Builder
	for _, line := range lines {
		framed.WriteString(bdr.Render("│") + line + bdr.Render("│") + "\n")
	}

	// ── Bottom border with HUD ─────────────────────────────────────────────────
	hudText := m.hudText()

	// "╰─ " = 3 chars  │  " " before dashes = 1  │  "─╯" = 2  → 6 reserved chars
	available := m.width - 6
	if available < 0 {
		available = 0
	}
	hudRunes := []rune(hudText)
	if len(hudRunes) > available {
		hudRunes = hudRunes[:available]
		hudText = string(hudRunes)
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

	return top + "\n" + framed.String() + bottom
}

// helpContent returns the help screen text, padded to fill the inner canvas
// (m.width-2 columns × m.height-2 rows) so the side borders align correctly.
func (m model) helpContent() string {
	w := max(m.width-2, 0)
	h := max(m.height-2, 0)

	accent := lipgloss.NewStyle().Foreground(lipgloss.Color("109")) // sage teal
	key := lipgloss.NewStyle().Foreground(lipgloss.Color("222"))    // warm off-white
	dim := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))    // dark grey

	helpLines := []string{
		"",
		accent.Render("  NCR on Terminal") + dim.Render("  ─  keybindings"),
		"",
		accent.Render("  Navigation"),
		"    " + key.Render("↑ k w") + dim.Render("  pan north    ") + key.Render("↓ j s") + dim.Render("  pan south"),
		"    " + key.Render("← h a") + dim.Render("  pan west     ") + key.Render("→ l d") + dim.Render("  pan east"),
		"",
		accent.Render("  Zoom"),
		"    " + key.Render("+ =") + dim.Render("         zoom in     ") + key.Render("- _") + dim.Render("       zoom out"),
		"    " + key.Render("scroll ↑") + dim.Render("     zoom in     ") + key.Render("scroll ↓") + dim.Render("   zoom out"),
		"",
		accent.Render("  Map symbols"),
		"    " + key.Render("M") + dim.Render("  metro station    ") + key.Render("T") + dim.Render("  rail/train station"),
		"    " + key.Render("+") + dim.Render("  hospital         ") + key.Render("f") + dim.Render("  restaurant / café"),
		"    " + key.Render("g") + dim.Render("  fuel station"),
		"",
		accent.Render("  Other"),
		"    " + key.Render("?") + dim.Render("   toggle this help screen"),
		"    " + key.Render("q") + dim.Render("   quit"),
		"",
		dim.Render("  ─────────────────────────────────────────────────"),
		dim.Render("  Tip: for AMOLED appearance, set your terminal's"),
		dim.Render("  background color to #000000 (pure black)."),
	}

	var sb strings.Builder
	for i := 0; i < h; i++ {
		var line string
		if i < len(helpLines) {
			line = helpLines[i]
		}
		// Pad the line to exactly w visual columns so the right │ border aligns.
		pad := w - lipgloss.Width(line)
		if pad < 0 {
			pad = 0
		}
		sb.WriteString(line + strings.Repeat(" ", pad) + "\n")
	}
	return sb.String()
}

// hudText returns the plain HUD string (styling applied separately in View).
func (m model) hudText() string {
	zoom := fmt.Sprintf("z:%d", m.zoom)
	coords := fmt.Sprintf("%.4f°N  %.4f°E", m.lat, m.lon)
	scale := zoomToScale(m.zoom)
	loading := ""
	if m.status == "Rendering..." {
		loading = "⠿"
	}
	parts := []string{zoom, "N↑", coords, scale}
	if loading != "" {
		parts = append(parts, loading)
	}
	return strings.Join(parts, "  │  ")
}

// renderCmd triggers an async map render.
// pixelW accounts for the two side-border columns (│ on each side).
func (m model) renderCmd() tea.Cmd {
	db := m.db
	lat, lon := m.lat, m.lon
	zoom := m.zoom
	pixelW := (m.width - 2) * 2  // -2: one │ column on each side
	pixelH := (m.height - 2) * 4 // -2: top border row + bottom HUD row

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
