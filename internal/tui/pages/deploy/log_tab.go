package deploy

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fsnotify/fsnotify"
	"github.com/kirychukyurii/wdeploy/internal/config"
	"github.com/kirychukyurii/wdeploy/internal/lib/logger"
	"github.com/kirychukyurii/wdeploy/internal/tui/common"
	"github.com/kirychukyurii/wdeploy/internal/tui/components/action"
	"github.com/kirychukyurii/wdeploy/internal/tui/components/code"
)

type LogMsg struct{}

// LogContentMsg is a message that contains the content of a file.
type LogContentMsg struct {
	content string
	ext     string
}

// Log is the readme component page.
type Log struct {
	common         common.Common
	code           *code.Code
	repo           action.Action
	spinner        spinner.Model
	currentContent LogContentMsg
	lineNumber     bool
	path           string

	sub chan struct{} // where we'll receive activity notifications

	cfg    config.Config
	logger logger.Logger
}

// NewLog creates a new config model.
func NewLog(common common.Common, cfg config.Config, logger logger.Logger) *Log {
	s := spinner.New()
	s.Spinner = spinner.Dot

	f := &Log{
		sub:        make(chan struct{}),
		common:     common,
		code:       code.New(common, "", ""),
		spinner:    s,
		lineNumber: true,

		cfg:    cfg,
		logger: logger,
	}

	f.code.SetShowLineNumber(false)
	return f
}

// SetSize implements common.Component.
func (r *Log) SetSize(width, height int) {
	r.common.SetSize(width, height)
	r.code.SetSize(width, height)
}

// ShortHelp implements help.KeyMap.
func (r *Log) ShortHelp() []key.Binding {
	b := []key.Binding{
		r.common.KeyMap.LeftRight,
		r.common.KeyMap.Select,
		r.common.KeyMap.UpDown,
		r.common.KeyMap.BackItem,
	}

	return b
}

// FullHelp implements help.KeyMap.
func (r *Log) FullHelp() [][]key.Binding {
	k := r.common.KeyMap
	b := [][]key.Binding{
		{
			k.Left,
			k.Right,
			k.Down,
			k.Up,
			r.code.KeyMap.PageDown,
			r.code.KeyMap.PageUp,
			r.code.KeyMap.HalfPageDown,
			r.code.KeyMap.HalfPageUp,
		},
		{
			k.Select,
		},
	}

	return b
}

// Init implements tea.Model.
func (r *Log) Init() tea.Cmd {
	r.logger.Zap.Debugf("(r *Log) Init()=%s", r.currentContent)
	r.currentContent.content, _ = r.cfg.GetAnsibleLogContent()

	r.code.GotoBottom()

	return tea.Batch(
		r.code.SetContent(r.currentContent.content, ".yml"),
		r.fileWatcher(r.sub),
		r.waitForActivity(r.sub), // wait for activity
	)
}

// Update implements tea.Model.
func (r *Log) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	r.logger.Zap.Debugf("Log.Update() msg.%T=%s", msg, msg)

	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if r.currentContent.content != "" {
			m, cmd := r.code.Update(msg)
			r.code = m.(*code.Code)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}

	case LogMsg:
		r.logger.Zap.Debug(fmt.Sprintf("Log.Update().LogContentMsg: %s", msg))
		r.currentContent.content, _ = r.cfg.GetAnsibleLogContent()
		r.code.SetContent(r.currentContent.content, ".yml")
		r.code.GotoBottom()
		cmds = append(cmds, r.waitForActivity(r.sub), updateStatusBarCmd)

	case RepoMsg:
		r.repo = action.Action(msg)
		cmd := r.Init()
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	c, cmd := r.code.Update(msg)
	r.code = c.(*code.Code)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return r, tea.Batch(cmds...)
}

// View implements tea.Model.
func (r *Log) View() string {
	view := lipgloss.JoinVertical(lipgloss.Top,
		r.code.View(),
	)

	return view
}

// StatusBarValue implements statusbar.StatusBar.
func (r *Log) StatusBarValue() string {
	return ""
}

// StatusBarInfo implements statusbar.StatusBar.
func (r *Log) StatusBarInfo() string {
	return fmt.Sprintf("☰ %.f%%", r.code.ScrollPercent()*100)
}

// StatusBarBranch implements statusbar.StatusBar.
func (r *Log) StatusBarBranch() string {
	return fmt.Sprintf("v%s", r.cfg.WebitelVersion)
}

func (r *Log) initSpinner() tea.Cmd {
	return r.spinner.Tick
}

func (r *Log) deployWebitel() tea.Cmd {
	r.initSpinner()

	return nil
}

func (r *Log) fileWatcher(sub chan struct{}) tea.Cmd {
	r.logger.Zap.Debug("setup watcher")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		r.logger.Zap.Error(err)
	}

	// provide the file name along with path to be watched
	err = watcher.Add(fmt.Sprintf("%s/ansible.log", r.cfg.LogDirectory))
	if err != nil {
		r.logger.Zap.Error(err)
	}

	// done := make(chan bool)
	// use goroutine to start the watcher
	return func() tea.Msg {
		defer func(watcher *fsnotify.Watcher) {
			r.logger.Zap.Debug("close watcher")
			err := watcher.Close()
			if err != nil {

			}
		}(watcher)

		for {
			select {
			case event := <-watcher.Events:
				// monitor only for write events
				if event.Op&fsnotify.Write == fsnotify.Write {
					sub <- struct{}{}

					r.logger.Zap.Debug(fmt.Sprintf("Modified file: %s", event.Name))
				}
			case err := <-watcher.Errors:
				if err != nil {
					r.logger.Zap.Error(err)
				}
			}
		}
	}

}

// A command that waits for the activity on a channel.
func (r *Log) waitForActivity(sub chan struct{}) tea.Cmd {
	return func() tea.Msg {
		return LogMsg(<-sub)
	}
}
