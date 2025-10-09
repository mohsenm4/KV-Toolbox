package configApp

import (
	"fmt"

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
	v.SetConfigName("data")
	v.SetConfigType("json")
	v.AddConfigPath(".")
	v.ReadInConfig()
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
			return fmt.Errorf("This database has already been added to your projects under the name '%s'", proj.Name), true
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
