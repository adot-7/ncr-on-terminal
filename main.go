package main

import (
	"fmt"
	"os"

	"teaTui/geo"
	"teaTui/render"
	"teaTui/tiles"

	"github.com/charmbracelet/log"

	tea "github.com/charmbracelet/bubbletea"
)

// model holds all application state.
// Rule: only Update() modifies the model. Nothing else.
type model struct {
	db     *tiles.DB
	lat    float64
	lon    float64
	zoom   int
	width  int // terminal columns
	height int // terminal rows
	frame  string
	status string
}

// These are your custom Msg types.
// Any Go type can be a Msg. We define specific types for clarity.
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
	// Nothing to do at startup; we wait for WindowSizeMsg
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		// Terminal size is now known (or changed).
		m.width = msg.Width
		m.height = msg.Height
		m.status = "Rendering..."
		return m, m.renderCmd()

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		// Pan (move the viewport)
		case "up", "k":
			m.lat += geo.PanAmount(m.zoom)
			m.status = fmt.Sprintf("lat=%.4f lon=%.4f z=%d", m.lat, m.lon, m.zoom)
			return m, m.renderCmd()
		case "down", "j":
			m.lat -= geo.PanAmount(m.zoom)
			m.status = fmt.Sprintf("lat=%.4f lon=%.4f z=%d", m.lat, m.lon, m.zoom)
			return m, m.renderCmd()
		case "left", "h":
			m.lon -= geo.PanAmount(m.zoom)
			m.status = fmt.Sprintf("lat=%.4f lon=%.4f z=%d", m.lat, m.lon, m.zoom)
			return m, m.renderCmd()
		case "right", "l":
			m.lon += geo.PanAmount(m.zoom)
			m.status = fmt.Sprintf("lat=%.4f lon=%.4f z=%d", m.lat, m.lon, m.zoom)
			return m, m.renderCmd()

		// Zoom in/out
		case "+", "=":
			if m.zoom < 15 {
				m.zoom++
			} else {
				return m, nil
			}
			m.status = fmt.Sprintf("lat=%.4f lon=%.4f z=%d", m.lat, m.lon, m.zoom)
			return m, m.renderCmd()
		case "-", "_":
			if m.zoom > 5 {
				m.zoom--
			} else {
				return m, nil
			}
			m.status = fmt.Sprintf("lat=%.4f lon=%.4f z=%d", m.lat, m.lon, m.zoom)
			return m, m.renderCmd()
		}

	case frameReadyMsg:
		// The async render finished. Store the frame.
		m.frame = string(msg)
		return m, nil

	case statusMsg:
		m.status = string(msg)
		return m, nil
	}

	return m, nil
}

// View is called after every Update. It must be fast.
// We just return the last computed frame string.
func (m model) View() string {
	if m.frame == "" {
		return m.status + "\n"
	}
	// Append status line at the bottom
	return m.frame + "\x1b[0m" + m.status + " (q=quit +/-=zoom arrows=pan)"
}

// renderCmd returns a Cmd that renders the map in a goroutine.
// Capture all values by copy — the goroutine may run after the model has changed.
func (m model) renderCmd() tea.Cmd {
	db := m.db
	lat, lon := m.lat, m.lon
	zoom := m.zoom
	// Braille pixel dimensions
	pixelW := m.width * 2
	pixelH := (m.height - 1) * 4 // -1 for the status line

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
		tea.WithAltScreen(),       // Use the alternate screen buffer (full terminal)
		tea.WithMouseCellMotion(), // Optional: enable mouse for later
	)

	if _, err := p.Run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
