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
	Resources  []string `json:"resources"`
	Plugins    []string `json:"plugins"`
	ServerCfg  string   `json:"server_cfg"`
	Author     string   `json:"author"`
	Repository string   `json:"repository"`
}

// IsOpenMPProject checks if the current directory is an open.mp project
func IsOpenMPProject() bool {
	// Check for project.json file
	if _, err := os.Stat("project.json"); !os.IsNotExist(err) {
		return true
	}

	// Check for server.cfg file
	if _, err := os.Stat("server.cfg"); !os.IsNotExist(err) {
		return true
	}

	// Check for pawn scripts
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
		Name:      "gamemode",
		Version:   "1.0.0",
		MainFile:  "gamemode.pwn",
		Resources: []string{},
		Plugins:   []string{},
		ServerCfg: "server.cfg",
	}

	// Check for main script file
	matches, err := filepath.Glob("*.pwn")
	if err == nil && len(matches) > 0 {
		config.MainFile = matches[0]
		config.Name = filepath.Base(matches[0])
		config.Name = config.Name[:len(config.Name)-4] // Remove .pwn extension
	}

	return config, nil
}

// CopyRequiredFiles copies necessary files to the build directory
func CopyRequiredFiles(buildDir string) error {
	// Get project configuration
	config, err := GetProjectConfig()
	if err != nil {
		return err
	}

	// Copy server.cfg
	if _, err := os.Stat(config.ServerCfg); !os.IsNotExist(err) {
		if err := copyFile(config.ServerCfg, filepath.Join(buildDir, "server.cfg")); err != nil {
			return fmt.Errorf("failed to copy server.cfg: %w", err)
		}
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
