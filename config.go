package main

import (
	"errors"

	"github.com/charmbracelet/lipgloss"
)

var (
	ErrorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0033"))
	InfoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ffcc"))
	WarningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffcc00"))

	TableBorderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("99"))

	TableHeaderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("99")).
				Bold(true).
				Align(lipgloss.Center)

	TableEvenRowStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#cccccc"))

	TableOddRowStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#999999"))

	ErrInvalidURL    = errors.New("Invalid url")
	ErrInvalidSchema = errors.New("Schema must be provided and must be either http or https")
)
