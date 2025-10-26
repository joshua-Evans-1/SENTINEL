package tui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/0xReLogic/SENTINEL/checker"
	"github.com/0xReLogic/SENTINEL/config"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/stopwatch"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Root struct {
	title		string
	height		int
	width		int
	KeyMap		KeyMap
	cfg  		*config.Config
	textInput 	textinput.Model
	table		TableModel
	stopwatch	stopwatch.Model
	status		[]checker.ServiceStatus
	currentView	string
	detailView 	DetailModel
}

type StatusMsg struct {
	status 		checker.ServiceStatus
}

type ChangeViewMsg struct {
	view 	string
}

func ( r Root ) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append( cmds, r.CheckServices() )
	cmds = append( cmds, r.stopwatch.Init() )	
	cmds = append( cmds, r.detailView.Init() )
	return tea.Batch( cmds... )
}

func ( r Root ) Update( msg tea.Msg ) ( tea.Model, tea.Cmd ) {
	var cmds []tea.Cmd
	switch msg := msg.( type ) {
		case tea.KeyMsg:
			switch {
				case key.Matches( msg, r.KeyMap.Quit ):
					if !r.textInput.Focused() {
						return r, tea.Quit
					}
				case key.Matches( msg, r.KeyMap.Back ):
					if r.textInput.Focused() {
						r.textInput.Blur()
						r.table.Focus()
						return r, nil
					}
					if r.currentView != "R00T" {
						r.table.Focus()
						return r, ChangeView( "R00T" )

					}
				case key.Matches( msg, r.KeyMap.FocusSearch ):
					r.textInput.Focus()
					r.table.Blur()
					return r, textinput.Blink
			}
		case tea.WindowSizeMsg:
			r = r.handleResize( msg.Height, msg.Width )
			return r, nil
		case StatusMsg:
			r = r.updateStatus( msg.status )
			r.table = r.table.UpdateRows( r.textInput.Value(), r.status )
			return r, nil
		case stopwatch.TickMsg:
			if r.stopwatch.Elapsed() >= time.Second * 5 {

				r.stopwatch = stopwatch.NewWithInterval( time.Second )
				var cmds []tea.Cmd
				cmds = append( cmds, r.CheckServices() )
				cmds = append( cmds, r.stopwatch.Init() )	
				return r, tea.Batch( cmds... )
			}
		case ChangeViewMsg:
			r.currentView = msg.view
			if msg.view != "R00T" {
				r.table.Blur()
				for _, stats := range r.status {
					if stats.Name == msg.view {
						r.detailView = NewDetailModel( stats )
					}
				}
			}
			return r, nil
	}
	var textInputCmd tea.Cmd
	r.textInput, textInputCmd = r.textInput.Update( msg )
	cmds = append( cmds, textInputCmd )

	r.table = r.table.UpdateRows( r.textInput.Value(), r.status )
	var tableCmd tea.Cmd
	r.table, tableCmd = r.table.Update( msg )
	cmds = append( cmds, tableCmd )

	var stopwatchCmd tea.Cmd
	r.stopwatch, stopwatchCmd = r.stopwatch.Update( msg )
	cmds = append( cmds, stopwatchCmd )
	
	var detailViewCmd tea.Cmd
	r.detailView, detailViewCmd = r.detailView.Update( msg )
	cmds = append( cmds, detailViewCmd )

	return r, tea.Batch( cmds... )
}

func ( r Root ) View() string {
	rootWindow := lipgloss.NewStyle().
		Height( r.height - 3 ).
		Width( r.width - 2 ).
		BorderStyle( lipgloss.RoundedBorder() ).
		BorderForeground( lipgloss.Color( "#ffffff" ) )
	var mainView string
	switch( r.currentView ) {
		case "R00T":
			mainView = lipgloss.JoinVertical( 
				lipgloss.Top,
				rootWindow.Render(  r.searchWindow(), r.table.View() ),
				r.statusbar(),
			)
			default:
				mainView = lipgloss.JoinVertical(
					lipgloss.Top,
					rootWindow.Render( r.detailView.View() ),
					r.statusbar(),
				)
		}	
	return mainView 
}

