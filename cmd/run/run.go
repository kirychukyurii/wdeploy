package run

import (
	"github.com/kirychukyurii/wdeploy/bootstrap"
	"github.com/kirychukyurii/wdeploy/internal/config"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var (
	inventoryFile string
	varsFile      string
	logLevel      string
	logFormat     string
	logFile       string
	user          string
	password      string
	inventoryType string
)

var (
	inventoryTemplateLocalhost = "./inventories/hosts/localhost.yml"
	inventoryTemplateCustom    = "./inventories/hosts/custom.yml"
)

func init() {
	pf := Command.PersistentFlags()
	pf.StringVarP(&logLevel, "log-level", "l",
		"debug", "log level: debug, info, warn, error, dpanic, panic, fatal")
	pf.StringVarP(&logFormat, "log-format", "F",
		"plain", "log format output: json, console")
	pf.StringVarP(&logFile, "log-path", "L",
		"./", "log file location")
	pf.StringVarP(&varsFile, "vars", "V",
		"./internal/templates/vars/vars.tpl", "variables file")
	pf.StringVarP(&inventoryFile, "inventory", "i",
		"./inventories/production/inventory.yml", "hosts file")
	pf.StringVarP(&user, "user", "u",
		"", "webitel repository user")
	pf.StringVarP(&password, "password", "p",
		"", "webitel repository password")
	pf.StringVarP(&inventoryType, "type", "t",
		"localhost", "inventory template type: localhost, custom")
}

var Command = &cobra.Command{
	Use:          "run",
	Short:        "Start Ansible Playbook",
	Example:      "wdeploy run -V inventories/production/group_vars/all.yml -i inventories/production/inventory.yml",
	SilenceUsage: true,
	PreRun: func(cmd *cobra.Command, args []string) {
		if inventoryType == "localhost" {
			inventoryFile = inventoryTemplateLocalhost
		} else {
			inventoryFile = inventoryTemplateCustom
		}

		config.SetProperties(logLevel, logFormat, logFile, varsFile, inventoryFile, user, password)
	},
	Run: func(cmd *cobra.Command, args []string) {
		runApplication()
	},
}

func runApplication() {
	fx.New(bootstrap.Module, fx.NopLogger).Run()
	//fx.New(bootstrap.Module).Run()
}
