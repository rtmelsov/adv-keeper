//go:build windows

package platformdirs

import (
	"os"
	"path/filepath"
	"unsafe"

	"golang.org/x/sys/windows"
)

func DownloadsDir() (string, error) {
	// FOLDERID_Downloads = 374DE290-123F-4565-9164-39C4925E467B
	guid := windows.KNOWNFOLDERID{
		Data1: 0x374de290, Data2: 0x123f, Data3: 0x4565,
		Data4: [8]byte{0x91, 0x64, 0x39, 0xc4, 0x92, 0x5e, 0x46, 0x7b},
	}
	var p *uint16
	if err := windows.SHGetKnownFolderPath(&guid, 0, 0, &p); err == nil && p != nil {
		defer windows.CoTaskMemFree(unsafe.Pointer(p))
		return windows.UTF16PtrToString(p), nil
	}
	// Фоллбэк
	profile := os.Getenv("USERPROFILE")
	if profile == "" {
		if home, err := os.UserHomeDir(); err == nil {
			profile = home
		}
	}
	return filepath.Join(profile, "Downloads"), nil
}
