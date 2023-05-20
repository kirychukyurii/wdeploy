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
		"debug", "log output level: debug, info, warn, error, dpanic, panic, fatal")
	pf.StringVarP(&config.DefaultConfig.LogFormat, "log-format", "F",
		"plain", "log output format: json, console")
	pf.StringVarP(&config.DefaultConfig.LogDirectory, "log-path", "L",
		"./", "log output to this directory")
	pf.StringVarP(&config.DefaultConfig.ConfigFiles[config.VarsConfig], "vars", "V",
		"", "specify Ansible variables file")
	pf.StringVarP(&config.DefaultConfig.ConfigFiles[config.InventoryConfig], "inventory", "i",
		"", "specify Ansible inventory host path")
	pf.StringVarP(&config.DefaultConfig.WebitelRepositoryUser, "user", "u",
		"", "specify Webitel Repository user")
	pf.StringVarP(&config.DefaultConfig.WebitelRepositoryPassword, "password", "p",
		"", "specify Webitel Repository password")
	pf.StringVarP(&config.DefaultConfig.InventoryType, "deploy-type", "t",
		"localhost", "specify Ansible inventory template type: localhost, custom")
}

var Command = &cobra.Command{
	Use:          "run",
	Short:        "Run wdeploy TUI",
	Example:      `wdeploy run --user "testUser" --password "testPassword" --deploy-type custom`,
	SilenceUsage: true,
	PreRun:       func(cmd *cobra.Command, args []string) {},
	Run: func(cmd *cobra.Command, args []string) {
		runApplication()
	},
}

func runApplication() {
	fx.New(bootstrap.Module, fx.NopLogger).Run()
	//fx.New(bootstrap.Module).Run()
}
