package tui

import (
	"github.com/0xReLogic/SENTINEL/checker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DetailModel struct {
	status		checker.ServiceStatus
	width		int 
	height		int
	focused		bool
}

func ( m DetailModel ) Init() tea.Cmd {
	return nil	
}

func ( m DetailModel ) Update( msg tea.Msg ) ( DetailModel, tea.Cmd ) {
	return m, nil
}

func ( m DetailModel ) View(  ) string {
	statusStyle := lipgloss.NewStyle().Height( m.height - 7 )

	return statusStyle.Render( m.status.String()  )
}

func NewDetailModel( status checker.ServiceStatus ) DetailModel {
	details := DetailModel{
		status: status,
	}	
	return details
}

func ( m DetailModel ) HandleResize( height, width int ) DetailModel {
	m.height = height
	m.width = width 
	return m
}

func ( m *DetailModel ) Focus() {
	m.focused = true
}

func ( m *DetailModel ) Blur() {
	m.focused = false
}

func ( m DetailModel ) Focused() bool {
	return m.focused
}

