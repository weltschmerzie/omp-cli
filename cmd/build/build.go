package build

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/weltschmerzie/omp-cli/internal/builder"
)

// BuildCmd represents the build command
var BuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build/compile the open.mp project",
	Long: `Build command compiles the open.mp project.
It will look for the project files in the current directory
and compile them according to open.mp specifications.
It uses project.json for project configuration and config.json for server settings.`,
	DisableFlagParsing:    false,
	DisableAutoGenTag:     true,
	DisableFlagsInUseLine: false,
	DisableSuggestions:    true,
	Run: func(cmd *cobra.Command, args []string) {
		// Get flags
		verbose, _ := cmd.Flags().GetBool("verbose")

		// Execute build
		if err := builder.Build(verbose); err != nil {
			fmt.Printf("Error building project: %v\n", err)
			return
		}

		fmt.Println("Project built successfully!")
	},
}

func init() {
	// Add flags
	BuildCmd.Flags().BoolP("verbose", "v", false, "Enable verbose output")
}
