package helpers

import (
	"errors"
	"fmt"
	"golang.org/x/text/unicode/norm"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

func SafeBase(name string) string {
	base := filepath.Base(name) // срезать пути
	base = strings.TrimSpace(base)

	// Нормализуем (избавляемся от странных сочетаний)
	base = norm.NFKC.String(base)

	// Удаляем управляющие/опасные символы и оставляем буквы/цифры + .-_ пробел
	buf := make([]rune, 0, len(base))
	for _, r := range base {
		switch {
		case r == '.', r == '-', r == '_', r == ' ':
			buf = append(buf, r)
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			buf = append(buf, r)
			// остальное — выкидываем
		}
	}
	if len(buf) == 0 {
		return "file"
	}

	// Ограничим длину
	if len(buf) > 128 {
		buf = buf[:128]
	}

	out := string(buf)

	// Защита от точек/пробелов на конце (Windows) и зарезервированных имён
	out = strings.TrimRight(out, " .")
	reserved := map[string]struct{}{
		"CON": {}, "PRN": {}, "AUX": {}, "NUL": {},
		"COM1": {}, "COM2": {}, "COM3": {}, "COM4": {}, "COM5": {}, "COM6": {}, "COM7": {}, "COM8": {}, "COM9": {},
		"LPT1": {}, "LPT2": {}, "LPT3": {}, "LPT4": {}, "LPT5": {}, "LPT6": {}, "LPT7": {}, "LPT8": {}, "LPT9": {},
	}
	up := strings.ToUpper(strings.TrimSuffix(out, filepath.Ext(out)))
	if _, bad := reserved[up]; bad {
		out = "_" + out
	}

	return out
}

func NextAvailableName(dir, name string) string {
	name = SafeBase(name)
	if name == "" {
		name = "file"
	}

	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)

	candidate := filepath.Join(dir, name)
	for i := 0; ; i++ {
		_, err := os.Stat(candidate)
		if errors.Is(err, os.ErrNotExist) {
			return filepath.Base(candidate)
		}
		// "file (1).txt"
		try := fmt.Sprintf("%s (%d)%s", base, i+1, ext)
		candidate = filepath.Join(dir, try)
	}
}
