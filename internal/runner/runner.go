package runner

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/weltschmerzie/omp-cli/pkg/utils"
)

// Run executes the open.mp project
func Run(debug bool, port int) error {
	// Check if we are in an open.mp project directory
	if !utils.IsOpenMPProject() {
		return errors.New("current directory is not an open.mp project")
	}

	// Check if the project is built
	buildDir := filepath.Join(".", "build")
	if _, err := os.Stat(buildDir); os.IsNotExist(err) {
		return errors.New("project is not built. Please run 'ompcli build' first")
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

	// Determine server executable based on OS
	var serverExe string
	switch runtime.GOOS {
	case "windows":
		serverExe = "omp-server.exe"
	case "linux", "darwin":
		serverExe = "omp-server"
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	// Check if server executable exists
	serverPath := filepath.Join(buildDir, serverExe)
	if _, err := os.Stat(serverPath); os.IsNotExist(err) {
		return fmt.Errorf("server executable not found at %s", serverPath)
	}

	// Check if the compiled gamemode exists
	gamemodePath := filepath.Join(buildDir, config.OutputFile)
	if _, err := os.Stat(gamemodePath); os.IsNotExist(err) {
		return fmt.Errorf("compiled gamemode not found at %s. Please run 'ompcli build' first", gamemodePath)
	}

	// Prepare command arguments
	args := []string{}

	// Add debug flag if enabled
	if debug {
		args = append(args, "--debug")
	}

	// Override port if specified
	if port != 0 && port != serverConfig.Port {
		args = append(args, "--port="+strconv.Itoa(port))
	} else {
		port = serverConfig.Port
	}

	// Add gamemode if not specified in config.json
	if serverConfig.Gamemode == "" {
		// Use relative path from build directory to gamemode file
		relativeGamemodePath := filepath.Join("gamemodes", filepath.Base(config.OutputFile))
		args = append(args, "--gamemode="+relativeGamemodePath)
	}

	// Create command
	cmd := exec.Command(serverPath, args...)

	// Set working directory to build directory
	cmd.Dir = buildDir

	// Connect standard I/O
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Run the server
	fmt.Printf("Starting open.mp server on port %d...\n", port)
	if debug {
		fmt.Println("Debug mode enabled")
	}
	fmt.Printf("Using gamemode: %s\n", filepath.Base(config.OutputFile))

	return cmd.Run()
}
