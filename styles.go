package main

import "github.com/charmbracelet/lipgloss"

var (
	colorAmber   = lipgloss.Color("#E0AF68")
	colorBlue    = lipgloss.Color("#7AA2F7")
	colorGreen   = lipgloss.Color("#9ECE6A")
	colorMuted   = lipgloss.Color("#565f89")
	colorFg      = lipgloss.Color("#C0CAF5")

	appStyle = lipgloss.NewStyle().
			Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(colorAmber).
			Bold(true).
			Padding(0, 1).
			MarginBottom(1)

	detailTitleStyle = lipgloss.NewStyle().
				Foreground(colorAmber).
				Bold(true).
				Padding(0, 1).
				Border(lipgloss.NormalBorder(), false, false, true, false).
				BorderForeground(colorMuted)

	metaStyle = lipgloss.NewStyle().
			Foreground(colorMuted).
			Italic(true).
			Padding(0, 1).
			MarginBottom(1)

	sectionTitleStyle = lipgloss.NewStyle().
				Foreground(colorBlue).
				Bold(true).
				Padding(0, 1)

	contentStyle = lipgloss.NewStyle().
			Foreground(colorFg).
			Padding(0, 2)

	viewportStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorMuted).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(colorMuted).
			Padding(1, 1, 0, 1)

	statusStyle = lipgloss.NewStyle().
			Foreground(colorGreen)
)