func ( r Root ) searchWindow( ) string {
	searchWindow := lipgloss.NewStyle().
		BorderLeft( false ).
		BorderRight( false ).
		BorderTop( false ).
		BorderBottom( true ).
		BorderStyle( lipgloss.ThickBorder() ).
		BorderForeground( lipgloss.Color( "#ffffff" ) ).
		Height( 1 ).
		Width( r.width - 2 ).Padding( 0, 2 )
	return searchWindow.Render( r.textInput.View() )	
}

func ( r Root ) statusbar(  ) string {
	upCount, downCount := r.UpCount()
	selectedService := ""
	if r.table.SelectedRow() != nil {
		selectedService = r.table.SelectedRow()[0]
	}

	itemStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color("#ffffff"))

	statusItems := []string{
		itemStyle.Render(" " + fmt.Sprintf( "%-8.8s", selectedService ) + " "),
		itemStyle.Render(" Services: " + strconv.Itoa(len(r.status)) + " "),
		itemStyle.Render(" Up: " + strconv.Itoa(upCount) + " "),
		itemStyle.Render(" Down: " + strconv.Itoa(downCount) + " "),
		itemStyle.Render(" Last Check: " + r.stopwatch.View() + " "),
		itemStyle.Render(" ? for help "),
	}

	// Calculate total content width
	totalContentWidth := 0
	for _, item := range statusItems {
		totalContentWidth += lipgloss.Width(item)
	}

	// Calculate available space for spacers
	availableSpace := r.width - totalContentWidth
	spacerCount := len(statusItems) - 1

	if spacerCount > 0 && availableSpace > 0 {
		spacerWidth := availableSpace / spacerCount
		extraSpace := availableSpace % spacerCount
		
		// Build the content with spacers
		var content strings.Builder
		for i, item := range statusItems {
			content.WriteString(item)
			if i < len(statusItems)-1 {
				spacerSize := spacerWidth
				if i < extraSpace {
					spacerSize++
				}
				spacer := strings.Repeat(" ", spacerSize)
				content.WriteString(spacer)
			}
		}
		return lipgloss.NewStyle().Width(r.width).Height(1).Render(content.String())
}

// Fallback if no space for spacers
statusBarContent := lipgloss.JoinHorizontal(lipgloss.Left, statusItems...)
return lipgloss.NewStyle().Width(r.width).Height(1).Render(statusBarContent)
}

func ( r Root ) handleResize( height, width int ) Root {
	r.height = height
	r.width = width
	r.textInput.Width = width - 8
	r.table = r.table.handleResize( r.height, r.width )
	r.table = r.table.CreateTable()
	r.table = r.table.UpdateRows( r.textInput.Value(), r.status )
	return r
}

func InitialRootModel( cfg *config.Config ) Root {
	r := Root {
		cfg: cfg,
		title: "SENTINEL Monitor",
		textInput: textinput.New(),
		KeyMap: DefaultKeyMap(),
		stopwatch: stopwatch.NewWithInterval( time.Second ),
	}

	r.textInput.Placeholder = "filter"
	r.textInput.CharLimit = 250
	r.textInput.Width = 40

	r.table = r.table.CreateTable()
	r.table = r.table.UpdateRows( r.textInput.Value(), r.status )
	r.table.SetCursor( 1 )
	r.currentView = "R00T"
	return r

}

func ( r Root ) updateStatus( services checker.ServiceStatus ) Root {
	var exists bool
	for i, service := range r.status {
		if service.Name == services.Name {
			r.status[i] = services
			exists = true
		}
	} 
	if !exists {
		r.status = append( r.status, services )
	}
	return r
}

func ( r Root ) CheckServices() tea.Cmd {
	var cmds []tea.Cmd
	for _, service := range r.cfg.Services {
		cmds = append( cmds, r.CheckService( service ) )
	}
	return tea.Batch( cmds... )
}

func ( r Root ) CheckService( service config.Service ) tea.Cmd {
	return func() tea.Msg {
		status := checker.CheckService( service.Name, service.URL )
		return StatusMsg{ status: status }
	}
}

func ( r Root ) UpCount(  ) ( int, int ) {
	var upCount, downCount int
	for _, service := range r.status {
		if service.IsUp {
			upCount++
		} else {
			downCount++
		}
	}
	return upCount, downCount
}

func ChangeView( view string ) tea.Cmd {
	return func() tea.Msg {
		return ChangeViewMsg{
			view: view,
		}
	}
}

