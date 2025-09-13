//go:build linux

package platformdirs

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func DownloadsDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// 1) Явная переменная окружения
	if dir := os.Getenv("XDG_DOWNLOAD_DIR"); dir != "" {
		dir = strings.ReplaceAll(dir, "$HOME", home)
		dir = strings.ReplaceAll(dir, "${HOME}", home)
		return filepath.Clean(dir), nil
	}

	// 2) ~/.config/user-dirs.dirs
	ud := filepath.Join(home, ".config", "user-dirs.dirs")
	if b, err := os.ReadFile(ud); err == nil {
		re := regexp.MustCompile(`(?m)^XDG_DOWNLOAD_DIR="?([^"\n]+)"?`)
		if m := re.FindStringSubmatch(string(b)); len(m) == 2 {
			dir := strings.ReplaceAll(m[1], "$HOME", home)
			dir = strings.ReplaceAll(dir, "${HOME}", home)
			return filepath.Clean(dir), nil
		}
	}

	// 3) Фоллбэк
	return filepath.Join(home, "Downloads"), nil
}
