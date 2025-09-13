//go:build darwin

package platformdirs

import (
	"os"
	"path/filepath"
)

func DownloadsDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	// В macOS каталог называется "Downloads" (локализуется только отображение)
	return filepath.Join(home, "Downloads"), nil
}
