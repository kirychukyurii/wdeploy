package deploy

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kirychukyurii/wdeploy/internal/config"
	"github.com/kirychukyurii/wdeploy/internal/lib/file"
	"github.com/kirychukyurii/wdeploy/internal/lib/logger"
	"github.com/kirychukyurii/wdeploy/internal/tui/common"
	"github.com/kirychukyurii/wdeploy/internal/tui/components/action"
	"github.com/kirychukyurii/wdeploy/internal/tui/components/code"
)

type LogMsg struct {
	message string
	sub     chan string
}

// LogContentMsg is a message that contains the content of a file.
type LogContentMsg struct {
	content string
	// ext     string
}

// Log is the readme component page.
type Log struct {
	common         common.Common
	code           *code.Code
	repo           action.Action
	spinner        spinner.Model
	currentContent LogContentMsg
	lineNumber     bool
	// path           string

	sub chan string // where we'll receive activity notifications

	cfg    config.Config
	logger logger.Logger
}

// NewLog creates a new config model.
func NewLog(common common.Common, cfg config.Config, logger logger.Logger) *Log {
	s := spinner.New()
	s.Spinner = spinner.Dot

	l := &Log{
		sub:        make(chan string),
		common:     common,
		code:       code.New(common, "", ""),
		spinner:    s,
		lineNumber: true,

		cfg:    cfg,
		logger: logger,
	}

	l.code.SetShowLineNumber(false)
	return l
}

// SetSize implements common.Component.
func (l *Log) SetSize(width, height int) {
	l.common.SetSize(width, height)
	l.code.SetSize(width, height)
}

// ShortHelp implements help.KeyMap.
func (l *Log) ShortHelp() []key.Binding {
	b := []key.Binding{
		l.common.KeyMap.LeftRight,
		l.common.KeyMap.Select,
		l.common.KeyMap.UpDown,
		l.common.KeyMap.BackItem,
	}

	return b
}

// FullHelp implements help.KeyMap.
func (l *Log) FullHelp() [][]key.Binding {
	k := l.common.KeyMap
	b := [][]key.Binding{
		{
			k.Left,
			k.Right,
			k.Down,
			k.Up,
			l.code.KeyMap.PageDown,
			l.code.KeyMap.PageUp,
			l.code.KeyMap.HalfPageDown,
			l.code.KeyMap.HalfPageUp,
		},
		{
			k.Select,
		},
	}

	return b
}

// Init implements tea.Model.
func (l *Log) Init() tea.Cmd {
	l.currentContent.content, _ = file.ReadFileContent(l.cfg.GetAnsibleLogLocation())
	l.code.GotoBottom()

	return tea.Batch(
		l.code.SetContent(l.currentContent.content, code.PlainTextExt),
	)
}

// Update implements tea.Model.
func (l *Log) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if l.currentContent.content != "" {
			m, cmd := l.code.Update(msg)
			l.code = m.(*code.Code)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}

	case LogMsg:
		l.currentContent.content += fmt.Sprintln(msg.message)
		l.code.SetContent(l.currentContent.content, code.PlainTextExt)
		l.code.GotoBottom()
		cmds = append(cmds, waitForActivity(msg.sub))

	case RepoMsg:
		l.repo = action.Action(msg)
		cmd := l.Init()
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	c, cmd := l.code.Update(msg)
	l.code = c.(*code.Code)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return l, tea.Batch(cmds...)
}

// View implements tea.Model.
func (l *Log) View() string {
	view := lipgloss.JoinVertical(lipgloss.Top,
		l.code.View(),
	)

	return view
}

// StatusBarValue implements statusbar.StatusBar.
func (l *Log) StatusBarValue() string {
	return l.cfg.GetAnsibleLogLocation()
}

// StatusBarInfo implements statusbar.StatusBar.
func (l *Log) StatusBarInfo() string {
	return fmt.Sprintf("☰ %.f%%", l.code.ScrollPercent()*100)
}

// StatusBarBranch implements statusbar.StatusBar.
func (l *Log) StatusBarBranch() string {
	return fmt.Sprintf("v%s", l.cfg.WebitelVersion)
}
