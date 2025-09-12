package tui

import (
	"fmt"
	"math"
	"strings"
	"time"
)

func humanizeBytes(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	div, exp := int64(unit), 0
	for x := n / unit; x >= unit; x /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(n)/float64(div), "KMGTPE"[exp])
}

func humanizeDuration(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	// форматы: 0:05, 1:23, 12:34:56
	sec := int(d.Seconds())
	h := sec / 3600
	m := (sec % 3600) / 60
	s := sec % 60
	if h > 0 {
		return fmt.Sprintf("%d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%d:%02d", m, s)
}

// Универсальный рендерер прогресса
func (m TuiModel) renderProgress(title string, done, total int64, startedAt time.Time, errStr string) string {
	var pct float64
	if total > 0 {
		pct = float64(done) / float64(total)
		if pct > 1 {
			pct = 1
		}
	}

	// Прогресс-бар
	const barW = 40
	filled := int(math.Round(pct * float64(barW)))
	if filled > barW {
		filled = barW
	}
	bar := "[" + strings.Repeat("█", filled) + strings.Repeat("░", barW-filled) + "]"

	// Скорость и ETA
	elapsed := time.Since(startedAt)
	var bps float64
	if elapsed > 0 {
		bps = float64(done) / elapsed.Seconds()
	}
	var eta time.Duration
	if bps > 0 && total > 0 {
		remaining := float64(total - done)
		eta = time.Duration(remaining/bps) * time.Second
	}

	line1 := fmt.Sprintf("%s %3.0f%%", bar, pct*100)
	line2 := fmt.Sprintf("%s / %s",
		humanizeBytes(done),
		func() string {
			if total > 0 {
				return humanizeBytes(total)
			}
			return "??"
		}(),
	)
	line3 := fmt.Sprintf("Скорость: %s/s   ETA: %s",
		func() string {
			if bps > 0 {
				return humanizeBytes(int64(bps))
			}
			return "—"
		}(),
		func() string {
			if eta > 0 {
				return humanizeDuration(eta)
			}
			return "—"
		}(),
	)

	errLine := ""
	if errStr != "" {
		errLine = "\nОшибка: " + errStr
	}

	return fmt.Sprintf("%s %s\n\n%s\n%s\n%s%s\n",
		m.Spin.View(), title, line1, line2, line3, errLine)
}

// Отдельные вьюхи
func (m TuiModel) uploadLoadingView() string {
	return m.renderProgress("Отправка файла… (ESC — отмена)", m.Uploaded, m.UploadTotal, m.UploadStart, m.Error)
}

func (m TuiModel) downloadLoadingView() string {
	return m.renderProgress("Скачивание файла… (ESC — отмена)", m.Downloaded, m.DownloadTotal, m.DownloadStart, m.Error)
}

// Общий вход в загрузочный экран
func (m TuiModel) loadingView() string {
	switch {
	case m.Uploading && m.Downloading:
		return m.uploadLoadingView() + "\n" + m.downloadLoadingView()
	case m.Uploading:
		return m.uploadLoadingView()
	case m.Downloading:
		return m.downloadLoadingView()
	default:
		return "" // или основной View()
	}
}
