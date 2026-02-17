# wfinfo-go

`wfinfo-go` is a lightweight, efficient replacement for [WFInfo](https://wfinfo.warframestat.us/) specifically designed for Linux users.
It automates the process of identifying Warframe relic rewards and fetching their current platinum values from [warframe.market](https://warframe.market/).

The tool is written entirely in Go and utilizes Tesseract for high-accuracy Optical Character Recognition (OCR).

## Features

- **Automatic Detection:** Monitors Warframe's `EE.log` in real-time to detect when a relic reward screen appears.
- **Smart Screen Capture:** Uses X11 (via `xgb`) to capture only the Warframe window, ensuring privacy and efficiency.
- **Robust OCR:** Employs a specialized image preprocessing pipeline to isolate and binarize text before processing with Tesseract.
- **Fuzzy Matching:** Implements the Smith-Waterman algorithm for local alignment, providing high resilience against OCR errors in item names.
- **Live Market Data:** Fetches up-to-date pricing information directly from the `warframe.market` API.
- **Resource Efficient:** Uses an event-driven architecture for log monitoring and optimizes OCR/API requests to minimize CPU and IO overhead.

## Dependencies

### Runtime Dependencies
If you are using a pre-built binary, you only need the Tesseract engine and its shared libraries.

- **Tesseract OCR**: The engine and English language data.
- **Shared Libraries**: `libtesseract` and `libleptonica`.

### Build Dependencies
If you are building from source, you need the following in addition to the runtime dependencies:

- **Go** (version 1.26 or later)
- **C Compiler**: `gcc` or `clang` for CGO bindings.
- **Development Headers**: Headers for `libtesseract` and `libleptonica`.
- **pkg-config**: Required for the Go build system to locate Tesseract.

## Installation

### Ubuntu/Debian

**To run:**
```bash
sudo apt-get install tesseract-ocr libtesseract5
```

**To build:**
```bash
sudo apt-get install golang tesseract-ocr libtesseract-dev pkg-config build-essential
```

### Fedora

**To run:**
```bash
sudo dnf install tesseract tesseract-langpack-eng
```

**To build:**
```bash
sudo dnf install golang tesseract-devel gcc pkg-config
```

### Arch

**To run:**
```bash
sudo pacman -S tesseract tesseract-data-eng
```

**To build:**
```bash
sudo pacman -S go tesseract gcc pkgconf
```

## Building

You can build the project using the provided `Taskfile.yml` (requires [go-task](https://taskfile.dev/)) or standard Go commands.

### Using Task

```bash
task build
```
This will place the binary in `bin/wfinfo-go`.

### Using Go directly

```bash
go build -o bin/wfinfo-go ./cmd/wfinfo-go
```

## Usage

Run the program from your terminal. It will stay active and watch for reward screens.

```bash
./bin/wfinfo-go [flags]
```

### Flags

- `-h`: Shows help information.
- `-d [PATH]`: Path to your Steam Library where Warframe is installed (defaults to `~/.local/share/Steam`).
- `-f [PATH]`: Direct path to `EE.log`. This flag takes precedence over `-d`.

### Example

```bash
./bin/wfinfo-go -d ~/.local/share/Steam
```

## How It Works

1. **Log Watching:** The application uses `fsnotify` to monitor `EE.log`. It listens for specific markers indicating the reward screen has initialized (e.g., `VoidProjections: OpenVoidProjectionRewardScreenRMI`).
2. **Window Capture:** Upon detection, it finds the Warframe window via X11 properties and captures its contents.
3. **Preprocessing:** The captured image is processed to identify text regions, isolate them based on color, and binarize the output to maximize OCR accuracy.
4. **OCR & Matching:** Tesseract extracts text from the processed image. The resulting strings are compared against a local cache of Warframe items using the Smith-Waterman algorithm to find the most likely matches.
5. **Market Integration:** For each identified item, the program queries `warframe.market` for current sell orders and prints the results to your terminal. Item data and market versions are cached locally in `~/.cache/wfm-go/` to reduce API load and improve startup time.

## Architecture

The program is designed with a focus on performance:
- **Parallel Processing:** Interleaves I/O-bound market API requests with CPU-bound OCR operations.
- **Event-Driven:** Avoids polling the filesystem, reducing idle resource usage.
- **Concurrency:** Uses Go's concurrency primitives (channels and goroutines) to handle detection and processing asynchronously.

## License

This project is licensed under the BSD 3-Clause License. See the [LICENSE](LICENSE) file for details.
