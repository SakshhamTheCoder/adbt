package adb

import (
	"os"
	"path/filepath"
	"time"
)

// DefaultSaveDir is where host-side files (screenshots, recordings, pulled files,
// logcat dumps, extracted APKs) are written: ~/Downloads when it exists, otherwise
// the home directory, falling back to the working directory. Built with filepath so
// it is correct on Windows, macOS, and Linux.
func DefaultSaveDir() string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return "."
	}
	downloads := filepath.Join(home, "Downloads")
	if info, err := os.Stat(downloads); err == nil && info.IsDir() {
		return downloads
	}
	return home
}

// TimestampedName builds a filename like "screenshot-20060102-150405.png".
func TimestampedName(prefix, ext string) string {
	return prefix + "-" + time.Now().Format("20060102-150405") + ext
}
