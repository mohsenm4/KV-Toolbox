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
	Comment     string `mapstructure:"comment"`
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

	// Determine executable directory
	execPath, _ := os.Executable()
	execDir := filepath.Dir(execPath)
	configFile := filepath.Join(execDir, "data.json")

	v.SetConfigFile(configFile)

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
		Comment:     data["Comment"],
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
