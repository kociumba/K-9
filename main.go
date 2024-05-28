package main

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	clog "github.com/charmbracelet/log"
)

var k9p = clog.New(os.Stderr)

func main() {

	styles := clog.DefaultStyles()
	styles.Levels[clog.InfoLevel] = lipgloss.NewStyle().
		SetString("K-9"). // tool branding
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("#6147FF")). // The tool color
		Foreground(lipgloss.Color("0"))
		// Border(lipgloss.Color("0"), 1)
	// Add a custom style for key `err`
	styles.Levels[clog.ErrorLevel] = lipgloss.NewStyle().
		SetString("K-9").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("#ed8796")).
		Foreground(lipgloss.Color("0"))
	styles.Keys["err"] = lipgloss.NewStyle().Foreground(lipgloss.Color("204"))
	styles.Values["err"] = lipgloss.NewStyle().Bold(true)

	k9p.SetStyles(styles)

	k9p.Info("🫡  K-9 starting...")

	parseConfig()
}
