package hosts

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
	path           string

	cfg    config.Config
	logger logger.Logger
	/*
		ref    RefMsg
		repo   git.GitRepo
	*/
}

// NewConfig creates a new config model.
func NewConfig(common common.Common, cfg config.Config, logger logger.Logger) *Config {
	f := &Config{
		common:     common,
		code:       code.New(common, "", ""),
		lineNumber: true,

		cfg:    cfg,
		logger: logger,
	}

	f.code.SetShowLineNumber(f.lineNumber)
	return f
}

// SetSize implements common.Component.
func (r *Config) SetSize(width, height int) {
	r.common.SetSize(width, height)
	r.code.SetSize(width, height)
}

// ShortHelp implements help.KeyMap.
func (r *Config) ShortHelp() []key.Binding {
	copyKey := r.common.KeyMap.Copy
	copyKey.SetHelp("c", "copy content")
	b := []key.Binding{
		r.common.KeyMap.UpDown,
		r.common.KeyMap.BackItem,
		r.common.KeyMap.EditItem,
		copyKey,
	}
	lexer := lexers.Match(r.currentContent.ext)
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
func (r *Config) FullHelp() [][]key.Binding {
	b := make([][]key.Binding, 0)

	copyKey := r.common.KeyMap.Copy
	copyKey.SetHelp("c", "copy content")
	k := r.code.KeyMap
	b = append(b, []key.Binding{
		r.common.KeyMap.BackItem,
		r.common.KeyMap.EditItem,
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
	lexer := lexers.Match(r.currentContent.ext)
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
func (r *Config) Init() tea.Cmd {
	inventoryConfig, err := file.ReadFileContent(r.cfg.ConfigFiles[config.InventoryConfig])
	if err != nil {
		return nil
	}

	if err = r.cfg.ReadToStruct(config.InventoryConfig); err != nil {
		return nil
	}

	r.code.GotoTop()
	return tea.Batch(
		r.code.SetContent(inventoryConfig, "yml"),
	)
}

// Update implements tea.Model.
func (r *Config) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, lineNo):
			r.lineNumber = !r.lineNumber
			r.code.SetShowLineNumber(r.lineNumber)
			cmds = append(cmds, r.code.SetContent(r.currentContent.content, r.currentContent.ext))
		case key.Matches(msg, r.common.KeyMap.EditItem):
			return r, r.editConfig()
		case key.Matches(msg, r.common.KeyMap.Select):
			return r, r.editConfig()
		}
	case FileContentMsg:
		r.currentContent = msg
		r.code.SetContent(msg.content, msg.ext)
		r.code.GotoTop()
		cmds = append(cmds, updateStatusBarCmd)
	case RepoMsg:
		r.repo = action.Action(msg)
		cmds = append(cmds, r.Init())

	}
	c, cmd := r.code.Update(msg)
	r.code = c.(*code.Code)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	return r, tea.Batch(cmds...)
}

// View implements tea.Model.
func (r *Config) View() string {
	return r.code.View()
}

// StatusBarValue implements statusbar.StatusBar.
func (r *Config) StatusBarValue() string {
	p := r.cfg.ConfigFiles[config.InventoryConfig]
	if p == "." {
		return ""
	}
	return p
}

// StatusBarInfo implements statusbar.StatusBar.
func (r *Config) StatusBarInfo() string {
	return fmt.Sprintf("☰ %.f%%", r.code.ScrollPercent()*100)
}

// StatusBarBranch implements statusbar.StatusBar.
func (r *Config) StatusBarBranch() string {
	return fmt.Sprintf("v%s", r.cfg.WebitelVersion)
}

func (r *Config) updateFileContent() tea.Msg {
	hostsConfig, err := file.ReadFileContent(r.cfg.ConfigFiles[config.InventoryConfig])
	if err != nil {
		return nil
	}

	if err = r.cfg.ReadToStruct(config.InventoryConfig); err != nil {
		return nil
	}

	return FileContentMsg{content: hostsConfig, ext: "yml"}
}

// editConfig opens the editor.
func (r *Config) editConfig() tea.Cmd {
	return tea.ExecProcess(editor.Cmd(r.cfg.ConfigFiles[config.InventoryConfig]), func(err error) tea.Msg {
		return r.updateFileContent()
	})
}
