package bootstrap

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kirychukyurii/wdeploy/internal/config"
	"github.com/kirychukyurii/wdeploy/internal/lib"
	"github.com/kirychukyurii/wdeploy/internal/lib/logger"
	"github.com/kirychukyurii/wdeploy/internal/tui"
	"github.com/kirychukyurii/wdeploy/internal/tui/common"
	"github.com/kirychukyurii/wdeploy/internal/tui/keymap"
	"github.com/kirychukyurii/wdeploy/internal/tui/styles"
	zone "github.com/lrstanley/bubblezone"
	"go.uber.org/fx"
	"golang.org/x/term"
	"os"
)

var Module = fx.Options(
	config.Module,
	lib.Module,
	fx.Invoke(bootstrap),
)

func bootstrap(lifecycle fx.Lifecycle, logger logger.Logger, config config.Config) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logger.Zap.Info("Starting Application")
			//fmt.Println(config)
			if config.WebitelRepositoryUser == "" || config.WebitelRepositoryPassword == "" {
				fmt.Println("403")
				os.Exit(1)
			}

			go func() {
				logger.Zap.Debug("Started goroutine")

				var opts []tea.ProgramOption

				// Always append alt screen program option.
				opts = append(opts, tea.WithAltScreen(), tea.WithMouseCellMotion())

				// Initialize and start app.
				width, height, err := term.GetSize(0)
				logger.Zap.Debug(fmt.Sprintf("Initial terminal size: width=%d, height=%d", width, height))
				if err != nil {
					logger.Zap.Fatalf("Failed to get terminal size: %s", err.Error())
				}

				c := common.Common{
					Styles: styles.DefaultStyles(),
					KeyMap: keymap.DefaultKeyMap(),
					Width:  width,
					Height: height,
					Zone:   zone.New(),
				}

				initialModel := tui.New(c, config, logger)

				p := tea.NewProgram(initialModel, opts...)
				if _, err := p.Run(); err != nil {
					logger.Zap.Fatalf("Failed to start: %s", err.Error())
				}
			}()

			return nil
		},
		OnStop: func(context.Context) error {
			logger.Zap.Info("Stopping Application")

			return nil
		},
	})
}
