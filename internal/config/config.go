package configApp

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Project struct {
	Name        string `mapstructure:"name"`
	FileAddress string `mapstructure:"fileAddress"`
	Databace    string `mapstructure:"databace"`
}

type JsonInformation struct {
	RecentProjects []Project `mapstructure:"recentProjects"`
}

type Config struct {
	v *viper.Viper
}

func NewConfig() *Config {
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

	return &Config{v: v}
}

func (c *Config) Load() (JsonInformation, error) {
	var state JsonInformation
	if err := c.v.Unmarshal(&state); err != nil {
		return state, err
	}
	return state, nil
}

func (c *Config) Write(state JsonInformation) error {
	c.v.Set("recentProjects", state.RecentProjects)
	return c.v.WriteConfig()
}

func (c *Config) Add(data map[string]string) (error, bool) {
	state, err := c.Load()
	if err != nil {
		return err, false
	}

	for _, proj := range state.RecentProjects {
		if data["Addres"] == proj.FileAddress {
			return fmt.Errorf("this database has already been added to your projects under the name '%s'", proj.Name), true
		}
	}

	newProject := Project{
		Name:        data["Name"],
		FileAddress: data["Addres"],
		Databace:    data["Database"],
	}
	state.RecentProjects = append(state.RecentProjects, newProject)
	return c.Write(state), false
}

func (c *Config) Remove(projectName string) error {
	state, err := c.Load()
	if err != nil {
		return err
	}

	for i, proj := range state.RecentProjects {
		if proj.Name == projectName {
			state.RecentProjects = append(state.RecentProjects[:i], state.RecentProjects[i+1:]...)
			break
		}
	}

	return c.Write(state)
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
