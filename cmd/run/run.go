package run

import (
	"github.com/kirychukyurii/wdeploy/bootstrap"
	"github.com/kirychukyurii/wdeploy/internal/config"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var logLevel string
var logFormat string
var logFile string
var varsFile string
var inventoryFile string

func init() {
	pf := Command.PersistentFlags()
	pf.StringVarP(&logLevel, "log-level", "l",
		"debug", "log level")
	pf.StringVarP(&logFormat, "log-format", "F",
		"plain", "log format output")
	pf.StringVarP(&logFile, "log-path", "L",
		"./", "log file location")
	pf.StringVarP(&varsFile, "vars", "V",
		"./inventories/production/group_vars/all.yml", "variables file")
	pf.StringVarP(&inventoryFile, "inventory", "i",
		"./inventories/production/inventory.yml", "hosts file")
}

var Command = &cobra.Command{
	Use:          "run",
	Short:        "Start Ansible Playbook",
	Example:      "wdeploy run -V inventories/production/group_vars/all.yml -i inventories/production/inventory.yml",
	SilenceUsage: true,
	PreRun: func(cmd *cobra.Command, args []string) {
		config.SetVarConfigPath(varsFile)
		config.SetHostsConfigPath(inventoryFile)
		config.SetLoggerProperties(logLevel, logFormat, logFile)
	},
	Run: func(cmd *cobra.Command, args []string) {
		runApplication()
	},
}

func runApplication() {
	fx.New(bootstrap.Module, fx.NopLogger).Run()
	//fx.New(bootstrap.Module).Run()
}
