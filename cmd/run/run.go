package run

import (
	"github.com/kirychukyurii/wdeploy/bootstrap"
	"github.com/kirychukyurii/wdeploy/internal/config"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

func init() {
	pf := Command.PersistentFlags()
	pf.StringVarP(&config.DefaultConfig.LogLevel, "log-level", "l",
		"debug", "log level: debug, info, warn, error, dpanic, panic, fatal")
	pf.StringVarP(&config.DefaultConfig.LogFormat, "log-format", "F",
		"plain", "log format output: json, console")
	pf.StringVarP(&config.DefaultConfig.LogDirectory, "log-path", "L",
		"./", "log file location")
	pf.StringVarP(&config.DefaultConfig.VarsFile, "vars", "V",
		"", "variables file")
	pf.StringVarP(&config.DefaultConfig.HostsFile, "inventory", "i",
		"", "hosts file")
	pf.StringVarP(&config.DefaultConfig.WebitelRepositoryUser, "user", "u",
		"", "webitel repository user")
	pf.StringVarP(&config.DefaultConfig.WebitelRepositoryPassword, "password", "p",
		"", "webitel repository password")
	pf.StringVarP(&config.DefaultConfig.InventoryType, "type", "t",
		"localhost", "inventory template type: localhost, custom")
}

var Command = &cobra.Command{
	Use:          "run",
	Short:        "Start Ansible Playbook",
	Example:      "wdeploy run -V inventories/production/group_vars/all.yml -i inventories/production/inventory.yml",
	SilenceUsage: true,
	PreRun: func(cmd *cobra.Command, args []string) {

	},
	Run: func(cmd *cobra.Command, args []string) {
		runApplication()
	},
}

func runApplication() {
	fx.New(bootstrap.Module, fx.NopLogger).Run()
	//fx.New(bootstrap.Module).Run()
}
