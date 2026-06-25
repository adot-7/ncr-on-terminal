# NCR on Terminal

> Delhi NCR, rendered in Braille characters. Entirely in your terminal.

<!-- Replace this line with the GIF once you've recorded it with demo.tape -->
![NCR on Terminal demo](demo.gif)

A terminal map viewer for the Delhi National Capital Region built with [BubbleTea](https://github.com/charmbracelet/bubbletea). Roads, buildings, metro stations, food joints, water, forests, and place labels - all drawn using Unicode Braille characters for sub-character pixel resolution.

---

## Try it instantly (no install required)

```sh
ssh ncr.akashparashar.dev
#Are you sure you want to continue connecting (yes/no/[fingerprint])? yes
```

> **Note:** The SSH demo server may be slow due to many concurrent users and may not always be running. See [Local install](#install) to run it yourself. 

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

Extract and [Run](#run). The binary has no dependencies.

**Or install with Go:**

```sh
go install github.com/adot-7/ncr-on-terminal@latest

```
Requires Go 1.21+. Binary lands in ~/go/bin/. 

**Or build from source:**

```sh
git clone https://github.com/adot-7/ncr-on-terminal
cd ncr-on-terminal
go build -o ncr-on-terminal .
```

Requires Go 1.21+.

---

## Setup: getting an MBTiles file

The app needs an [OpenMapTiles](https://openmaptiles.org)-compatible `.mbtiles` vector tile file.

### Option A — Download the ready-made Delhi NCR file (recommended)

The Delhi NCR tile file is distributed directly from the [Releases page](https://github.com/adot-7/ncr-on-terminal/releases/latest) as a release asset:

```sh
# Download (~50 MB)
wget https://github.com/adot-7/ncr-on-terminal/releases/download/v0.1.0/delhi-ncr.mbtiles

# Put it in mapdata/
mkdir -p mapdata
mv delhi-ncr.mbtiles mapdata/
```

> **Data attribution:** Tile data sourced from [BBBike.org extracts](https://extract.bbbike.org/) — free city extracts from OpenStreetMap.
> Map data © [OpenStreetMap contributors](https://www.openstreetmap.org/copyright), licensed [ODbL](https://opendatacommons.org/licenses/odbl/).

### Option B — Generate tiles for any city or region with tilemaker

If you want tiles for a different city (or want to generate Delhi tiles yourself):

**Step 1 — Get a OSM extract from BBBike**

Go to [extract.bbbike.org](https://extract.bbbike.org/), draw your area on the map, select **PBF format**, and request the extract. You'll receive a download link by email (usually within a few minutes for city-sized areas).

Alternatively, pre-made city extracts are available at [download.bbbike.org/osm/bbbike/](https://download.bbbike.org/osm/bbbike/):

```sh
# Example: Delhi
wget https://download.bbbike.org/osm/bbbike/Delhi/Delhi.osm.pbf
```

**Step 2 — Install tilemaker**

```sh
# macOS
brew install tilemaker

# Linux — download the latest release binary
wget https://github.com/systemed/tilemaker/releases/latest/download/tilemaker-ubuntu-22.04.zip
unzip tilemaker-ubuntu-22.04.zip
```

Or [build from source](https://github.com/systemed/tilemaker#building).

**Step 3 — Clone the OpenMapTiles config files**

```sh
git clone https://github.com/systemed/tilemaker
```

The repo includes `resources/config-openmaptiles.json` and `resources/process-openmaptiles.lua` — the config files that produce tiles in the schema this app expects.

**Step 4 — Convert to MBTiles**

```sh
cd tilemaker
./tilemaker \
  --input /path/to/Delhi.osm.pbf \
  --output ../mapdata/delhi-ncr.mbtiles \
  --config resources/config-openmaptiles.json \
  --process resources/process-openmaptiles.lua
```

Conversion takes 2–5 minutes for a city-sized area. Output is ~40–60 MB.

> **Any OpenMapTiles-compatible `.mbtiles` file works.** Adjust the bounding box and filename for your city.

### Put the file in mapdata/

```
ncr-on-terminal/
└── mapdata/
    └── delhi-ncr.mbtiles   ← either downloaded or generated
```

---

## Run

```sh
./ncr-on-terminal mapdata/delhi-ncr.mbtiles
```

---

## Controls

| Key             | Action                             |
| --------------- | ---------------------------------- |
| `↑` `↓` `←` `→` | Pan the map                        |
| `w` `a` `s` `d` | Pan (WASD)                         |
| `k` `h` `j` `l` | Pan (vim-style)                    |
| `+` / `=`       | Zoom in                            |
| `-` / `_`       | Zoom out                           |
| Scroll wheel    | Zoom in / out (except in ssh mode) |
| `?`             | Toggle help screen                 |
| `q` / `Ctrl+C`  | Quit                               |

---

## Map symbols

| Symbol      | Meaning                 |
| ----------- | ----------------------- |
| `M` (cyan)  | Metro / subway station  |
| `T` (blue)  | Railway / train station |
| `+` (red)   | Hospital                |
| `🍴`, etc   | Restaurant / food       |
| `g` (green) | Fuel station            |

Coloured road names, place labels (towns, suburbs), and building outlines appear as you zoom in.

---

## Tips

- **AMOLED look:** set your terminal background to `#000000` (pure black). The app's colour palette is designed for dark backgrounds.
- **Best zoom range:** zoom 10–13 for an overview of Delhi; zoom 14–15 for street-level detail.
- Labels are in Latin script (`name:latin` from OSM). Switch to Hindi by editing `featureName()` in `render/renderer.go`.

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

## mapscii

[mapscii](https://github.com/rastapasta/mapscii) is great but fetches tiles 
from a remote server on every pan. It felt slow, network-dependent, and breaks offline.

NCR on Terminal reads from a local `.mbtiles` file (SQLite on disk). 
Every tile read is a microsecond filesystem lookup. No network, no tile server, 
works on a plane (idk why you would need to see a map in braille on terminal on plane when you can just look down)

## Roadmap

- [ ] **metro route planning** — kinda fun, will be released in the coming week
- [ ] Better label collision avoidance (priority-sorted placement)
- [ ] Configurable starting location via flags

---

## License

MIT — do whatever you want with it.

Map data © [OpenStreetMap contributors](https://www.openstreetmap.org/copyright) (ODbL).
Tile source: [BBBike.org extracts](https://extract.bbbike.org/).

---

*Built by [Akash Parashar](https://akashparashar.dev) · [GitHub](https://github.com/adot-7/ncr-on-terminal)*
