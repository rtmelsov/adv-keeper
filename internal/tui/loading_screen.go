package tui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/rtmelsov/adv-keeper/internal/ui"
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

func (m TuiModel) renderProgressStyled(title string, done, total int64, startedAt time.Time, errStr string, maxWidth int) string {
	if maxWidth < 24 {
		maxWidth = 24
	} // safety

	// percent
	var pct float64
	if total > 0 {
		pct = float64(done) / float64(total)
		if pct > 1 {
			pct = 1
		}
	}

	// bar width derived from maxWidth (leave room for " 100%")
	barW := max(10, min(60, maxWidth-12))
	filled := int(math.Round(pct * float64(barW)))
	if filled > barW {
		filled = barW
	}

	// colored bar
	bar := "[" +
		ui.StBarFill.Render(strings.Repeat("█", filled)) +
		ui.StBarEmpty.Render(strings.Repeat("░", barW-filled)) +
		"]"

	// speed/eta
	elapsed := time.Since(startedAt)
	var bps float64
	if elapsed > 0 {
		bps = float64(done) / elapsed.Seconds()
	}
	var eta time.Duration
	if bps > 0 && total > 0 {
		eta = time.Duration((float64(total-done) / bps)) * time.Second
	}

	// lines (styled)
	lineTitle := ui.StTitle.Render(ui.StValue.Render(m.Spin.View()) + " " + title)
	line1 := fmt.Sprintf("%s %s", bar, ui.StValue.Render(fmt.Sprintf("%3.0f%%", pct*100)))
	line2 := fmt.Sprintf("%s %s %s %s",
		ui.StLabel.Render("Done:"),
		ui.StValue.Render(humanizeBytes(done)),
		ui.StLabel.Render(" / Total:"),
		ui.StValue.Render(func() string {
			if total > 0 {
				return humanizeBytes(total)
			}
			return "??"
		}()),
	)
	line3 := fmt.Sprintf("%s %s   %s %s",
		ui.StLabel.Render("Speed:"),
		ui.StValue.Render(func() string {
			if bps > 0 {
				return humanizeBytes(int64(bps)) + "/s"
			}
			return "—"
		}()),
		ui.StLabel.Render("ETA:"),
		ui.StValue.Render(func() string {
			if eta > 0 {
				return humanizeDuration(eta)
			}
			return "—"
		}()),
	)

	errLine := ""
	if errStr != "" {
		errLine = "\n" + ui.StErr.Render("Ошибка: "+errStr)
	}

	return lipgloss.NewStyle().
		MaxWidth(maxWidth).
		Render(lineTitle + "\n\n" + line1 + "\n" + line2 + "\n" + line3 + errLine)
}

// your loader page uses the styled renderer
func (m TuiModel) loadingViewStyled(maxWidth int) string {
	switch {
	case m.Uploading && m.Downloading:
		return "ERROR: UPLOADING AND DOWNLOADING IN SAME TIME"
	case m.Uploading:
		return m.renderProgressStyled("Отправка файла… (ESC — отмена)", m.Uploaded, m.UploadTotal, m.UploadStart, m.Error, maxWidth)
	case m.Downloading:
		return m.renderProgressStyled("Скачивание файла… (ESC — отмена)", m.Downloaded, m.DownloadTotal, m.DownloadStart, m.Error, maxWidth)
	default:
		return "EMPTY SREAMING LOADING"
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
