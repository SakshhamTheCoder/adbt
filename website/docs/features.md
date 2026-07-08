---
sidebar_position: 3
---

# Core Features

## Dashboard
- The main entry point for managing devices.

![Dashboard](/img/screenshots/dashboard_2.png)

## Device Management & Info
- **Device Info**: View detailed stats (Battery, Storage, Resolution, Android Version).
- **Power Controls**: Reboot, Recovery, Bootloader, and Screen Toggle.
- **Scrcpy Integration**: Launch screen mirroring with a single keypress.
- **Screen Capture**: Take a screenshot (`p`) or start/stop a screen recording (`v`); files are saved to your host's Downloads folder.

<!-- TODO(screenshot): refresh device_info_2.png to show the Screenshot and Screen record actions. -->
![Device Info](/img/screenshots/device_info_2.png)

## Performance Monitor
- **Real-time Stats**: CPU, Memory, and Network usage monitoring.
- **Visual Graphs**: Live progress bars for system resource consumption.
- **Battery**: Live level, status, temperature, voltage, health, and plug type.
- **Thermal**: Per-zone temperature readouts (where the device exposes them).

<!-- TODO(screenshot): refresh performance_2.png to show the Battery and Thermal sections. -->
![Performance Monitor](/img/screenshots/performance_2.png)

## App Manager
- **List & Search**: Browse all installed applications.
- **Filtering**: Toggle between User and System apps.
- **Details**: The highlighted app shows its version, APK size, and target SDK.
- **Multi-select**: Mark several apps with `Space` and uninstall them in one batch.
- **Actions**:
  - Launch App
  - Force Stop
  - Clear Data
  - Uninstall
  - Install APK (`i`)
  - Extract APK to host (`e`)

<!-- TODO(screenshot): refresh app_manager_2.png to show the details pane and selection markers. -->
![App Manager](/img/screenshots/app_manager_2.png)

## File Explorer
- **Browse**: Navigate the device file system seamlessly.
- **Transfer**: Pull files from the device (`p`) and push files from your host (`u`).
- **Manage**: Create new directories (`n`) and delete files or folders with confirmation (`d`).

<!-- TODO(screenshot): refresh file_explorer_2.png to show the push / new folder actions. -->
![File Explorer](/img/screenshots/file_explorer_2.png)

## Logcat Viewer
- **Live Streaming**: Real-time log capture with a pause/resume toggle.
- **Filtering**: Filter by log level, and narrow to a single app by package name or PID.
- **Search**: Text search with real-time highlighting.
- **Save**: Write the current buffer to a timestamped file on your host (`w`).

![Logcat Viewer](/img/screenshots/logcat_2.png)

## Input Sender
- **Send Text**: Type a line and press `Enter` to send it to the device.
- **Key Events**: Send Back, Home, and Recents, plus Enter/Backspace/Tab and arrow keys.
- **Live Capture**: Toggle a mode (`Tab`) that forwards each keystroke to the device as you type.

> 📷 _Screenshot pending — add `website/static/img/screenshots/input.png`, then replace this note with `![Input Sender](/img/screenshots/input.png)`._

## Intent Tester
- **Activity & Broadcast**: Send intents to test deep links and receiver behavior.
- **Suggestions**: Quick access to common Android intent actions.

![Intent Tester](/img/screenshots/intent_tester.png)

## Port Forwarding
- **Forward & Reverse**: Manage network connections between host and device.
- **Multiple Schemes**: Supports tcp, localabstract, localreserved, and more.

![Port Forwarding](/img/screenshots/port_manager.png)

