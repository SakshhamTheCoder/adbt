# adbt — Android Debug Bridge TUI

**adbt** is a modern, keyboard-driven **Terminal User Interface (TUI)** for interacting with Android devices using the Android Debug Bridge (ADB).
It provides a structured and interactive alternative to raw `adb` commands while remaining fast, predictable, and safe for daily use.

Built in **Go** using **Charm’s Bubble Tea framework**, `adbt` focuses on clarity, correctness, and real-time device interaction entirely from the terminal.

---

## Screenshots

Screenshots reflect the current codebase and actual UI.

| Dashboard                               | Device Info                                 |
| --------------------------------------- | ------------------------------------------- |
| ![Dashboard](screenshots/dashboard.png) | ![Device Info](screenshots/device_info.png) |

---

## Features

### Device Management

-   Automatic detection of connected Android devices
-   Device selection with connection state and Android version
-   Automatic selection when exactly one device is connected
-   Wireless ADB pairing using IP address, port, and PIN

### Dashboard

-   Central hub displaying current device status
-   Keyboard shortcuts for all major features
-   Device-aware actions disabled when no device is selected

### Device Info and Controls

-   Displays comprehensive device details:
    -   Model, Serial, Android Version
    -   Battery Level and Status
    -   Storage Usage
    -   Display Resolution and Density
    -   IP Address
-   Supported actions:

    -   Start scrcpy
    -   Toggle Wi-Fi
    -   Toggle screen power
    -   Reboot device
    -   Reboot to recovery
    -   Reboot to bootloader

-   Mandatory confirmation prompts for destructive actions

### Logcat Viewer

-   Real-time streaming of `adb logcat`
-   Color-coded logs by priority (Verify, Debug, Info, Warning, Error, Fatal)
-   Filter by priority (press `f` to cycle levels)
-   Text search and highlighting (press `/`)
-   Pause and resume logging
-   Clear log buffer
-   Retains the last 1000 log lines for performance

### App Manager

-   Lists installed applications on the device
-   Filter by app type (User / System / All) using `f`
-   Search by package name using `/`
-   Launch, Force Stop, Uninstall, Clear Data actions
-   Reloads the application list on demand

### File Browser

-   Browse the device file system starting at `/sdcard`
-   Navigate directories
-   Pull files to local computer (Downloads folder)
-   Delete files and directories with confirmation
-   Refresh directory contents
-   Scrollable viewport for large directories

---

## Keyboard Navigation

### Global Keys

| Key           | Action               |
| ------------- | -------------------- |
| `q`, `Ctrl+C` | Quit the application |
| `esc`         | Go back or cancel    |
| `↑ ↓`, `j k`  | Navigate lists       |
| `enter`       | Select or confirm    |

### Dashboard Shortcuts

| Key | Action      |
| --- | ----------- |
| `d` | Devices     |
| `i` | Device Info |
| `l` | Logcat      |
| `a` | Apps        |
| `f` | Files       |

### Devices Shortcuts

| Key | Action          |
| --- | --------------- |
| `w` | Wireless Pair   |
| `r` | Refresh List    |
| `enter` | Select Device   |

### App Manager Shortcuts

| Key | Action          |
| --- | --------------- |
| `/` | Search apps     |
| `f` | Filter (User/System) |
| `enter` | Launch app      |
| `s` | Force Stop      |
| `u` | Uninstall       |
| `x` | Clear Data      |
| `r` | Refresh list    |

### File Browser Shortcuts

| Key | Action          |
| --- | --------------- |
| `enter` | Open / Enter    |
| `backspace` | Go Up Directory |
| `p` | Pull File       |
| `d` | Delete          |
| `r` | Refresh         |



### Logcat Shortcuts

| Key | Action          |
| --- | --------------- |
| `/` | Search logs     |
| `f` | Filter priority |
| `c` | Clear buffer    |
| `s` | Start / Stop    |

---

## Project Structure

The project follows Bubble Tea’s Elm-style architecture, separating UI, state, and ADB logic.

### Entry Point

-   `cmd/adbt/main.go`
    Initializes the Bubble Tea program and starts the UI loop.

### ADB Layer (`internal/adb`)

All direct interaction with the `adb` binary is handled here.

-   `client.go` — Command execution helpers and output parsing
-   `devices.go` — Device discovery and wireless pairing
-   `device_info.go` — Device control actions (reboot, Wi-Fi, screen, scrcpy)
-   `logcat.go` — Streaming logcat implementation
-   `files.go` — File system operations
-   `app_manager.go` — Installed applications listing

### Global State (`internal/state`)

-   `app.go` — Holds selected device, device list, and terminal dimensions

### UI Layer (`internal/ui`)

Responsible for all TUI rendering and interaction.

-   `app.go` — Root Bubble Tea model and screen router

#### Screens (`internal/ui/screens`)

Each screen is an isolated Bubble Tea model.

-   `dashboard.go` — Main dashboard and quick actions
-   `devices.go` — Device selection and wireless pairing
-   `device_info.go` — Device details and controls
-   `logcat.go` — Live logcat viewer
-   `files.go` — File browser
-   `app_manager.go` — App manager
-   `actions.go` — Action-to-screen resolution
-   `messages.go` — Screen-switching messages

#### UI Components (`internal/ui/components`)

Reusable UI building blocks.

-   `layout.go` — Layout and viewport handling
-   `header.go` — Header rendering
-   `list.go` — Device list rendering
-   `file_list.go` — File list rendering
-   `keyvalue.go` — Key-value table rendering
-   `confirm.go` — Confirmation modal
-   `form.go` — Input form modal
-   `toast.go` — Toast notifications
-   `styles.go` — Lip Gloss styling definitions

### CI / Release

-   `.github/workflows/release.yml` — Multi-OS build and GitHub release workflow

---

## Requirements

-   Go 1.21 or newer
-   Android SDK Platform Tools (`adb` available in PATH)
-   Android device or emulator with USB debugging enabled
-   Optional: `scrcpy` for screen mirroring

---

## Installation

### Download Prebuilt Binaries

Prebuilt binaries for Linux, macOS, and Windows are available on the GitHub Releases page.

```bash
chmod +x adbt
./adbt
```

### Build from Source

```bash
git clone https://github.com/SakshhamTheCoder/adbt
cd adbt
go build -o adbt ./cmd/adbt
./adbt
```

Optional global install:

```bash
sudo mv adbt /usr/local/bin
```

---

## Contributing

Contributions are welcome.
Fork the repository, create a feature branch, and submit a pull request for review.

---

## License

MIT License

