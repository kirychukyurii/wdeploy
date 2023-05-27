package deploy

import (
	"bytes"
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kirychukyurii/wdeploy/internal/config"
	"github.com/kirychukyurii/wdeploy/internal/lib/logger"
	tview "github.com/kirychukyurii/wdeploy/internal/templates/view"
	"github.com/kirychukyurii/wdeploy/internal/tui/common"
	"github.com/kirychukyurii/wdeploy/internal/tui/components/action"
	"github.com/kirychukyurii/wdeploy/internal/tui/components/code"
	"github.com/kirychukyurii/wdeploy/internal/tui/components/dialog"
	"text/template"
)

type ReadmeMsg struct{}

// FileContentMsg is a message that contains the content of a file.
type FileContentMsg struct {
	content string
}

// View is the readme component page.
type View struct {
	common         common.Common
	code           *code.Code
	dialog         *dialog.Dialog
	repo           action.Action
	spinner        spinner.Model
	currentContent FileContentMsg
	lineNumber     bool

	cfg    config.Config
	logger logger.Logger
}

// NewView creates a new config model.
func NewView(common common.Common, cfg config.Config, logger logger.Logger) *View {
	s := spinner.New()
	s.Spinner = spinner.Dot

	v := &View{
		common:     common,
		code:       code.New(common, "", ""),
		dialog:     dialog.New(common, "Are you sure want to deploy Webitel?", []string{"Deploy", "Cancel"}),
		spinner:    s,
		lineNumber: true,

		cfg:    cfg,
		logger: logger,
	}

	v.code.SetShowLineNumber(v.lineNumber)
	return v
}

// SetSize implements common.Component.
func (v *View) SetSize(width, height int) {
	v.common.SetSize(width, height)
	hm := v.common.Styles.Dialog.Box.GetHorizontalFrameSize()
	v.code.SetSize(width, height-hm-3)
	v.dialog.SetSize(width, hm)
}

// ShortHelp implements help.KeyMap.
func (v *View) ShortHelp() []key.Binding {
	b := []key.Binding{
		v.common.KeyMap.LeftRight,
		v.common.KeyMap.Select,
		v.common.KeyMap.UpDown,
		v.common.KeyMap.BackItem,
	}

	return b
}

// FullHelp implements help.KeyMap.
func (v *View) FullHelp() [][]key.Binding {
	k := v.common.KeyMap
	b := [][]key.Binding{
		{
			k.Left,
			k.Right,
			k.Down,
			k.Up,
			v.code.KeyMap.PageDown,
			v.code.KeyMap.PageUp,
			v.code.KeyMap.HalfPageDown,
			v.code.KeyMap.HalfPageUp,
		},
		{
			k.Select,
		},
	}

	return b
}

// Init implements tea.Model.
func (v *View) Init() tea.Cmd {
	var buf bytes.Buffer

	tpl, err := template.New("").Parse(tview.Tmpl)
	if err != nil {
		v.logger.Zap.Debug(err)
	}

	err = tpl.Execute(&buf, v.cfg)
	if err != nil {
		v.logger.Zap.Debug(err)
	}

	view := buf.String()

	v.code.GotoTop()
	return tea.Batch(
		v.dialog.Init(),
		v.code.SetContent(view, ".md"),
	)
}

// Update implements tea.Model.
func (v *View) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if v.currentContent.content != "" {
			m, cmd := v.code.Update(msg)
			v.code = m.(*code.Code)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	case RepoMsg:
		v.repo = action.Action(msg)
		cmds = append(cmds, v.Init())
	}
	d, cmd := v.dialog.Update(msg)
	v.dialog = d.(*dialog.Dialog)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	c, cmd := v.code.Update(msg)
	v.code = c.(*code.Code)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return v, tea.Batch(cmds...)
}

// View implements tea.Model.
func (v *View) View() string {
	view := lipgloss.JoinVertical(lipgloss.Top,
		v.code.View(),
		v.dialog.View(),
	)

	return view
}

// StatusBarValue implements statusbar.StatusBar.
func (v *View) StatusBarValue() string {
	return ""
}

// StatusBarInfo implements statusbar.StatusBar.
func (v *View) StatusBarInfo() string {
	return fmt.Sprintf("☰ %.f%%", v.code.ScrollPercent()*100)
}

// StatusBarBranch implements statusbar.StatusBar.
func (v *View) StatusBarBranch() string {
	return fmt.Sprintf("v%s", v.cfg.WebitelVersion)
}
