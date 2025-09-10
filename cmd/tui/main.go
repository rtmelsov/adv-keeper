package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/rtmelsov/adv-keeper/internal/tui"
)

var (
	version   = "dev"
	buildDate = "-"
)

func main() {
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()
	if *showVersion {
		fmt.Printf("adv-keeper %s (built %s)\n", version, buildDate)
		return
	}
	p := tea.NewProgram(tui.InitialModel())
	final, err := p.Run()
	if err != nil {
		fmt.Printf("Alas, there's been an error: %v\n", err)
		os.Exit(1)
	}

	// После выхода выведем список выбранных покупок
	m := final.(tui.TuiModel)
	if len(m.Selected) == 0 {
		fmt.Println("Ничего не выбрано.")
		return
	}

	// Соберём в стабильном порядке
	idxs := make([]int, 0, len(m.Selected))
	for i := range m.Selected {
		idxs = append(idxs, i)
	}
	sort.Ints(idxs)
}
