package builder

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/weltschmerzie/omp-cli/pkg/utils"
)

// BuildResult represents the result of a build operation
type BuildResult struct {
	Success  bool
	Errors   []string
	Warnings []string
}

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

	// Get server configuration
	serverConfig, err := utils.GetServerConfig()
	if err != nil {
		return fmt.Errorf("failed to get server configuration: %w", err)
	}

	if verbose {
		fmt.Println("Building open.mp project...")
		fmt.Printf("Project name: %s\n", config.Name)
		fmt.Printf("Project version: %s\n", config.Version)
		fmt.Printf("Server hostname: %s\n", serverConfig.Hostname)
		fmt.Printf("Using pawncc from: %s\n", config.PawnccPath)
		fmt.Printf("Main file: %s\n", config.MainFile)
		fmt.Printf("Output file: %s\n", config.OutputFile)
	}

	// Create build directory if it doesn't exist
	buildDir := filepath.Join(".", "build")
	if err := os.MkdirAll(buildDir, 0755); err != nil {
		return fmt.Errorf("failed to create build directory: %w", err)
	}

	// Create gamemodes directory in build directory
	gamemodesDir := filepath.Join(buildDir, "gamemodes")
	if err := os.MkdirAll(gamemodesDir, 0755); err != nil {
		return fmt.Errorf("failed to create gamemodes directory: %w", err)
	}

	// Determine pawncc executable path
	var pawnccExe string
	if config.PawnccPath != "" {
		if runtime.GOOS == "windows" {
			pawnccExe = filepath.Join(config.PawnccPath, "pawncc.exe")
		} else {
			pawnccExe = filepath.Join(config.PawnccPath, "pawncc")
		}
	} else {
		// Fallback to just "pawncc" and rely on PATH
		if runtime.GOOS == "windows" {
			pawnccExe = "pawncc.exe"
		} else {
			pawnccExe = "pawncc"
		}
	}

	// Check if pawncc exists
	if _, err := os.Stat(pawnccExe); os.IsNotExist(err) && config.PawnccPath != "" {
		// If not found at specified path, try to find in PATH
		if verbose {
			fmt.Printf("Warning: pawncc not found at %s, trying to find in PATH\n", pawnccExe)
		}

		// Fallback to just "pawncc" and rely on PATH
		if runtime.GOOS == "windows" {
			pawnccExe = "pawncc.exe"
		} else {
			pawnccExe = "pawncc"
		}
	}

	// Determine output file path
	outputPath := filepath.Join(buildDir, config.OutputFile)
	outputDir := filepath.Dir(outputPath)

	// Make sure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create command with output file option
	cmd := exec.Command(pawnccExe, "-o"+outputPath, config.MainFile)

	// Create buffers to capture output
	var stdout, stderr bytes.Buffer

	// Set up output streams
	if verbose {
		// For verbose mode, we want to see output in real-time and also capture it
		cmd.Stdout = io.MultiWriter(os.Stdout, &stdout)
		cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)
	} else {
		// For non-verbose mode, just capture the output
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
	}

	// Execute the compiler
	err = cmd.Run()

	// Parse the output for errors and warnings
	result := parseBuildOutput(stdout.String(), stderr.String())

	// Display the result
	if len(result.Errors) > 0 {
		fmt.Printf("\nBuild failed with %d errors and %d warnings.\n", len(result.Errors), len(result.Warnings))

		if !verbose {
			// Only show errors and warnings if not in verbose mode (to avoid duplication)
			fmt.Println("\nErrors:")
			for _, errMsg := range result.Errors {
				fmt.Printf("  - %s\n", errMsg)
			}

			if len(result.Warnings) > 0 {
				fmt.Println("\nWarnings:")
				for _, warnMsg := range result.Warnings {
					fmt.Printf("  - %s\n", warnMsg)
				}
			}
		}

		return fmt.Errorf("compilation failed with %d errors", len(result.Errors))
	} else if len(result.Warnings) > 0 {
		fmt.Printf("\nBuild completed with %d warnings.\n", len(result.Warnings))

		if !verbose {
			// Only show warnings if not in verbose mode (to avoid duplication)
			fmt.Println("\nWarnings:")
			for _, warnMsg := range result.Warnings {
				fmt.Printf("  - %s\n", warnMsg)
			}
		}
	} else {
		fmt.Println("\nBuild completed successfully with 0 errors and 0 warnings.")
	}

	if err != nil {
		return fmt.Errorf("compilation process failed: %w", err)
	}

	// Copy necessary files to build directory
	if err := utils.CopyRequiredFiles(buildDir); err != nil {
		return fmt.Errorf("failed to copy required files: %w", err)
	}

	return nil
}

// parseBuildOutput parses the compiler output to extract errors and warnings
func parseBuildOutput(stdout, stderr string) BuildResult {
	result := BuildResult{
		Success:  true,
		Errors:   []string{},
		Warnings: []string{},
	}

	// Combine stdout and stderr
	output := stdout + "\n" + stderr

	// Regular expressions for error and warning detection
	errorRegex := regexp.MustCompile(`(?i)(error|fatal error|undefined symbol|cannot find|not found).*`)
	warningRegex := regexp.MustCompile(`(?i)(warning|note|suggestion).*`)

	// Scan through the output line by line
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()

		// Check for errors
		if errorRegex.MatchString(line) {
			result.Errors = append(result.Errors, strings.TrimSpace(line))
			result.Success = false
		} else if warningRegex.MatchString(line) {
			result.Warnings = append(result.Warnings, strings.TrimSpace(line))
		}
	}

	return result
}
