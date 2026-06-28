# demo.tape — ncr-on-terminal v1.0.0
# Delhi NCR in your terminal — braille vector tile map viewer
#
# Prerequisites:
#   go build -o ncr-on-terminal .
#   (mapdata/delhi-ncr.mbtiles must be present)
#
# Install vhs:  go install github.com/charmbracelet/vhs@latest
# Run:          vhs demo.tape
#
# Zoom keys used:
#   + / -   →  ±0.25  (coarse keyboard zoom)
#   , / .   →  ±0.10  (fine zoom, mirrors scroll wheel — see main.go)

Output demo.gif

Set Width        1280
Set Height       800
Set FontSize     20
Set Padding      10
Set TypingSpeed  40ms
Set WaitTimeout  30s
Set Theme        { "name": "AMOLED-Whimsy", "black": "#535178", "red": "#ef6487", "green": "#5eca89", "yellow": "#fdd877", "blue": "#65aef7", "magenta": "#aa7ff0", "cyan": "#43c1be", "white": "#ffffff", "brightBlack": "#535178", "brightRed": "#ef6487", "brightGreen": "#5eca89", "brightYellow": "#fdd877", "brightBlue": "#65aef7", "brightMagenta": "#aa7ff0", "brightCyan": "#43c1be", "brightWhite": "#ffffff", "background": "#000000", "foreground": "#b3b0d6", "selection": "#3d3c58", "cursor": "#b3b0d6" }
Set WindowBar    Colorful
Set BorderRadius 10

# ── Launch ────────────────────────────────────────────────────────────
Show
Type "./ncr-on-terminal mapdata/delhi-ncr.mbtiles"
Enter
Sleep 3s

# ── Opening overview — let the full NCR breathe ───────────────────────
Sleep 1.5s

# ── Pan east toward central Delhi ────────────────────────────────────
Right 3
Sleep 110ms
Right 3
Sleep 110ms
Up 2
Sleep 700ms

# ── Fine zoom in  (, = +0.1 per press) — smooth, continuous feel ──────
Type ","
Sleep 50ms
Type ","
Sleep 50ms
Type ","
Sleep 50ms
Type ","
Sleep 50ms
Type ","
Sleep 50ms
Type ","
Sleep 50ms
Type ","
Sleep 50ms
Type ","
Sleep 50ms
Type ","
Sleep 50ms
Type ","
Sleep 50ms
Type ","
Sleep 50ms
Type ","
Sleep 800ms

# ── Coarse keyboard zoom in  (+ = +0.25 per press) — snappy steps ────
Type "+"
Sleep 130ms
Type "+"
Sleep 130ms
Type "+"
Sleep 130ms
Type "+"
Sleep 800ms

# ── Pan south-west toward Dwarka / Gurugram border ───────────────────
Left 4
Sleep 110ms
Left 4
Sleep 110ms
Down 4
Sleep 110ms
Down 3
Sleep 700ms

# ── Fine zoom out  (. = -0.1 per press) ──────────────────────────────
Type "."
Sleep 50ms
Type "."
Sleep 50ms
Type "."
Sleep 50ms
Type "."
Sleep 50ms
Type "."
Sleep 50ms
Type "."
Sleep 50ms
Type "."
Sleep 50ms
Type "."
Sleep 50ms
Type "."
Sleep 50ms
Type "."
Sleep 50ms
Type "."
Sleep 50ms
Type "."
Sleep 800ms

# ── Coarse keyboard zoom out  (- = -0.25 per press) ──────────────────
Type "-"
Sleep 130ms
Type "-"
Sleep 130ms
Type "-"
Sleep 130ms
Type "-"
Sleep 800ms

# ── Help overlay ──────────────────────────────────────────────────────
Type "?"
Sleep 2.5s

# ── Close help → quit ─────────────────────────────────────────────────
Type "?"
Sleep 150ms
Type "q"
Sleep 200ms
