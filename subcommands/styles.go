package subcommands

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	greenStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("10"))
	
	errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("9"))
	
	yellowStyle2 = lipgloss.NewStyle().
		Foreground(lipgloss.Color("11"))
	
	infoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("12"))
	
	cyanStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("14"))
	
	cyanBold = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("14"))
	
	magentaBold = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("13"))
	
	passStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("10")).
		Bold(true)
	
	failStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("9")).
		Bold(true)
	
	testHeaderStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("14")).
		Bold(true).
		Underline(true)
	
	testInfoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("12"))
	
	templateTitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		Underline(true)
	
	configItemStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("14"))
	
	configValueStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("10"))
	
	configTitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		MarginBottom(1)
	
	ruleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("14"))
	
	boldStyle = lipgloss.NewStyle().
		Bold(true)
	
	titleDaily = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("14"))
	
	panelStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("14")).
		Padding(1, 2)
	
	tableHeaderStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("14"))
)