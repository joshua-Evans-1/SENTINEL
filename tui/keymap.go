package tui

import ( 

	"github.com/charmbracelet/bubbles/key"

)

type KeyMap struct {
	Quit 		key.Binding
	Up 			key.Binding
	Down		key.Binding
	ViewService	key.Binding
	Refresh		key.Binding
	Back		key.Binding
	ShowHelp	key.Binding
	FocusSearch	key.Binding
	
	LineUp       key.Binding
	LineDown     key.Binding
	PageUp       key.Binding
	PageDown     key.Binding
	HalfPageUp   key.Binding
	HalfPageDown key.Binding
	GotoTop      key.Binding
	GotoBottom   key.Binding
}

func DefaultKeyMap() KeyMap {
	const spacebar = " "
	return KeyMap {
		Quit: key.NewBinding(
			key.WithKeys( "q", "ctrl+c" ),
			key.WithHelp( "q/<C-c>", "Exit the application" ),
		),
		Up: key.NewBinding(
			key.WithKeys( "k", "up" ),        
			key.WithHelp( "↑/k", "move up" ), 
		),
		Down: key.NewBinding(
			key.WithKeys( "j", "down" ),
			key.WithHelp( "↓/j", "move down" ),
		),
		ViewService: key.NewBinding(
			key.WithKeys( "enter" ),
			key.WithHelp( "enter", "View service details" ),
		),
		Refresh: key.NewBinding(
			key.WithKeys( "r" ),
			key.WithHelp( "r", "Force refresh" ),
		),
		Back: key.NewBinding(
			key.WithKeys( "esc" ),
			key.WithHelp( "esc", "Back to list" ),
		),
		ShowHelp: key.NewBinding(
			key.WithKeys( "?" ),
			key.WithHelp( "?", "Show help" ),
		),
		FocusSearch: key.NewBinding(
			key.WithKeys( "/" ),
			key.WithHelp( "/", "Focus search bar" ),
		),

		LineUp: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		LineDown: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("b", "pgup"),
			key.WithHelp("b/pgup", "page up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("f", "pgdown", spacebar),
			key.WithHelp("f/pgdn", "page down"),
		),
		HalfPageUp: key.NewBinding(
			key.WithKeys("u", "ctrl+u"),
			key.WithHelp("u", "½ page up"),
		),
		HalfPageDown: key.NewBinding(
			key.WithKeys("d", "ctrl+d"),
			key.WithHelp("d", "½ page down"),
		),
		GotoTop: key.NewBinding(
			key.WithKeys("home", "g"),
			key.WithHelp("g/home", "go to start"),
		),
		GotoBottom: key.NewBinding(
			key.WithKeys("end", "G"),
			key.WithHelp("G/end", "go to end"),
		),
	}
}

// ShortHelp implements the KeyMap interface.
func (km KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{km.LineUp, km.LineDown}
}

// FullHelp implements the KeyMap interface.
func (km KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{km.LineUp, km.LineDown, km.GotoTop, km.GotoBottom},
		{km.PageUp, km.PageDown, km.HalfPageUp, km.HalfPageDown},
	}
}
