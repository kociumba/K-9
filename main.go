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
		SetString("K-9"). // tool branding
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("#6147FF")). // The tool color
		Foreground(lipgloss.Color("0"))
		// Border(lipgloss.Color("0"), 1)
	// Add a custom style for key `err`
	styles.Keys["err"] = lipgloss.NewStyle().Foreground(lipgloss.Color("204"))
	styles.Values["err"] = lipgloss.NewStyle().Bold(true)

	log.SetStyles(styles)

	log.Info("🫡  K-9 starting...")

	parseConfig()
}
