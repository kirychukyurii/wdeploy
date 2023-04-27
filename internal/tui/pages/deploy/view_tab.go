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
	"github.com/kirychukyurii/wdeploy/internal/tui/components/editor"
	"text/template"
)

type ReadmeMsg struct{}

// FileContentMsg is a message that contains the content of a file.
type FileContentMsg struct {
	content string
	ext     string
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
	path           string

	cfg    config.Config
	logger logger.Logger
	/*
		ref    RefMsg
		repo   git.GitRepo
	*/
}

// NewView creates a new config model.
func NewView(common common.Common, cfg config.Config, logger logger.Logger) *View {
	s := spinner.New()
	s.Spinner = spinner.Dot

	f := &View{
		common:     common,
		code:       code.New(common, "", ""),
		dialog:     dialog.New(common, "Are you sure want to deploy Webitel?", []string{"Deploy", "Cancel"}),
		spinner:    s,
		lineNumber: true,

		cfg:    cfg,
		logger: logger,
	}

	f.code.SetShowLineNumber(f.lineNumber)
	return f
}

// SetSize implements common.Component.
func (r *View) SetSize(width, height int) {
	r.common.SetSize(width, height)
	hm := r.common.Styles.Repo.Body.GetVerticalFrameSize() +
		r.common.Styles.Repo.Header.GetHeight() +
		r.common.Styles.Repo.Header.GetVerticalFrameSize() +
		r.common.Styles.StatusBar.GetHeight() +
		r.common.Styles.Dialog.Box.GetHeight() +
		r.common.Styles.Dialog.Box.GetVerticalFrameSize()

	r.logger.Zap.Debug(fmt.Sprintf("width=%d height=%d hm=%d", width, height, hm))

	r.dialog.SetSize(width, r.common.Styles.Dialog.Box.GetHeight())
	r.code.SetSize(width, height-hm)
}

// ShortHelp implements help.KeyMap.
func (r *View) ShortHelp() []key.Binding {
	b := []key.Binding{
		r.common.KeyMap.LeftRight,
		r.common.KeyMap.Select,
	}

	return b
}

// FullHelp implements help.KeyMap.
func (r *View) FullHelp() [][]key.Binding {
	k := r.common.KeyMap
	b := [][]key.Binding{
		{
			k.Left,
			k.Right,
		},
		{
			k.Select,
		},
	}

	return b
}

// Init implements tea.Model.
func (r *View) Init() tea.Cmd {
	var buf bytes.Buffer

	tpl, err := template.New("").Parse(tview.Tmpl)
	if err != nil {
		r.logger.Zap.Debug(err)
	}

	err = tpl.Execute(&buf, r.cfg)
	if err != nil {
		r.logger.Zap.Debug(err)
	}

	view := buf.String()

	/*
		for _, host := range r.cfg.Inventory.Inventory.Hosts {
			view += fmt.Sprintf("%s", host.AnsibleHost)
			view += fmt.Sprintf(" (%s)", strings.Join(host.WebitelServices, ", "))
			view += "\n"
		}

	*/

	r.code.GotoTop()
	return tea.Batch(
		r.dialog.Init(),
		r.code.SetContent(view, "markdown"),
	)
}

// Update implements tea.Model.
func (r *View) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	case dialog.SelectDialogButtonMsg:
		fmt.Println(msg)
		if msg == 0 {
			// cmds = append(cmds, r.spinner.Update(""))
		}
	case RepoMsg:
		r.repo = action.Action(msg)
		cmds = append(cmds, r.Init())
	}
	_, cmd := r.dialog.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return r, tea.Batch(cmds...)
}

// View implements tea.Model.
func (r *View) View() string {
	view := lipgloss.JoinVertical(lipgloss.Top,
		r.dialog.View(),
		"\n",
		r.code.View(),
	)

	return view
}

// StatusBarValue implements statusbar.StatusBar.
func (r *View) StatusBarValue() string {
	return ""
}

// StatusBarInfo implements statusbar.StatusBar.
func (r *View) StatusBarInfo() string {
	return fmt.Sprintf("☰ %.f%%", r.code.ScrollPercent()*100)
}

// StatusBarBranch implements statusbar.StatusBar.
func (r *View) StatusBarBranch() string {
	return fmt.Sprintf("v%s", r.cfg.WebitelVersion)
}

func (r *View) updateFileContent() tea.Msg {
	hostsConfig, err := r.cfg.GetHostsConfigContent()
	if err != nil {
		return nil
	}

	return FileContentMsg{content: hostsConfig, ext: "yml"}
}

// editConfig opens the editor.
func (r *View) editConfig() tea.Cmd {
	return tea.ExecProcess(editor.Cmd(r.cfg.VarsFile), func(err error) tea.Msg {
		return r.updateFileContent()
	})
}

func (r *View) initSpinner() tea.Cmd {
	return r.spinner.Tick
}

func (r *View) deployWebitel() tea.Cmd {
	r.initSpinner()

	return nil
}
