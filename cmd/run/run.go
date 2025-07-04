package run

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/weltschmerzie/omp-cli/internal/runner"
)

// RunCmd represents the run command
var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the open.mp project",
	Long: `Run command executes the open.mp project.
It will look for the compiled project in the current directory
and run it according to open.mp specifications using config.json.`,
	DisableFlagParsing:    false,
	DisableAutoGenTag:     true,
	DisableFlagsInUseLine: false,
	DisableSuggestions:    true,
	Run: func(cmd *cobra.Command, args []string) {
		// Get flags
		debug, _ := cmd.Flags().GetBool("debug")
		port, _ := cmd.Flags().GetInt("port")

		// Execute run
		if err := runner.Run(debug, port); err != nil {
			fmt.Printf("Error running project: %v\n", err)
			return
		}

		fmt.Println("Project is running. Press Ctrl+C to stop.")
	},
}

func init() {
	// Add flags
	RunCmd.Flags().BoolP("debug", "d", false, "Enable debug mode")
	RunCmd.Flags().IntP("port", "p", 7777, "Port to run the server on")
}
