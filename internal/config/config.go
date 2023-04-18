package config

import (
	"bufio"
	"fmt"
	"github.com/adrg/xdg"
	"go.uber.org/fx"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var Module = fx.Options(
	fx.Provide(NewConfig),
)

type Config struct {
	Home      string
	VarsFile  string
	HostsFile string
	LoggerConfig
	Variables
	Inventory
}

var (
	defaultVarsConfigPath  = "./internal/templates/vars/vars.tpl"
	defaultHostsConfigPath = "./inventories/production/inventory.yml"
	defaultLogLevel        = "info"
	defaultLogFormat       = "console"
	defaultLogDirectory    = "./"
	defaultWebitelUser     = ""
	defaultWebitelPass     = ""
)

func SetProperties(logLevel, logFormat, logFile, varsFile, inventoryFile, user, password string) {
	defaultLogLevel = logLevel
	defaultLogFormat = logFormat
	defaultLogDirectory = logFile
	defaultVarsConfigPath = varsFile
	defaultHostsConfigPath = inventoryFile
	defaultWebitelUser = user
	defaultWebitelPass = password
}

func NewConfig() Config {
	config := Config{
		Home:      defaultHome(),
		VarsFile:  defaultVarsConfigPath,
		HostsFile: defaultHostsConfigPath,
		LoggerConfig: LoggerConfig{
			LogLevel:     defaultLogLevel,
			LogFormat:    defaultLogFormat,
			LogDirectory: defaultLogDirectory,
		},
		Variables: Variables{
			WebitelRepositoryUser:     defaultWebitelUser,
			WebitelRepositoryPassword: defaultWebitelPass,
		},
	}

	config = createVarsConfigFromTpl(config)

	if err := config.readVarsToStruct(); err != nil {
		fmt.Println(err)
		return config
	}

	return config
}

func createVarsConfigFromTpl(config Config) Config {
	tpl, err := template.ParseFiles(config.VarsFile)
	if err != nil {
		fmt.Println(err)
	}

	cfgPath := getVarsConfigPath(config.WebitelRepositoryUser)
	config.VarsFile = cfgPath

	file, err := os.Create(cfgPath)
	if err != nil {
		fmt.Println(err)
	}

	err = tpl.Execute(file, config)
	if err != nil {
		fmt.Println(err)
	}

	return config
}

func getVarsConfigPath(user string) string {
	path, err := xdg.DataFile(fmt.Sprintf("wdeploy/%s/vars/all.yml", user))
	if err != nil {
		fmt.Println(err)
	}

	return path
}

func (c *Config) readVarsToStruct() error {
	file, err := os.Open(c.VarsFile)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := yaml.NewDecoder(file).Decode(&c.Variables); err != nil {
		return err
	}

	return nil
}

func (c *Config) GetVarsConfigContent() (fullText string, err error) {
	var text []string

	file, err := os.Open(c.VarsFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	fullText = strings.Join(text, "\n")
	if err := c.readVarsToStruct(); err != nil {
		return "", err
	}

	return fullText, nil
}

// default helpers for the configuration.
// We use $XDG_DATA_HOME to avoid cluttering the user's home directory.
func defaultHome() string {
	return filepath.Join(xdg.DataHome, "wdeploy")
}
