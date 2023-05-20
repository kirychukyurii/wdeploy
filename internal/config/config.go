package config

import (
	"fmt"
	"github.com/adrg/xdg"
	"github.com/kirychukyurii/wdeploy/internal/constants"
	"github.com/kirychukyurii/wdeploy/internal/lib/file"
	"github.com/kirychukyurii/wdeploy/internal/templates/inventory/custom"
	"github.com/kirychukyurii/wdeploy/internal/templates/inventory/localhost"
	"github.com/kirychukyurii/wdeploy/internal/templates/vars"
	"go.uber.org/fx"
	"gopkg.in/yaml.v3"
	"path/filepath"
	"regexp"
	"text/template"
)

var Module = fx.Options(
	fx.Provide(New),
)

const (
	VarsConfig int = iota
	InventoryConfig
	lastsConfig
)

type Config struct {
	PlaybookRepositoryUrl string
	PlaybookTempDir       string
	ConfigFiles           []string
	InventoryType         string
	LoggerConfig
	Variables
	Inventory
}

var DefaultConfig = Config{
	PlaybookRepositoryUrl: "https://github.com/kirychukyurii/wansible",
	ConfigFiles:           make([]string, 2),
	InventoryType:         "custom",
	LoggerConfig: LoggerConfig{
		LogLevel:     "info",
		LogFormat:    "console",
		LogDirectory: "./",
	},
	Variables: Variables{
		WebitelRepositoryUser:     "",
		WebitelRepositoryPassword: "",
	},
}

func New() Config {
	config := DefaultConfig
	home := config.getUserLocalHome()

	configFilesType := make(map[int]string, lastsConfig)
	configFilesType[VarsConfig] = "vars"
	configFilesType[InventoryConfig] = "inventory"

	for i, v := range config.ConfigFiles {
		if v == "" {
			if err := file.EnsureDir(filepath.Join(home, configFilesType[i])); err != nil {
				fmt.Println("file.EnsureDir(): " + err.Error())
			}

			config.ConfigFiles[i] = filepath.Join(home, configFilesType[i], "all.yml")
			fmt.Printf("%s: %s\n", configFilesType[i], config.ConfigFiles[i])
		}

		if !file.IsFile(config.ConfigFiles[i]) {
			if err := config.createConfigFromTpl(i); err != nil {
				fmt.Println("config.createConfigFromTpl(i): " + err.Error())
			}
		}

		if err := config.ReadToStruct(i); err != nil {
			fmt.Println("config.ReadToStruct(i): ", err.Error())
		}
	}

	config.LogDirectory = filepath.Join(home, "logs")
	if err := file.EnsureDir(config.LogDirectory); err != nil {
		fmt.Println("file.EnsureDir(): " + err.Error())
	}

	ansibleLogLocation := config.GetAnsibleLogLocation()

	if !file.IsFile(ansibleLogLocation) {
		f, err := file.Create(ansibleLogLocation)
		if err != nil {
			fmt.Println(err)
		}

		defer file.Close(f)
	}

	return config
}

func (c *Config) getUserLocalHome() string {
	trimUser := regexp.MustCompile("[^a-zA-Z0-9]+").ReplaceAllString(c.WebitelRepositoryUser, "")
	return filepath.Join(xdg.DataHome, constants.AppName, trimUser)
}

func (c *Config) createConfigFromTpl(configFileType int) error {
	var tmpl string

	switch configFileType {
	case VarsConfig:
		tmpl = vars.Tmpl
	case InventoryConfig:
		if c.InventoryType == "local" {
			tmpl = localhost.Tmpl
		} else {
			tmpl = custom.Tmpl
		}
	}

	t, err := template.New("").Parse(tmpl)
	if err != nil {
		return err
	}

	f, err := file.Create(c.ConfigFiles[configFileType])
	defer file.Close(f)
	if err != nil {
		return err
	}

	err = t.Execute(f, c)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) GetAnsibleLogLocation() string {
	return filepath.Join(c.LogDirectory, "ansible.log")
}

func (c *Config) ReadToStruct(configFileType int) error {
	f, err := file.Open(c.ConfigFiles[configFileType])
	defer file.Close(f)
	if err != nil {
		return err
	}

	switch configFileType {
	case VarsConfig:
		if err = yaml.NewDecoder(f).Decode(&c.Variables); err != nil {
			return err
		}
	case InventoryConfig:
		if err = yaml.NewDecoder(f).Decode(&c.Inventory); err != nil {
			return err
		}
	}

	return nil
}
