# adbt: A Terminal UI for Android Debug Bridge

**adbt** is a modern, interactive terminal application built on the Android Debug Bridge (ADB). It provides a structured, keyboard-driven interface for managing Android devices directly within the terminal environment.

Developed using Go and the Charm Bracelet's Bubble Tea framework, `adbt` emphasizes clarity, performance, and correctness, offering a powerful alternative to raw command-line ADB interactions.

---

## Screenshots

|                    Screen View                     |
| :------------------------------------------------: |
|                **Dashboard Screen**                |
|   ![Dashboard Screen](screenshots/dashboard.png)   |
|               **Device Info Screen**               |
| ![Device Info Screen](screenshots/device_info.png) |

## 1\. Key Features

### Device Management and Control

-   **Device Discovery:** Automatically lists all connected Android devices and displays detailed information including Model, Serial, Android version, and current state.
-   **Automatic Selection:** Automatically selects the device when only a single device is connected.
-   **Destructive Actions:** Provides one-click access to critical device functions with mandatory confirmation prompts:
    -   Reboot device (`adb reboot`).
    -   Reboot to Recovery (`adb reboot recovery`).
    -   Reboot to Bootloader (`adb reboot bootloader`).
-   **Screen Mirroring:** Quick action to launch the external **scrcpy** tool for the selected device.

### Interactive Tools

-   **Live Logcat Viewer:** Streams and displays device log output in real-time within the TUI. Features include the ability to **clear** logs and **start/stop** the incoming stream.
-   **Dashboard:** A centralized screen providing an overview of the selected device's status and quick access to all major features via single-key shortcuts.
-   **ADB Shell:** Designated screen for interactive shell access (currently defined as a feature, with TUI implementation pending).

### Architecture and User Experience

-   **Keyboard-First Design:** The application is entirely keyboard-navigable, using standard key bindings:
    -   **Navigation:** `↑ / ↓` or `j / k`
    -   **Selection/Confirmation:** `enter`
    -   **Back/Cancel:** `esc`
    -   **Quick Actions:** Single-key shortcuts (e.g., `d`, `l`, `i`).
-   **Modal Input Handling:** Confirmation prompts and dialogs fully capture input, preventing key presses from affecting the underlying screen state, ensuring reliable execution of critical commands.
-   **Modular UI:** Built upon reusable Bubble Tea components (layouts, toasts, confirmation modals) for a consistent and robust visual experience.

---

## 2\. Technologies Used

