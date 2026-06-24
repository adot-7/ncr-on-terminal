# NCR on Terminal

> Delhi NCR, rendered in Braille characters — entirely in your terminal.

<!-- Replace this line with the GIF once you've recorded it with demo.tape -->
![NCR on Terminal demo](demo.gif)

A terminal map viewer for the Delhi National Capital Region built with [BubbleTea](https://github.com/charmbracelet/bubbletea). Roads, water, forests, buildings, metro stations, and place labels — all drawn using Unicode Braille characters for sub-character pixel resolution.

---

## Try it instantly (no install required)

```sh
ssh ncr.akashparashar.dev -p 2222
```

> **Note:** The SSH demo server may not always be running. See [Local install](#install) to run it yourself.

---

## Install

**Download a binary** from [Releases](https://github.com/adot-7/ncr-on-terminal/releases/latest):

| Platform | Download |
|---|---|
| macOS (Apple Silicon) | `ncr-on-terminal_*_darwin_arm64.tar.gz` |
| macOS (Intel) | `ncr-on-terminal_*_darwin_amd64.tar.gz` |
| Linux (x86_64) | `ncr-on-terminal_*_linux_amd64.tar.gz` |
| Linux (ARM64) | `ncr-on-terminal_*_linux_arm64.tar.gz` |
| Windows | `ncr-on-terminal_*_windows_amd64.zip` |

Extract and run. The binary has no dependencies — no Docker, no runtime, nothing else.

**Or build from source:**

```sh
git clone https://github.com/adot-7/ncr-on-terminal
cd ncr-on-terminal
go build -o ncr-on-terminal .
```

Requires Go 1.21+.

---

## Setup: getting an MBTiles file

The app needs an [OpenMapTiles](https://openmaptiles.org)-compatible `.mbtiles` vector tile file. This file is not included in the repo (it's ~50 MB).

### Recommended: generate with Planetiler

[Planetiler](https://github.com/onthegomap/planetiler) generates an MBTiles file from OpenStreetMap data in minutes.

```sh
# 1. Download the India OSM extract (~720 MB)
wget https://download.geofabrik.de/asia/india-latest.osm.pbf

# 2. Download Planetiler (requires Java 17+)
wget https://github.com/onthegomap/planetiler/releases/latest/download/planetiler.jar

# 3. Generate Delhi NCR tiles (~50 MB output, takes 3–5 min)
java -jar planetiler.jar \
  --osm-path=india-latest.osm.pbf \
  --bounds=76.84,28.30,77.35,29.00 \
  --output=mapdata/delhi-ncr.mbtiles
```

**Any OMT-compatible `.mbtiles` file works.** The bounding box above covers Delhi NCR. Adjust for any other city.

### Put the file in mapdata/

```
ncr-on-terminal/
└── mapdata/
    └── delhi-ncr.mbtiles   ← put it here
```

---

## Run

```sh
./ncr-on-terminal mapdata/delhi-ncr.mbtiles
```

---

## Controls

| Key | Action |
|---|---|
| `↑` `↓` `←` `→` | Pan the map |
| `w` `a` `s` `d` | Pan (WASD) |
| `k` `h` `j` `l` | Pan (vim-style) |
| `+` / `=` | Zoom in |
| `-` / `_` | Zoom out |
| Scroll wheel | Zoom in / out |
| `?` | Toggle help screen |
| `q` / `Ctrl+C` | Quit |

---

## Map symbols

| Symbol | Meaning |
|---|---|
| `M` (cyan) | Metro / subway station |
| `T` (blue) | Railway / train station |
| `+` (red) | Hospital |
| `f` (amber) | Restaurant / food |
| `g` (green) | Fuel station |

Coloured road names, place labels (towns, suburbs), and building outlines appear as you zoom in.

---

## Tips

- **AMOLED look:** set your terminal background to `#000000` (pure black). The app's colour palette is designed for dark backgrounds.
- **Best zoom range:** zoom 10–13 for an overview of Delhi; zoom 14–15 for street-level detail.
- Labels are in Latin script (`name:latin` property from OSM). Switch to Hindi in the style config if preferred.

---

## Running the SSH demo server

Let anyone try the map with a single `ssh` command — no install on their end.

**Requires:** `go get charm.land/wish/v2 && go mod tidy`

```sh
# Generate a stable SSH host key (do this once; keep it between restarts)
ssh-keygen -t ed25519 -f ssh_host_ed25519_key -N ""

# Run the server
go run ./cmd/sshserver \
  --addr :2222 \
  --host-key ssh_host_ed25519_key \
  --tiles mapdata/delhi-ncr.mbtiles
```

Or with Docker:

```sh
docker build -f Dockerfile.sshserver -t ncr-sshserver .

docker run -d -p 2222:2222 \
  -v $(pwd)/ssh_host_ed25519_key:/app/ssh_host_ed25519_key:ro \
  -v $(pwd)/mapdata:/app/mapdata:ro \
  ncr-sshserver
```

Point a DNS A record at your server and share:

```
ssh yourname.dev -p 2222
```

---

## How it works

```
MBTiles (SQLite) → MVT protobuf tiles → orb geometry
    ↓
DouglasPeucker simplification
    ↓
Tile → screen coordinate transform
    ↓
Braille buffer (each cell = 2×4 sub-pixel dots)
    ↓
ANSI-escaped Unicode output → BubbleTea View()
```

Key packages:

| Package | Role |
|---|---|
| [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) | TUI framework (Elm architecture) |
| [charmbracelet/wish](https://github.com/charmbracelet/wish) | SSH server for the demo |
| [paulmach/orb](https://github.com/paulmach/orb) | MVT decoding + geometry |
| [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) | Pure-Go SQLite for MBTiles |
| [charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss) | Terminal styling / HUD |

---

## Roadmap

- [ ] **GTFS metro route planning** — a fork of this repo will add the Delhi Metro route network using real GTFS data, letting you plan trips entirely in the terminal
- [ ] Better label collision avoidance (priority-sorted placement)
- [ ] Configurable starting location (flags)
- [ ] Other Indian cities (Bengaluru, Mumbai, Chennai bounding boxes)

---

## License

MIT — do whatever you want with it.

---

*Built by [Akash Parashar](https://akashparashar.dev) · [GitHub](https://github.com/adot-7/ncr-on-terminal)*
