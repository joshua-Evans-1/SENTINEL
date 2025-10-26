package cmd

import (
	"fmt"
	"os"

	"github.com/0xReLogic/SENTINEL/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// tuiCmd represents the tui command
var tuiCmd = &cobra.Command{
	Use:   cmdNameTui,
	Short: descTuiShort,
	Long:  descTuiLong,
	Run: func(cmd *cobra.Command, args []string) {
		// load configuration
		cfg, err := loadConfig(configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, errLoadingConfig, err)
			os.Exit(exitConfigError)
		}
		
		p := tea.NewProgram(
			tui.InitialRootModel( cfg ),
			tea.WithAltScreen(),
			tea.WithMouseCellMotion(),
		)

		if _, err := p.Run(); err != nil {
			print( fmt.Errorf( "error : %v", err ) )
		}
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}

