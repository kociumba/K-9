package main

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	clog "github.com/charmbracelet/log"
)

var log = clog.New(os.Stderr)

func main() {

	styles := clog.DefaultStyles()
	styles.Levels[clog.InfoLevel] = lipgloss.NewStyle().
		SetString("INFO").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("#5fffd7")).
		Foreground(lipgloss.Color("0"))
	// Add a custom style for key `err`
	styles.Keys["err"] = lipgloss.NewStyle().Foreground(lipgloss.Color("204"))
	styles.Values["err"] = lipgloss.NewStyle().Bold(true)

	log.SetStyles(styles)

	log.Info("🫡  K-9 starting...")

	parseConfig()
}
