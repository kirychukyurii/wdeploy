package vars

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kirychukyurii/wdeploy/internal/config"
	"github.com/kirychukyurii/wdeploy/internal/lib/logger"
	"github.com/kirychukyurii/wdeploy/internal/tui/common"
	"github.com/kirychukyurii/wdeploy/internal/tui/components/action"
	"github.com/kirychukyurii/wdeploy/internal/tui/components/footer"
	"github.com/kirychukyurii/wdeploy/internal/tui/components/statusbar"
	"github.com/kirychukyurii/wdeploy/internal/tui/components/tabs"
)

type tab int

const (
	configTab tab = iota

	lastTab
)

func (t tab) String() string {
	return []string{
		"Config",
	}[t]
}

// ResetURLMsg is a message to reset the URL string.
type ResetURLMsg struct{}

// UpdateStatusBarMsg updates the status bar.
type UpdateStatusBarMsg struct{}

// RepoMsg is a message that contains a git.Repository.
type RepoMsg action.Action

// BackMsg is a message to go back to the previous view.
type BackMsg struct{}

// Vars is a view for a git repository.
type Vars struct {
	common       common.Common
	selectedRepo action.Action
	statusbar    *statusbar.StatusBar

	activeTab tab
	tabs      *tabs.Tabs
	panes     []common.Component

	cfg    config.Config
	logger logger.Logger
}

// New returns a new Repo.
func New(c common.Common, cfg config.Config, logger logger.Logger) *Vars {
	sb := statusbar.New(c)
	ts := make([]string, lastTab)
	// Tabs must match the order of tab constants above.
	for i, t := range []tab{configTab} {
		ts[i] = t.String()
	}
	tb := tabs.New(c, ts)

	config := NewConfig(c, cfg, logger)

	// Make sure the order matches the order of tab constants above.
	panes := []common.Component{
		config,
	}

	v := &Vars{
		common:    c,
		statusbar: sb,
		tabs:      tb,
		panes:     panes,
		cfg:       cfg,
		logger:    logger,
	}
	return v
}

// SetSize implements common.Component.
func (v *Vars) SetSize(width, height int) {
	v.common.SetSize(width, height)
	hm := v.common.Styles.Repo.Body.GetVerticalFrameSize() +
		v.common.Styles.Repo.Header.GetHeight() +
		v.common.Styles.Repo.Header.GetVerticalFrameSize() +
		v.common.Styles.StatusBar.GetHeight()
	v.tabs.SetSize(width, height-hm)
	v.statusbar.SetSize(width, height-hm)
	for _, p := range v.panes {
		p.SetSize(width, height-hm)
	}
}

func (v *Vars) commonHelp() []key.Binding {
	b := make([]key.Binding, 0)
	back := v.common.KeyMap.Back
	back.SetHelp("esc", "back to menu")
	tab := v.common.KeyMap.Section
	tab.SetHelp("tab", "switch tab")
	b = append(b, back)
	b = append(b, tab)
	return b
}

// ShortHelp implements help.KeyMap.
func (v *Vars) ShortHelp() []key.Binding {
	b := v.commonHelp()
	b = append(b, v.panes[v.activeTab].(help.KeyMap).ShortHelp()...)
	return b
}

// FullHelp implements help.KeyMap.
func (v *Vars) FullHelp() [][]key.Binding {
	b := make([][]key.Binding, 0)
	b = append(b, v.commonHelp())
	b = append(b, v.panes[v.activeTab].(help.KeyMap).FullHelp()...)
	return b
}

// Init implements tea.View.
func (v *Vars) Init() tea.Cmd {
	return tea.Batch(
		v.tabs.Init(),
		v.statusbar.Init(),
	)
}

