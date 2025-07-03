package builder

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/yourusername/omp-cli/pkg/utils"
)

// Build compiles the open.mp project
func Build(verbose bool) error {
	// Check if we are in an open.mp project directory
	if !utils.IsOpenMPProject() {
		return errors.New("current directory is not an open.mp project")
	}

	// Get project configuration
	config, err := utils.GetProjectConfig()
	if err != nil {
		return fmt.Errorf("failed to get project configuration: %w", err)
	}

	if verbose {
		fmt.Println("Building open.mp project...")
		fmt.Printf("Project name: %s\n", config.Name)
		fmt.Printf("Project version: %s\n", config.Version)
	}

	// Create build directory if it doesn't exist
	buildDir := filepath.Join(".", "build")
	if err := os.MkdirAll(buildDir, 0755); err != nil {
		return fmt.Errorf("failed to create build directory: %w", err)
	}

	// Determine compiler command based on OS
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("pawncc", "-D"+buildDir, config.MainFile)
	case "linux", "darwin":
		cmd = exec.Command("pawncc", "-D"+buildDir, config.MainFile)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	// Set output for verbose mode
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	// Execute the compiler
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("compilation failed: %w", err)
	}

	// Copy necessary files to build directory
	if err := utils.CopyRequiredFiles(buildDir); err != nil {
		return fmt.Errorf("failed to copy required files: %w", err)
	}

	return nil
}
