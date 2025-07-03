package cmd

import (
	"github.com/spf13/cobra"
	buildCmd "github.com/weltschmerzie/omp-cli/cmd/build"
	runCmd "github.com/weltschmerzie/omp-cli/cmd/run"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ompcli",
	Short: "A CLI tool for open.mp projects",
	Long: `ompcli is a command line interface tool for open.mp projects.
It allows you to build and run open.mp projects easily.
	
For example:
  ompcli build - Builds/compiles the open.mp project
  ompcli run   - Runs the open.mp project`,
}

func init() {
	// Add subcommands
	RootCmd.AddCommand(buildCmd.BuildCmd)
	RootCmd.AddCommand(runCmd.RunCmd)
}
