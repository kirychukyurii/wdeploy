package inventory

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

// Inventory is a view for a git repository.
type Inventory struct {
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
func New(c common.Common, cfg config.Config, logger logger.Logger) *Inventory {
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

	i := &Inventory{
		common:    c,
		statusbar: sb,
		tabs:      tb,
		panes:     panes,
		cfg:       cfg,
		logger:    logger,
	}

	return i
}

// SetSize implements common.Component.
func (i *Inventory) SetSize(width, height int) {
	i.common.SetSize(width, height)
	hm := i.common.Styles.Repo.Body.GetVerticalFrameSize() +
		i.common.Styles.Repo.Header.GetHeight() +
		i.common.Styles.Repo.Header.GetVerticalFrameSize() +
		i.common.Styles.StatusBar.GetHeight()
	i.tabs.SetSize(width, height-hm)
	i.statusbar.SetSize(width, height-hm)
	for _, p := range i.panes {
		p.SetSize(width, height-hm)
	}

}

func (i *Inventory) commonHelp() []key.Binding {
	b := make([]key.Binding, 0)
	back := i.common.KeyMap.Back
	back.SetHelp("esc", "back to menu")
	tab := i.common.KeyMap.Section
	tab.SetHelp("tab", "switch tab")
	b = append(b, back)
	b = append(b, tab)
	return b
}

// ShortHelp implements help.KeyMap.
func (i *Inventory) ShortHelp() []key.Binding {
	b := i.commonHelp()
	b = append(b, i.panes[i.activeTab].(help.KeyMap).ShortHelp()...)
	return b
}

// FullHelp implements help.KeyMap.
func (i *Inventory) FullHelp() [][]key.Binding {
	b := make([][]key.Binding, 0)
	b = append(b, i.commonHelp())
	b = append(b, i.panes[i.activeTab].(help.KeyMap).FullHelp()...)
	return b
}

// Init implements tea.View.
func (i *Inventory) Init() tea.Cmd {
	return tea.Batch(
		i.tabs.Init(),
		i.statusbar.Init(),
	)
}

// Update implements tea.Model.
func (i *Inventory) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case RepoMsg:
		i.activeTab = 0
		i.selectedRepo = action.Action(msg) //git.GitRepo(msg)
		cmds = append(cmds,
			i.tabs.Init(),
			i.updateStatusBarCmd,
			i.updateModels(msg),
		)
	case tabs.SelectTabMsg:
		i.activeTab = tab(msg)
		t, cmd := i.tabs.Update(msg)
		i.tabs = t.(*tabs.Tabs)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	case tabs.ActiveTabMsg:
		i.activeTab = tab(msg)
		cmds = append(cmds,
			i.updateStatusBarCmd,
		)
	case tea.KeyMsg, tea.MouseMsg:
		t, cmd := i.tabs.Update(msg)
		i.tabs = t.(*tabs.Tabs)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		cmds = append(cmds, i.updateStatusBarCmd)
		switch msg := msg.(type) {
		case tea.MouseMsg:
			switch msg.Type {
			case tea.MouseLeft:
				switch {
				case i.common.Zone.Get("repo-help").InBounds(msg):
					cmds = append(cmds, footer.ToggleFooterCmd)
				}
			case tea.MouseRight:
				switch {
				case i.common.Zone.Get("repo-main").InBounds(msg):
					cmds = append(cmds, backCmd)
				}
			}
		}
	// The Log bubble is the only bubble that uses a spinner, so this is fine
	// for now. We need to pass the TickMsg to the Log bubble when the Log is
	// loading but not the current selected tab so that the spinner works.
	case UpdateStatusBarMsg:
		cmds = append(cmds, i.updateStatusBarCmd)
	case tea.WindowSizeMsg:
		cmds = append(cmds, i.updateModels(msg))
	}
	s, cmd := i.statusbar.Update(msg)
	i.statusbar = s.(*statusbar.StatusBar)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	m, cmd := i.panes[i.activeTab].Update(msg)
	i.panes[i.activeTab] = m.(common.Component)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return i, tea.Batch(cmds...)
}

// View implements tea.Model.
func (i *Inventory) View() string {
	s := i.common.Styles.Repo.Base.Copy().
		Width(i.common.Width).
		Height(i.common.Height)
	repoBodyStyle := i.common.Styles.Repo.Body.Copy()
	hm := repoBodyStyle.GetVerticalFrameSize() +
		i.common.Styles.Repo.Header.GetHeight() +
		i.common.Styles.Repo.Header.GetVerticalFrameSize() +
		i.common.Styles.StatusBar.GetHeight() +
		i.common.Styles.Tabs.GetHeight() +
		i.common.Styles.Tabs.GetVerticalFrameSize()
	mainStyle := repoBodyStyle.
		Height(i.common.Height - hm)
	main := i.common.Zone.Mark(
		"repo-main",
		mainStyle.Render(i.panes[i.activeTab].View()),
	)
	view := lipgloss.JoinVertical(lipgloss.Top,
		i.headerView(),
		i.tabs.View(),
		main,
		i.statusbar.View(),
	)

	return s.Render(view)
}

func (i *Inventory) headerView() string {
	if i.selectedRepo == nil {
		return ""
	}
	truncate := lipgloss.NewStyle().MaxWidth(i.common.Width)
	name := i.common.Styles.Repo.HeaderName.Render(i.selectedRepo.Title())
	desc := i.selectedRepo.Description()
	if desc == "" {
		desc = name
		name = ""
	} else {
		desc = i.common.Styles.Repo.HeaderDesc.Render(desc)
	}
	urlStyle := i.common.Styles.URLStyle.Copy().
		Width(i.common.Width - lipgloss.Width(desc) - 1).
		Align(lipgloss.Right)
	url := i.selectedRepo.ID()

	url = common.TruncateString(url, i.common.Width-lipgloss.Width(desc)-1)
	url = i.common.Zone.Mark(
		fmt.Sprintf("%s-url", i.selectedRepo.ID()),
		urlStyle.Render(url),
	)
	style := i.common.Styles.Repo.Header.Copy().Width(i.common.Width)

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

func (i *Inventory) updateStatusBarCmd() tea.Msg {
	/*
		if r.selectedRepo == nil {
			return nil
		}
	*/

	value := i.panes[i.activeTab].(statusbar.Model).StatusBarValue()
	info := i.panes[i.activeTab].(statusbar.Model).StatusBarInfo()
	branch := i.panes[i.activeTab].(statusbar.Model).StatusBarBranch()

	return statusbar.StatusBarMsg{
		Key:    i.selectedRepo.ID(),
		Value:  value,
		Info:   info,
		Branch: branch,
	}
}

func (i *Inventory) updateModels(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 0)
	for n, b := range i.panes {
		m, cmd := b.Update(msg)
		i.panes[n] = m.(common.Component)
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
