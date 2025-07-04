package init

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// Project represents the structure of project.json
type Project struct {
	Name       string   `json:"name"`
	Version    string   `json:"version"`
	MainFile   string   `json:"main_file"`
	OutputFile string   `json:"output_file"`
	Resources  []string `json:"resources"`
	Plugins    []string `json:"plugins"`
	ServerCfg  string   `json:"server_cfg"`
	Author     string   `json:"author"`
	Repository string   `json:"repository"`
	PawnccPath string   `json:"pawncc_path"`
}

// Server represents the structure of config.json
type Server struct {
	Hostname     string   `json:"hostname"`
	Port         int      `json:"port"`
	MaxPlayers   int      `json:"maxplayers"`
	Language     string   `json:"language"`
	Gamemode     string   `json:"gamemode"`
	Plugins      []string `json:"plugins"`
	WebURL       string   `json:"weburl"`
	RCONPassword string   `json:"rcon_password"`
	Password     string   `json:"password"`
}

// InitCmd represents the init command
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new open.mp project",
	Long: `Init command creates a new open.mp project by generating
the necessary configuration files (project.json and config.json).`,
	DisableFlagParsing:    false,
	DisableAutoGenTag:     true,
	DisableFlagsInUseLine: false,
	DisableSuggestions:    true,
	Run: func(cmd *cobra.Command, args []string) {
		// Get flags
		name, _ := cmd.Flags().GetString("name")
		author, _ := cmd.Flags().GetString("author")
		pawnccPath, _ := cmd.Flags().GetString("pawncc-path")

		// If name is not provided, use current directory name
		if name == "" {
			currentDir, err := os.Getwd()
			if err == nil {
				name = filepath.Base(currentDir)
			} else {
				name = "gamemode"
			}
		}

		// If pawncc path is not provided, use default
		if pawnccPath == "" {
			pawnccPath = "qawno"
		}

		// Create project.json
		if err := createProjectJson(name, author, pawnccPath); err != nil {
			fmt.Printf("Error creating project.json: %v\n", err)
			return
		}

		// Create config.json
		if err := createConfigJson(name); err != nil {
			fmt.Printf("Error creating config.json: %v\n", err)
			return
		}

		// Create gamemodes directory if it doesn't exist
		if err := os.MkdirAll("gamemodes", 0755); err != nil {
			fmt.Printf("Warning: Failed to create gamemodes directory: %v\n", err)
		}

		fmt.Println("Project initialized successfully!")
		fmt.Println("Created files:")
		fmt.Println("- project.json")
		fmt.Println("- config.json")
		fmt.Println("- gamemodes/ directory")
	},
}

func init() {
	// Add flags
	InitCmd.Flags().StringP("name", "n", "", "Project name (default: current directory name)")
	InitCmd.Flags().StringP("author", "a", "", "Project author")
	InitCmd.Flags().String("pawncc-path", "", "Path to pawncc compiler (default: qawno)")
}

// createProjectJson creates a project.json file
func createProjectJson(name, author, pawnccPath string) error {
	// Check if file already exists
	if _, err := os.Stat("project.json"); !os.IsNotExist(err) {
		return fmt.Errorf("project.json already exists")
	}

	// Create project configuration
	project := Project{
		Name:       name,
		Version:    "1.0.0",
		MainFile:   filepath.Join("gamemodes", name+".pwn"),
		OutputFile: filepath.Join("gamemodes", name+".amx"),
		Resources:  []string{},
		Plugins:    []string{},
		ServerCfg:  "config.json",
		Author:     author,
		Repository: "",
		PawnccPath: pawnccPath,
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	if err := os.WriteFile("project.json", jsonData, 0644); err != nil {
		return err
	}

	return nil
}

// createConfigJson creates a config.json file
func createConfigJson(name string) error {
	// Check if file already exists
	if _, err := os.Stat("config.json"); !os.IsNotExist(err) {
		return fmt.Errorf("config.json already exists")
	}

	// Create server configuration
	server := Server{
		Hostname:     name + " Server",
		Port:         7777,
		MaxPlayers:   50,
		Language:     "English",
		Gamemode:     name,
		Plugins:      []string{},
		WebURL:       "open.mp",
		RCONPassword: "changeme",
		Password:     "",
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(server, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	if err := os.WriteFile("config.json", jsonData, 0644); err != nil {
		return err
	}

	return nil
}
