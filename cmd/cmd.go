package cmd

import (
	"errors"
	"github.com/kirychukyurii/wdeploy/cmd/man"
	"github.com/kirychukyurii/wdeploy/cmd/run"
	"github.com/kirychukyurii/wdeploy/internal/constants"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	Command.AddCommand(run.Command)
	Command.AddCommand(man.Command)
}

var Command = &cobra.Command{
	Use:          "wdeploy",
	Short:        "wdeploy - easily deploy Webitel for your instances",
	SilenceUsage: true,
	Long: `wdeploy is a application that allows you to easily deploy Webitel services on your own instances.
Just specify needed variables and hosts configuration in TUI and take a coffee, wdeploy will do the rest for you`,
	Version: constants.Version,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New(
				"requires at least one arg, " +
					"you can view the available parameters through `--help`",
			)
		}
		return nil
	},
	PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
	Run:               func(cmd *cobra.Command, args []string) {},
}

func Execute() {
	if err := Command.Execute(); err != nil {
		os.Exit(-1)
	}
}
