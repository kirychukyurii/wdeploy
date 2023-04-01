package man

import (
	"fmt"
	"github.com/kirychukyurii/wdeploy/cmd/run"
	mcobra "github.com/muesli/mango-cobra"
	"github.com/muesli/roff"
	"github.com/spf13/cobra"
)

var (
	Command = &cobra.Command{
		Use:    "man",
		Short:  "Generate man pages",
		Args:   cobra.NoArgs,
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			manPage, err := mcobra.NewManPage(1, run.Command) //.
			if err != nil {
				return err
			}

			manPage = manPage.WithSection("Copyright", "(C) 2023 Webitel")
			fmt.Println(manPage.Build(roff.NewDocument()))
			return nil
		},
	}
)
