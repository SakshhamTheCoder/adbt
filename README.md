<p align="center">
  <img src="website/static/img/logoadbt.png" alt="adbt Logo" width="120" />
</p>

<h1 align="center">adbt</h1>

<p align="center">
  <b>A modern, keyboard-driven Terminal User Interface (TUI) for Android Debug Bridge.</b>
</p>

<p align="center">
  <a href="https://adbt-tui.vercel.app"><b>Documentation Website »</b></a>
</p>

<p align="center">
  <img src="https://img.shields.io/github/license/SakshhamTheCoder/adbt?style=flat-square&color=blue" alt="License" />
  <img src="https://img.shields.io/github/v/release/SakshhamTheCoder/adbt?style=flat-square&color=purple" alt="Latest Release" />
  <img src="https://img.shields.io/badge/built%20with-Bubble%20Tea-brightgreen?style=flat-square" alt="Built with Bubble Tea" />
</p>

---

**adbt** is a fast, keyboard-driven terminal user interface for interacting with Android devices via ADB. Built in Go using Charm’s **Bubble Tea** framework, it provides a structured and interactive command-line workspace that helps you avoid complex raw shell commands and context-switching.

For detailed guides, shortcuts, and troubleshooting, visit the [adbt Documentation Website](https://adbt-tui.vercel.app).

---

## ⚡ Core Features

- **Device Management**: Real-time connected device detection, wireless pairing (QR code logic/PIN), detailed device info, and power/reboot controls.
- **Scrcpy Integration**: Launch high-performance screen mirroring with a single keypress.
- **Performance Monitor**: Real-time visual graphs for device CPU, memory, and network usage.
- **App Manager**: Browse, search, and filter user/system apps, with quick controls to launch, force-stop, clear data, or uninstall.
- **File Explorer**: Seamlessly browse the device filesystem, pull files to your host machine, and manage directories.
- **Logcat Viewer**: Live stream device logs with severity filters (Debug, Info, Error, Fatal) and text search highlighting.
- **Intent Tester**: Construct and send custom activity/broadcast intents to test deep links and receiver behavior.
- **Port Forwarding**: Easily configure forward and reverse network connections between your host and device.

<details>
<summary><b>📷 View Screenshots Gallery</b></summary>
<br />

| Dashboard | Device Info |
| :---: | :---: |
| ![Dashboard](website/static/img/screenshots/dashboard_2.png) | ![Device Info](website/static/img/screenshots/device_info_2.png) |

| App Manager | File Explorer |
| :---: | :---: |
| ![App Manager](website/static/img/screenshots/app_manager_2.png) | ![File Explorer](website/static/img/screenshots/file_explorer_2.png) |

| Logcat Viewer | Performance Monitor |
| :---: | :---: |
| ![Logcat](website/static/img/screenshots/logcat_2.png) | ![Performance](website/static/img/screenshots/performance_2.png) |

| Intent Tester | Port Forwarding |
| :---: | :---: |
| ![Intent Tester](website/static/img/screenshots/intent_tester.png) | ![Port Forwarding](website/static/img/screenshots/port_manager.png) |

</details>

---

## 📦 Quick Installation

### Homebrew (macOS / Linux)
```bash
brew install --cask SakshhamTheCoder/tap/adbt
# or
brew tap SakshhamTheCoder/tap && brew install --cask adbt
```

### Scoop (Windows)
```bash
scoop bucket add SakshhamTheCoder https://github.com/SakshhamTheCoder/scoop-bucket
scoop install adbt
```

### AUR (Arch Linux)
```bash
yay -S adbt-bin
# or
paru -S adbt-bin
```

For manual binary downloads (`.deb`, `.rpm`, or generic archives), check the [Releases Page](https://github.com/SakshhamTheCoder/adbt/releases) or the [Installation Docs](https://adbt-tui.vercel.app/docs/installation).

---

## ⌨️ Essential Navigation

Once installed, plug in your device and run:
```bash
adbt
```

| Key | Action |
| :--- | :--- |
| `q` / `Ctrl+C` | Quit `adbt` |
| `Esc` / `Backspace` | Back / Go Up a directory |
| `↑` `↓` / `k` `j` | Navigate lists and menus |
| `Enter` | Select / Confirm / Launch |

For the full set of module-specific hotkeys (App Manager, File Explorer, Logcat, etc.), refer to the [Keyboard Shortcuts Cheat Sheet](https://adbt-tui.vercel.app/docs/shortcuts).

---

## 🛠️ Troubleshooting & Support

If you encounter connection issues, unauthorized device screens, or missing ADB paths:
- Consult the [Troubleshooting & FAQs Guide](https://adbt-tui.vercel.app/docs/troubleshooting).
- Open an issue on our [GitHub Issue Tracker](https://github.com/SakshhamTheCoder/adbt/issues).

## 📄 License

This project is licensed under the [MIT License](LICENSE).
