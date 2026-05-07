---
sidebar_position: 5
---

# Troubleshooting & FAQs

If you run into issues while using `adbt`, here are some of the most common solutions.

### 1. Device shows as "Unauthorized"
When connecting a device over USB or Wi-Fi for the first time, Android requires you to authorize the host computer.
- **Solution**: Unlock your Android device's screen. A prompt will appear asking "Allow USB debugging?". Check "Always allow from this computer" and tap **Allow**. If `adbt` does not automatically update the status, unplug the device and plug it back in.

### 2. "adb command not found" or `adbt` fails to launch
`adbt` relies on having the Android Debug Bridge (`adb`) installed and accessible in your system's `PATH`.
- **macOS/Linux**: `brew install --cask android-platform-tools`
- **Windows**: Install ADB via Scoop or download the zip from the Android developer portal.

### 3. Screen mirroring (Scrcpy) doesn't start
`adbt` has native support to launch `scrcpy` with a single keypress. If it fails to launch, you likely do not have `scrcpy` installed on your host machine.
- **Solution**: Install scrcpy. Visit [scrcpy's GitHub repository](https://github.com/Genymobile/scrcpy) for installation instructions specific to your OS.

### 4. Wireless pairing (Wi-Fi) fails
Android 11+ supports wireless pairing via a pairing code and IP address.
- Ensure your host computer and the Android device are connected to the **exact same Wi-Fi network**.
- Verify you are entering the correct IP, Port, and Pairing PIN shown in the "Wireless debugging" menu in Developer Options.
- If it still fails, try restarting the ADB server manually with `adb kill-server` and `adb start-server` and try again in `adbt`.

### 5. Cannot pull or push files in File Explorer
Certain directories (like `/data/data/`) are restricted by android permissions and require **root access**. `adbt` runs with shell user privileges (`shell`), so you can only explore and transfer files from locations accessible to the shell user, such as `/sdcard/`.

If you have root access and want to browse root directories, you must restart `adbd` as root (via `adb root`) before launching `adbt`.
