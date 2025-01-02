package main

import (
	"errors"

	"github.com/charmbracelet/lipgloss"
)

var (
	ErrorStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0033"))
	InfoStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ffcc"))
	WarningStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffcc00"))
	ErrInvalidURL    = errors.New("Invalid url")
	ErrInvalidSchema = errors.New("Schema must be provided and must be either http or https")
)
