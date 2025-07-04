package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ProjectConfig represents the configuration of an open.mp project
type ProjectConfig struct {
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

// ServerConfig represents the configuration of an open.mp server
type ServerConfig struct {
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

// IsOpenMPProject checks if the current directory is an open.mp project
func IsOpenMPProject() bool {
	// Check for project.json file
	if _, err := os.Stat("project.json"); !os.IsNotExist(err) {
		return true
	}

	// Check for config.json file
	if _, err := os.Stat("config.json"); !os.IsNotExist(err) {
		return true
	}

	// Check for pawn scripts in gamemodes directory
	if _, err := os.Stat("gamemodes"); !os.IsNotExist(err) {
		matches, err := filepath.Glob(filepath.Join("gamemodes", "*.pwn"))
		if err == nil && len(matches) > 0 {
			return true
		}
	}

	// Check for pawn scripts in root directory (legacy support)
	matches, err := filepath.Glob("*.pwn")
	if err == nil && len(matches) > 0 {
		return true
	}

	return false
}

// GetProjectConfig reads and parses the project configuration
func GetProjectConfig() (*ProjectConfig, error) {
	// Try to read project.json first
	if _, err := os.Stat("project.json"); !os.IsNotExist(err) {
		data, err := os.ReadFile("project.json")
		if err != nil {
			return nil, fmt.Errorf("failed to read project.json: %w", err)
		}

		var config ProjectConfig
		if err := json.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse project.json: %w", err)
		}

		return &config, nil
	}

	// If project.json doesn't exist, try to infer configuration
	config := &ProjectConfig{
		Name:       "gamemode",
		Version:    "1.0.0",
		MainFile:   filepath.Join("gamemodes", "gamemode.pwn"),
		OutputFile: filepath.Join("gamemodes", "gamemode.amx"),
		Resources:  []string{},
		Plugins:    []string{},
		ServerCfg:  "config.json", // Updated to config.json
		PawnccPath: "qawno",       // Default pawncc path
	}

	// Check for main script file in gamemodes directory
	if _, err := os.Stat("gamemodes"); !os.IsNotExist(err) {
		matches, err := filepath.Glob(filepath.Join("gamemodes", "*.pwn"))
		if err == nil && len(matches) > 0 {
			config.MainFile = matches[0]
			baseName := filepath.Base(matches[0])
			nameWithoutExt := baseName[:len(baseName)-4] // Remove .pwn extension
			config.Name = nameWithoutExt
			config.OutputFile = filepath.Join("gamemodes", nameWithoutExt+".amx")
		}
	} else {
		// Legacy support: Check for main script file in root directory
		matches, err := filepath.Glob("*.pwn")
		if err == nil && len(matches) > 0 {
			config.MainFile = matches[0]
			baseName := filepath.Base(matches[0])
			nameWithoutExt := baseName[:len(baseName)-4] // Remove .pwn extension
			config.Name = nameWithoutExt
			config.OutputFile = filepath.Join("gamemodes", nameWithoutExt+".amx")
		}
	}

	return config, nil
}

// GetServerConfig reads and parses the server configuration
func GetServerConfig() (*ServerConfig, error) {
	// Try to read config.json
	if _, err := os.Stat("config.json"); !os.IsNotExist(err) {
		data, err := os.ReadFile("config.json")
		if err != nil {
			return nil, fmt.Errorf("failed to read config.json: %w", err)
		}

		var config ServerConfig
		if err := json.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse config.json: %w", err)
		}

		return &config, nil
	}

	// If config.json doesn't exist, return default configuration
	return &ServerConfig{
		Hostname:   "Open.MP Server",
		Port:       7777,
		MaxPlayers: 50,
		Language:   "English",
		Gamemode:   "gamemode",
		Plugins:    []string{},
		WebURL:     "open.mp",
	}, nil
}

// CopyRequiredFiles copies necessary files to the build directory
func CopyRequiredFiles(buildDir string) error {
	// Get project configuration
	config, err := GetProjectConfig()
	if err != nil {
		return err
	}

	// Copy config.json
	if _, err := os.Stat("config.json"); !os.IsNotExist(err) {
		if err := copyFile("config.json", filepath.Join(buildDir, "config.json")); err != nil {
			return fmt.Errorf("failed to copy config.json: %w", err)
		}
	} else if _, err := os.Stat(config.ServerCfg); !os.IsNotExist(err) {
		// For backward compatibility, also check for server.cfg
		if err := copyFile(config.ServerCfg, filepath.Join(buildDir, "config.json")); err != nil {
			return fmt.Errorf("failed to copy server configuration: %w", err)
		}
	}

	// Create gamemodes directory in build directory
	gamemodesDir := filepath.Join(buildDir, "gamemodes")
	if err := os.MkdirAll(gamemodesDir, 0755); err != nil {
		return fmt.Errorf("failed to create gamemodes directory: %w", err)
	}

	// Copy resources
	for _, resource := range config.Resources {
		if _, err := os.Stat(resource); !os.IsNotExist(err) {
			destPath := filepath.Join(buildDir, filepath.Base(resource))
			if err := copyFile(resource, destPath); err != nil {
				return fmt.Errorf("failed to copy resource %s: %w", resource, err)
			}
		}
	}

	// Copy plugins
	pluginsDir := filepath.Join(buildDir, "plugins")
	if err := os.MkdirAll(pluginsDir, 0755); err != nil {
		return fmt.Errorf("failed to create plugins directory: %w", err)
	}

	for _, plugin := range config.Plugins {
		if _, err := os.Stat(plugin); !os.IsNotExist(err) {
			destPath := filepath.Join(pluginsDir, filepath.Base(plugin))
			if err := copyFile(plugin, destPath); err != nil {
				return fmt.Errorf("failed to copy plugin %s: %w", plugin, err)
			}
		}
	}

	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}
