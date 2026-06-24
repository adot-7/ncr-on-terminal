# demo.tape — record a demo GIF with vhs
#
# Install vhs: go install github.com/charmbracelet/vhs@latest
# Run:        vhs demo.tape
#
# Prerequisites: build the binary first:
#   go build -o ncr-on-terminal .
#   (and have mapdata/delhi-ncr.mbtiles present)

Output demo.gif
Set Width 1260
Set Height 650
Set FontSize 14
Set TypingSpeed 80ms
Set WaitTimeout 30s
Set Theme "Catppuccin Mocha"

Hide
Type "./ncr-on-terminal mapdata/delhi-ncr.mbtiles"
Show
Enter
Sleep 3s

# Pan east toward central Delhi
Right 4
Sleep 600ms
Down 2
Sleep 500ms

# Zoom in
Type "+"
Sleep 1.5s
Type "+"
Sleep 1.5s

# Pan around at higher zoom
Right 2
Sleep 600ms
Up 1
Sleep 600ms

# Zoom out
Type "-"
Sleep 1s

# Toggle help screen
Type "?"
Sleep 2.5s

# Close help and quit
Type "?"
Sleep 500ms
Type "q"
Sleep 500ms
