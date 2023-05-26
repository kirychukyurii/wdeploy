package vars

import (
	"fmt"
	"github.com/alecthomas/chroma/lexers"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kirychukyurii/wdeploy/internal/config"
	"github.com/kirychukyurii/wdeploy/internal/lib/file"
	"github.com/kirychukyurii/wdeploy/internal/lib/logger"
	"github.com/kirychukyurii/wdeploy/internal/tui/common"
	"github.com/kirychukyurii/wdeploy/internal/tui/components/action"
	"github.com/kirychukyurii/wdeploy/internal/tui/components/code"
	"github.com/kirychukyurii/wdeploy/internal/tui/components/editor"
)

var (
	lineNo = key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "toggle line numbers"),
	)
)

type ReadmeMsg struct{}

// FileContentMsg is a message that contains the content of a file.
type FileContentMsg struct {
	content string
	ext     string
}

// Config is the readme component page.
type Config struct {
	common         common.Common
	code           *code.Code
	repo           action.Action
	currentContent FileContentMsg
	lineNumber     bool

	cfg    config.Config
	logger logger.Logger
}

// NewConfig creates a new config model.
func NewConfig(common common.Common, cfg config.Config, logger logger.Logger) *Config {
	c := &Config{
		common:     common,
		code:       code.New(common, "", ""),
		lineNumber: true,

		cfg:    cfg,
		logger: logger,
	}

	c.code.SetShowLineNumber(c.lineNumber)
	return c
}

// SetSize implements common.Component.
func (c *Config) SetSize(width, height int) {
	c.common.SetSize(width, height)
	c.code.SetSize(width, height)
}

// ShortHelp implements help.KeyMap.
func (c *Config) ShortHelp() []key.Binding {
	copyKey := c.common.KeyMap.Copy
	copyKey.SetHelp("c", "copy content")
	b := []key.Binding{
		c.common.KeyMap.UpDown,
		c.common.KeyMap.BackItem,
		c.common.KeyMap.EditItem,
		copyKey,
	}
	lexer := lexers.Match(c.currentContent.ext)
	lang := ""
	if lexer != nil && lexer.Config() != nil {
		lang = lexer.Config().Name
	}
	if lang != "markdown" {
		b = append(b, lineNo)
	}

	return b
}

// FullHelp implements help.KeyMap.
func (c *Config) FullHelp() [][]key.Binding {
	b := make([][]key.Binding, 0)

	copyKey := c.common.KeyMap.Copy
	copyKey.SetHelp("c", "copy content")
	k := c.code.KeyMap
	b = append(b, []key.Binding{
		c.common.KeyMap.BackItem,
		c.common.KeyMap.EditItem,
	})
	b = append(b, [][]key.Binding{
		{
			k.PageDown,
			k.PageUp,
			k.HalfPageDown,
			k.HalfPageUp,
		},
	}...)
	lc := []key.Binding{
		k.Down,
		k.Up,
		copyKey,
	}
	lexer := lexers.Match(c.currentContent.ext)
	lang := ""
	if lexer != nil && lexer.Config() != nil {
		lang = lexer.Config().Name
	}
	if lang != "markdown" {
		lc = append(lc, lineNo)
	}
	b = append(b, lc)

	return b
}

// Init implements tea.Model.
func (c *Config) Init() tea.Cmd {
	varsConfig, err := file.ReadFileContent(c.cfg.ConfigFiles[config.VarsConfig])
	if err != nil {
		return nil
	}

	if err = c.cfg.ReadToStruct(config.VarsConfig); err != nil {
		return nil
	}

	c.code.GotoTop()
	return tea.Batch(
		c.code.SetContent(varsConfig, ".yml"),
	)
}

// Update implements tea.Model.
func (c *Config) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if c.currentContent.content != "" {
			m, cmd := c.code.Update(msg)
			c.code = m.(*code.Code)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, lineNo):
			c.lineNumber = !c.lineNumber
			c.code.SetShowLineNumber(c.lineNumber)
			cmds = append(cmds, c.code.SetContent(c.currentContent.content, c.currentContent.ext))
		case key.Matches(msg, c.common.KeyMap.EditItem):
			return c, c.editConfig()
		case key.Matches(msg, c.common.KeyMap.Select):
			return c, c.editConfig()
		}
	case FileContentMsg:
		c.currentContent = msg
		c.code.SetContent(msg.content, msg.ext)
		c.code.GotoTop()
		cmds = append(cmds, updateStatusBarCmd)
	case RepoMsg:
		c.repo = action.Action(msg)
		cmds = append(cmds, c.Init())

	}
	co, cmd := c.code.Update(msg)
	c.code = co.(*code.Code)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	return c, tea.Batch(cmds...)
}

// View implements tea.Model.
func (c *Config) View() string {
	return c.code.View()
}

// StatusBarValue implements statusbar.StatusBar.
func (c *Config) StatusBarValue() string {
	return c.cfg.ConfigFiles[config.VarsConfig]
}

// StatusBarInfo implements statusbar.StatusBar.
func (c *Config) StatusBarInfo() string {
	return fmt.Sprintf("☰ %.f%%", c.code.ScrollPercent()*100)
}

// StatusBarBranch implements statusbar.StatusBar.
func (c *Config) StatusBarBranch() string {
	return fmt.Sprintf("v%s", c.cfg.WebitelVersion)
}

func (c *Config) updateFileContent() tea.Msg {
	varsConfig, err := file.ReadFileContent(c.cfg.ConfigFiles[config.VarsConfig])
	if err != nil {
		return nil
	}

	if err = c.cfg.ReadToStruct(config.VarsConfig); err != nil {
		return nil
	}

	return FileContentMsg{content: varsConfig, ext: ".yml"}
}

// editConfig opens the editor.
func (c *Config) editConfig() tea.Cmd {
	return tea.ExecProcess(editor.Cmd(c.cfg.ConfigFiles[config.VarsConfig]), func(err error) tea.Msg {
		return c.updateFileContent()
	})
}
