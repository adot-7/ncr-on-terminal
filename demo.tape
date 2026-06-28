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
Set FontSize 20
Set Padding 15
Set TypingSpeed 40ms
Set WaitTimeout 30s
Set Theme { "name": "Whimsy", "black": "#535178", "red": "#ef6487", "green": "#5eca89", "yellow": "#fdd877", "blue": "#65aef7", "magenta": "#aa7ff0", "cyan": "#43c1be", "white": "#ffffff", "brightBlack": "#535178", "brightRed": "#ef6487", "brightGreen": "#5eca89", "brightYellow": "#fdd877", "brightBlue": "#65aef7", "brightMagenta": "#aa7ff0", "brightCyan": "#43c1be", "brightWhite": "#ffffff", "background": "#000000", "foreground": "#b3b0d6", "selection": "#3d3c58", "cursor": "#b3b0d6" }
Set WindowBar Colorful
Set BorderRadius 10

Show
Type "./ncr-on-terminal mapdata/delhi-ncr.mbtiles"
Enter
Sleep 2.5s

# Pan east toward central Delhi
Up 2
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
Down 5
Sleep 100ms


# Zoom in
Type "+"
Sleep 500ms
Type "+"
Sleep 500ms
Type "+"
Sleep 500ms

# Zoom out
Type "-"
Sleep 200ms
Type "-"
Sleep 200ms
Type "-"
Sleep 200ms
Type "-"
Sleep 200ms
Type "-"
Sleep 200ms
Type "-"
Sleep 200ms
Type "-"
Sleep 200ms
Type "-"
Sleep 200ms
Type "-"
Sleep 750ms

# Toggle help screen
Type "?"
Sleep 2s

# Close help and quit
Type "?"
Sleep 100ms
Type "q"
Sleep 50ms
