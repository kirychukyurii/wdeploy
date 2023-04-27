package config

import (
	"bufio"
	"fmt"
	"github.com/adrg/xdg"
	"github.com/kirychukyurii/wdeploy/internal/pkg/file"
	"github.com/kirychukyurii/wdeploy/internal/templates/inventory/custom"
	"github.com/kirychukyurii/wdeploy/internal/templates/inventory/localhost"
	"github.com/kirychukyurii/wdeploy/internal/templates/vars"
	"go.uber.org/fx"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var Module = fx.Options(
	fx.Provide(New),
)

type Config struct {
	PlaybookFile  string
	VarsFile      string
	HostsFile     string
	InventoryType string
	LoggerConfig
	Variables
	Inventory
}

var DefaultConfig = Config{
	PlaybookFile:  "playbook.yml",
	VarsFile:      "",
	HostsFile:     "",
	InventoryType: "custom",
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

	if config.VarsFile == "" {
		config.VarsFile = getConfigPath("vars", config.WebitelRepositoryUser)
	}
	if config.HostsFile == "" {
		config.HostsFile = getConfigPath("hosts", config.WebitelRepositoryUser)
	}

	fmt.Println(config.VarsFile)
	fmt.Println(config.HostsFile)

	if !file.IsFile(config.VarsFile) {
		config = createVarsConfigFromTpl(config)
	}

	if err := config.readVarsToStruct(); err != nil {
		fmt.Println(err)
		return config
	}

	if !file.IsFile(config.HostsFile) {
		config = createHostsConfigFromTpl(config)
	}

	if err := config.readHostsToStruct(); err != nil {
		fmt.Println(err)

		return config
	}

	return config
}

func createVarsConfigFromTpl(config Config) Config {
	cfgPath, err := xdg.DataFile(fmt.Sprintf("wdeploy/%s/vars/vars.yml", config.WebitelRepositoryUser))
	if err != nil {
		fmt.Println(err)
	}

	tpl, err := template.New("").Parse(vars.Tmpl)
	if err != nil {
		fmt.Println(err)
	}

	f, err := os.Create(cfgPath)
	if err != nil {
		fmt.Println(err)
	}

	err = tpl.Execute(f, config)
	if err != nil {
		fmt.Println(err)
	}

	return config
}

func createHostsConfigFromTpl(config Config) Config {
	cfgPath, err := xdg.DataFile(fmt.Sprintf("wdeploy/%s/hosts/hosts.yml", config.WebitelRepositoryUser))
	if err != nil {
		fmt.Println(err)
	}

	hostsTpl := custom.Tmpl
	if config.InventoryType == "local" {
		hostsTpl = localhost.Tmpl
	}

	tpl, err := template.New("").Parse(hostsTpl)
	if err != nil {
		fmt.Println(err)
	}

	f, err := os.Create(cfgPath)
	if err != nil {
		fmt.Println(err)
	}

	err = tpl.Execute(f, config)
	if err != nil {
		fmt.Println(err)
	}

	return config
}

func getConfigPath(cfgType, user string) (cfgPath string) {
	return filepath.Join(xdg.DataHome, fmt.Sprintf("wdeploy/%s/%s/%s.yml", user, cfgType, cfgType))
}

func (c *Config) readVarsToStruct() error {
	f, err := os.Open(c.VarsFile)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(&c.Variables); err != nil {
		return err
	}

	return nil
}

func (c *Config) readHostsToStruct() error {
	f, err := os.Open(c.HostsFile)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(&c.Inventory); err != nil {
		return err
	}

	return nil
}

func (c *Config) GetVarsConfigContent() (fullText string, err error) {
	var text []string

	f, err := os.Open(c.VarsFile)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
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

func (c *Config) GetHostsConfigContent() (fullText string, err error) {
	var text []string

	f, err := os.Open(c.HostsFile)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	fullText = strings.Join(text, "\n")
	if err := c.readHostsToStruct(); err != nil {
		return "", err
	}

	return fullText, nil
}