// Update implements tea.Model.
func (v *Vars) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case RepoMsg:
		v.activeTab = 0
		v.selectedRepo = action.Action(msg) //git.GitRepo(msg)
		cmds = append(cmds,
			v.tabs.Init(),
			v.updateStatusBarCmd,
			v.updateModels(msg),
		)
	case tabs.SelectTabMsg:
		v.activeTab = tab(msg)
		t, cmd := v.tabs.Update(msg)
		v.tabs = t.(*tabs.Tabs)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	case tabs.ActiveTabMsg:
		v.activeTab = tab(msg)
		cmds = append(cmds,
			v.updateStatusBarCmd,
		)
	case tea.KeyMsg, tea.MouseMsg:
		t, cmd := v.tabs.Update(msg)
		v.tabs = t.(*tabs.Tabs)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		cmds = append(cmds, v.updateStatusBarCmd)
		switch msg := msg.(type) {
		case tea.MouseMsg:
			switch msg.Type {
			case tea.MouseLeft:
				switch {
				case v.common.Zone.Get("repo-help").InBounds(msg):
					cmds = append(cmds, footer.ToggleFooterCmd)
				}
			case tea.MouseRight:
				switch {
				case v.common.Zone.Get("repo-main").InBounds(msg):
					cmds = append(cmds, backCmd)
				}
			}
		}
	// The Log bubble is the only bubble that uses a spinner, so this is fine
	// for now. We need to pass the TickMsg to the Log bubble when the Log is
	// loading but not the current selected tab so that the spinner works.
	case UpdateStatusBarMsg:
		cmds = append(cmds, v.updateStatusBarCmd)
	case tea.WindowSizeMsg:
		cmds = append(cmds, v.updateModels(msg))
	}
	s, cmd := v.statusbar.Update(msg)
	v.statusbar = s.(*statusbar.StatusBar)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	m, cmd := v.panes[v.activeTab].Update(msg)
	v.panes[v.activeTab] = m.(common.Component)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return v, tea.Batch(cmds...)
}

// View implements tea.Model.
func (v *Vars) View() string {
	s := v.common.Styles.Repo.Base.Copy().
		Width(v.common.Width).
		Height(v.common.Height)
	repoBodyStyle := v.common.Styles.Repo.Body.Copy()
	hm := repoBodyStyle.GetVerticalFrameSize() +
		v.common.Styles.Repo.Header.GetHeight() +
		v.common.Styles.Repo.Header.GetVerticalFrameSize() +
		v.common.Styles.StatusBar.GetHeight() +
		v.common.Styles.Tabs.GetHeight() +
		v.common.Styles.Tabs.GetVerticalFrameSize()
	mainStyle := repoBodyStyle.
		Height(v.common.Height - hm)
	main := v.common.Zone.Mark(
		"repo-main",
		mainStyle.Render(v.panes[v.activeTab].View()),
	)
	view := lipgloss.JoinVertical(lipgloss.Top,
		v.headerView(),
		v.tabs.View(),
		main,
		v.statusbar.View(),
	)

	return s.Render(view)
}

func (v *Vars) headerView() string {
	if v.selectedRepo == nil {
		return ""
	}
	truncate := lipgloss.NewStyle().MaxWidth(v.common.Width)
	name := v.common.Styles.Repo.HeaderName.Render(v.selectedRepo.Title())
	desc := v.selectedRepo.Description()
	if desc == "" {
		desc = name
		name = ""
	} else {
		desc = v.common.Styles.Repo.HeaderDesc.Render(desc)
	}
	urlStyle := v.common.Styles.URLStyle.Copy().
		Width(v.common.Width - lipgloss.Width(desc) - 1).
		Align(lipgloss.Right)
	url := v.selectedRepo.ID()

	url = common.TruncateString(url, v.common.Width-lipgloss.Width(desc)-1)
	url = v.common.Zone.Mark(
		fmt.Sprintf("%s-url", v.selectedRepo.ID()),
		urlStyle.Render(url),
	)
	style := v.common.Styles.Repo.Header.Copy().Width(v.common.Width)

	return style.Render(
		lipgloss.JoinVertical(lipgloss.Top,
			truncate.Render(name),
			truncate.Render(lipgloss.JoinHorizontal(lipgloss.Left,
				desc,
				url,
			)),
		),
	)
}

func (v *Vars) updateStatusBarCmd() tea.Msg {
	value := v.panes[v.activeTab].(statusbar.Model).StatusBarValue()
	info := v.panes[v.activeTab].(statusbar.Model).StatusBarInfo()
	branch := v.panes[v.activeTab].(statusbar.Model).StatusBarBranch()

	return statusbar.StatusBarMsg{
		Key:    v.selectedRepo.ID(),
		Value:  value,
		Info:   info,
		Branch: branch,
	}
}

func (v *Vars) updateModels(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 0)
	for i, b := range v.panes {
		m, cmd := b.Update(msg)
		v.panes[i] = m.(common.Component)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return tea.Batch(cmds...)
}

func updateStatusBarCmd() tea.Msg {
	return UpdateStatusBarMsg{}
}

func backCmd() tea.Msg {
	return BackMsg{}
}
