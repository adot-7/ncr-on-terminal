# demo.tape — record a demo GIF with vhs
#
# Install vhs: go install github.com/charmbracelet/vhs@latest
# Run:        vhs demo.tape
#
# Prerequisites: build the binary first:
#   go build -o ncr-on-terminal .
#   (and have mapdata/delhi-ncr.mbtiles present)

Output demo.gif
Set Width 1280
Set Height 800
Set FontSize 18
Set TypingSpeed 50ms
Set WaitTimeout 30s
Set Theme { "name": "Whimsy", "black": "#535178", "red": "#ef6487", "green": "#5eca89", "yellow": "#fdd877", "blue": "#65aef7", "magenta": "#aa7ff0", "cyan": "#43c1be", "white": "#ffffff", "brightBlack": "#535178", "brightRed": "#ef6487", "brightGreen": "#5eca89", "brightYellow": "#fdd877", "brightBlue": "#65aef7", "brightMagenta": "#aa7ff0", "brightCyan": "#43c1be", "brightWhite": "#ffffff", "background": "#000000", "foreground": "#b3b0d6", "selection": "#3d3c58", "cursor": "#b3b0d6" }
Set WindowBar Colorful
Set BorderRadius 10

Show
Type "./ncr-on-terminal mapdata/delhi-ncr-big.mbtiles"
Enter
Sleep 1s

# Pan east toward central Delhi
Right 4
Sleep 100ms
Down 2
Sleep 100ms

# Zoom in
Type "+"
Sleep 500ms
Type "+"
Sleep 500ms
Type "+"
Sleep 500ms
Type "+"
Sleep 500ms


# Zoom out
Type "-"
Sleep 500ms
Type "-"
Sleep 500ms
Type "-"
Sleep 500ms

# Pan west toward dwarka
Left 5
Sleep 100ms
Left 5
Sleep 100ms
Left 5
Sleep 100ms
Sleep 1s

# Zoom in
Type "+"
Sleep 500ms
Type "+"
Sleep 500ms
Type "+"
Sleep 500ms

# Zoom out
#Type "-"
#Sleep 500ms
#Type "-"
#Sleep 500ms
#Type "-"

# Toggle help screen
Type "?"
Sleep 2.5s

# Close help and quit
Type "?"
Sleep 50ms
Type "q"
Sleep 50ms