| Category               | Technology                                               | Purpose                                                                     |
| :--------------------- | :------------------------------------------------------- | :-------------------------------------------------------------------------- |
| **Language**           | Go (Golang)                                              | Primary application development language (Go 1.21+).                        |
| **TUI Framework**      | [Bubble Tea](https://github.com/charmbracelet/bubbletea) | Implementation of the Terminal UI using The Elm Architecture.               |
| **Styling**            | [Lip Gloss](https://github.com/charmbracelet/lipgloss)   | Used for advanced terminal styling, layout, and visual component rendering. |
| **Tooling**            | Android SDK Platform Tools (`adb`)                       | The essential dependency for all device communication.                      |
| **Tooling (Optional)** | `scrcpy`                                                 | External utility used to enable the screen mirroring feature.               |

---

## 3\. Prerequisites and Installation

### Requirements

To run and build `adbt`, the following are mandatory:

1.  **Go:** Version 1.21 or newer.
2.  **Android SDK Platform Tools:** The `adb` executable must be installed and accessible in your system's `$PATH`.
3.  **An Android Device:** A physical device or emulator must be connected with USB debugging enabled.

### Detailed Installation Steps (Build from Source)

The following steps guide you through cloning the repository, building the executable, and making it globally accessible.

1.  **Clone the Repository:**

    ```bash
    git clone https://github.com/<your-username>/adbt
    cd adbt
    ```

2.  **Build the Executable:**
    The `go build` command compiles the source code, targeting the entry point in `cmd/adbt/main.go`.

    ```bash
    go build -o adbt ./cmd/adbt
    ```

    This creates an executable file named `adbt` in the project root.

3.  **Run Locally:**
    To execute the application from the project directory:

    ```bash
    ./adbt
    ```

4.  **(Optional) Install to PATH:**
    To use `adbt` globally from any terminal directory, move the executable to a directory included in your `$PATH` (e.g., `/usr/local/bin`):

    ```bash
    sudo mv adbt /usr/local/bin
    ```

    You can now launch the application using only the command name:

    ```bash
    adbt
    ```

---

## 4\. Usage and Key Commands

### Launching

Ensure an ADB-enabled device is connected before launching.

```bash
adbt
```

### Core Navigation and Global Bindings

| Key Binding     | Screen    | Action                                                       |
| :-------------- | :-------- | :----------------------------------------------------------- |
| `q` or `Ctrl+c` | Global    | Quits the application.                                       |
| `esc`           | Global    | Returns to the Dashboard or cancels an active dialog/prompt. |
| `enter` / `y`   | Prompt    | Confirms an action within a dialog.                          |
| `esc` / `n`     | Prompt    | Cancels an action within a dialog.                           |
| `d`             | Dashboard | Switches to the **Device Selection** screen.                 |

### Dashboard (Menu) Actions

The Dashboard acts as the primary hub for device interaction.

| Key | Menu Item   | Device Requirement | Function                                                             |
| :-- | :---------- | :----------------- | :------------------------------------------------------------------- |
| `d` | Devices     | No                 | View and select connected devices.                                   |
| `i` | Device Info | Yes                | View detailed properties and control device states (reboot, scrcpy). |
| `l` | Logcat      | Yes                | Access the live log viewer.                                          |
| `s` | Shell       | Yes                | Access the interactive ADB shell.                                    |

### Logcat Viewer Actions

| Key | Function   | Description                                     |
| :-- | :--------- | :---------------------------------------------- |
| `c` | Clear      | Clears the log history currently displayed.     |
| `s` | Start/Stop | Pauses or resumes the real-time logging stream. |

### Device Control (Device Info Screen) Actions

Actions on the Device Info screen trigger immediate confirmation prompts.

| Key | Action               | ADB Command Executed (if confirmed) |
| :-- | :------------------- | :---------------------------------- |
| `s` | Start Scrcpy         | `scrcpy -s <serial>`                |
| `r` | Reboot               | `adb -s <serial> reboot`            |
| `R` | Reboot to Recovery   | `adb -s <serial> reboot recovery`   |
| `b` | Reboot to Bootloader | `adb -s <serial> reboot bootloader` |

---

## 5\. Project Structure

The codebase is organized to separate concerns between the UI framework, global state, and the underlying ADB logic, following the Bubble Tea pattern.

```text
adbt/
├── cmd/
│   └── adbt/
│       └── main.go       # Application entry point; initializes the Bubble Tea program.
└── internal/
    ├── adb/              # Contains all functions that execute 'adb' commands.
    │   ├── client.go     # Low-level command execution and error handling (ExecuteCommand).
    │   ├── devices.go    # Logic for listing devices and fetching properties.
    │   ├── logcat.go     # Manages the streaming logcat process using channels and messages.
    │   └── device_info.go# Wrappers for device control commands (reboot, scrcpy).
    ├── state/            # Defines the global, mutable application data.
    │   └── app.go        # AppState: holds the selected device, device list, and screen dimensions.
    └── ui/               # All TUI components and screen logic.
        ├── app.go        # The root TUI model and router; manages screen switching.
        ├── components/   # Reusable UI primitives (views, styles, modals).
        │   ├── styles.go # Defines all terminal styling using Lip Gloss.
        │   └── confirm.go# Logic for the mandatory confirmation prompt modal.
        └── screens/      # Dedicated Bubble Tea models for each primary view.
            ├── dashboard.go# The main menu and status overview.
            ├── devices.go  # The device selection list.
            ├── logcat.go   # The logcat stream viewer.
            └── actions.go  # Helper functions to resolve actions into screen switching commands.
```

---

## 6\. Contributing

We welcome contributions, issue reports, and discussions. If you are interested in terminal UIs, Go development, Bubble Tea, or Android tooling, you are encouraged to get involved.

Please follow standard GitHub practices: fork the repository, create a feature branch, and submit a pull request for review.

---

## 7\. License

This project is released under the **MIT License**.

