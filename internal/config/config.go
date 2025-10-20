package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Project struct {
	Name        string `mapstructure:"name"`
	Comment     string `mapstructure:"comment"`
	FileAddress string `mapstructure:"fileAddress"`
	Database    string `mapstructure:"database"`
}
type Config struct {
	RecentProjects []Project `mapstructure:"recentProjects"`
}

func LoadConfig() *Config {
	v := viper.New()
	v.SetConfigType("json")

	// Determine configuration file path using standard locations
	configFile, err := getConfigFilePath()
	fmt.Println("ConfigPath :", configFile)
	if err != nil {
		fmt.Printf("Warning: Could not determine config file path: %v. Using fallback.\n", err)
		configFile = "data.json" // fallback to current directory
	}

	v.SetConfigFile(configFile)

	fmt.Println("Using configuration file:", v.ConfigFileUsed())

	// Try to read existing config
	if err := v.ReadInConfig(); err != nil {
		var pathErr *os.PathError
		if errors.As(err, &pathErr) && os.IsNotExist(pathErr.Err) {

			fmt.Println("data.json not found, creating new one...")
			v.SafeWriteConfigAs(configFile)
			fmt.Println("data.json created at:", configFile)
		} else {
			fmt.Println("Error reading config file:", err)
		}
	}

	var config Config

	err = v.Unmarshal(&config)
	if err != nil {
		fmt.Println("Error unmarshalling config:", err)
		return nil
	}

	return &config
}

func (c *Config) Write() error {
	c.v.Set("recentProjects", c.RecentProjects)
	return c.v.WriteConfig()
}

func (c *Config) Add(data map[string]string) (error, bool) {

	for _, proj := range c.RecentProjects {
		if data["Addres"] == proj.FileAddress {
			return fmt.Errorf("this database has already been added to your projects under the name '%s'", proj.Name), true
		}
	}

	newProject := Project{
		Name:        data["Name"],
		Comment:     data["Comment"],
		FileAddress: data["Addres"],
		Database:    data["Database"],
	}
	c.RecentProjects = append(c.RecentProjects, newProject)
	return nil, false
}

func (c *Config) Remove(projectName string) {
	for i, proj := range c.RecentProjects {
		if proj.Name == projectName {
			c.RecentProjects = append(c.RecentProjects[:i], c.RecentProjects[i+1:]...)
			break
		}
	}

	return
}

// getConfigFilePath determines the best location for the config file
func getConfigFilePath() (string, error) {
	// Priority 1: User config directory (most appropriate for user data)
	if userConfigDir, err := os.UserConfigDir(); err == nil {
		appConfigDir := filepath.Join(userConfigDir, "KV-Toolbox")
		if err := os.MkdirAll(appConfigDir, 0755); err == nil {
			configFile := filepath.Join(appConfigDir, "data.json")
			// Check if we can write to this directory
			if checkWritePermission(appConfigDir) == nil {
				return configFile, nil
			}
		}
	}

	// Priority 2: Home directory fallback
	if homeDir, err := os.UserHomeDir(); err == nil {
		configFile := filepath.Join(homeDir, ".kv-toolbox", "data.json")
		configDir := filepath.Dir(configFile)
		if err := os.MkdirAll(configDir, 0755); err == nil {
			if checkWritePermission(configDir) == nil {
				return configFile, nil
			}
		}
	}

	// Priority 3: Current working directory (last resort)
	if wd, err := os.Getwd(); err == nil {
		configFile := filepath.Join(wd, "data.json")
		if checkWritePermission(wd) == nil {
			return configFile, nil
		}
	}

	// Final fallback: current directory (even if we can't verify write permission)
	return "data.json", fmt.Errorf("unable to find suitable config directory, using current directory")
}

// checkWritePermission checks if we can write to a directory
func checkWritePermission(dir string) error {
	testFile := filepath.Join(dir, ".write_test")
	file, err := os.Create(testFile)
	if err != nil {
		return err
	}
	file.Close()
	return os.Remove(testFile)
}
