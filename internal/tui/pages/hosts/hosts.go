package hosts

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

type state int

const (
	loadingState state = iota
	loadedState
)

type tab int

const (
	configTab tab = iota
	//testTab
	lastTab
)

func (t tab) String() string {
	return []string{
		"Config",
		//"Test",
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
	//test := NewConfig(c, cfg, logger)

	// Make sure the order matches the order of tab constants above.
	panes := []common.Component{
		config,
		//	test,
	}

	r := &Inventory{
		common:    c,
		statusbar: sb,
		tabs:      tb,
		panes:     panes,
		cfg:       cfg,
		logger:    logger,
	}
	return r
}

// SetSize implements common.Component.
func (r *Inventory) SetSize(width, height int) {
	r.common.SetSize(width, height)
	hm := r.common.Styles.Repo.Body.GetVerticalFrameSize() +
		r.common.Styles.Repo.Header.GetHeight() +
		r.common.Styles.Repo.Header.GetVerticalFrameSize() +
		r.common.Styles.StatusBar.GetHeight()
	r.tabs.SetSize(width, height-hm)
	r.statusbar.SetSize(width, height-hm)
	for _, p := range r.panes {
		p.SetSize(width, height-hm)
	}

}

func (r *Inventory) commonHelp() []key.Binding {
	b := make([]key.Binding, 0)
	back := r.common.KeyMap.Back
	back.SetHelp("esc", "back to menu")
	tab := r.common.KeyMap.Section
	tab.SetHelp("tab", "switch tab")
	b = append(b, back)
	b = append(b, tab)
	return b
}

// ShortHelp implements help.KeyMap.
func (r *Inventory) ShortHelp() []key.Binding {
	b := r.commonHelp()
	b = append(b, r.panes[r.activeTab].(help.KeyMap).ShortHelp()...)
	return b
}

// FullHelp implements help.KeyMap.
func (r *Inventory) FullHelp() [][]key.Binding {
	b := make([][]key.Binding, 0)
	b = append(b, r.commonHelp())
	b = append(b, r.panes[r.activeTab].(help.KeyMap).FullHelp()...)
	return b
}

// Init implements tea.View.
func (r *Inventory) Init() tea.Cmd {
	return tea.Batch(
		r.tabs.Init(),
		r.statusbar.Init(),
	)
}

// Update implements tea.Model.
func (r *Inventory) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	r.logger.Zap.Debug(msg)
	switch msg := msg.(type) {
	case RepoMsg:
		r.activeTab = 0
		r.selectedRepo = action.Action(msg) //git.GitRepo(msg)
		cmds = append(cmds,
			r.tabs.Init(),
			r.updateStatusBarCmd,
			r.updateModels(msg),
		)
	case tabs.SelectTabMsg:
		r.activeTab = tab(msg)
		t, cmd := r.tabs.Update(msg)
		r.tabs = t.(*tabs.Tabs)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	case tabs.ActiveTabMsg:
		r.activeTab = tab(msg)
		cmds = append(cmds,
			r.updateStatusBarCmd,
		)
	case tea.KeyMsg, tea.MouseMsg:
		t, cmd := r.tabs.Update(msg)
		r.tabs = t.(*tabs.Tabs)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		cmds = append(cmds, r.updateStatusBarCmd)
		switch msg := msg.(type) {
		case tea.MouseMsg:
			switch msg.Type {
			case tea.MouseLeft:
				switch {
				case r.common.Zone.Get("repo-help").InBounds(msg):
					cmds = append(cmds, footer.ToggleFooterCmd)
				}
			case tea.MouseRight:
				switch {
				case r.common.Zone.Get("repo-main").InBounds(msg):
					cmds = append(cmds, backCmd)
				}
			}
		}
	// The Log bubble is the only bubble that uses a spinner, so this is fine
	// for now. We need to pass the TickMsg to the Log bubble when the Log is
	// loading but not the current selected tab so that the spinner works.
	case UpdateStatusBarMsg:
		cmds = append(cmds, r.updateStatusBarCmd)
	case tea.WindowSizeMsg:
		cmds = append(cmds, r.updateModels(msg))
	}
	s, cmd := r.statusbar.Update(msg)
	r.statusbar = s.(*statusbar.StatusBar)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	m, cmd := r.panes[r.activeTab].Update(msg)
	r.panes[r.activeTab] = m.(common.Component)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return r, tea.Batch(cmds...)
}

// View implements tea.Model.
func (r *Inventory) View() string {
	//fmt.Println("test")
	/*
		s := r.common.Styles.Repo.Base.Copy().
			Width(r.common.Width).
			Height(r.common.Height)
		view := lipgloss.JoinVertical(lipgloss.Top,
			r.headerView(),
			r.statusbar.View(),
		)
		return s.Render(view)
	*/

	s := r.common.Styles.Repo.Base.Copy().
		Width(r.common.Width).
		Height(r.common.Height)
	repoBodyStyle := r.common.Styles.Repo.Body.Copy()
	hm := repoBodyStyle.GetVerticalFrameSize() +
		r.common.Styles.Repo.Header.GetHeight() +
		r.common.Styles.Repo.Header.GetVerticalFrameSize() +
		r.common.Styles.StatusBar.GetHeight() +
		r.common.Styles.Tabs.GetHeight() +
		r.common.Styles.Tabs.GetVerticalFrameSize()
	mainStyle := repoBodyStyle.
		Height(r.common.Height - hm)
	main := r.common.Zone.Mark(
		"repo-main",
		mainStyle.Render(r.panes[r.activeTab].View()),
	)
	view := lipgloss.JoinVertical(lipgloss.Top,
		r.headerView(),
		r.tabs.View(),
		main,
		r.statusbar.View(),
	)
	return s.Render(view)
}

func (r *Inventory) headerView() string {
	/*
		if r.selectedRepo == nil {
			return ""
		}
		truncate := lipgloss.NewStyle().MaxWidth(r.common.Width)
		name := r.common.Styles.Repo.HeaderName.Render(r.selectedRepo.Title())
		desc := r.selectedRepo.Description()
		if desc == "" {
			desc = name
			name = ""
		} else {
			desc = r.common.Styles.Repo.HeaderDesc.Render(desc)
		}
		urlStyle := r.common.Styles.URLStyle.Copy().
			Width(r.common.Width - lipgloss.Width(desc) - 1).
			Align(lipgloss.Right)
		url := r.selectedRepo.ID()

		url = common.TruncateString(url, r.common.Width-lipgloss.Width(desc)-1)
		url = r.common.Zone.Mark(
			fmt.Sprintf("%s-url", r.selectedRepo.ID()),
			urlStyle.Render(url),
		)
		style := r.common.Styles.Repo.Header.Copy().Width(r.common.Width)
		return style.Render(
			lipgloss.JoinVertical(lipgloss.Top,
				truncate.Render(name),
				truncate.Render(lipgloss.JoinHorizontal(lipgloss.Left,
					desc,
					url,
				)),
			),
		)
	*/
	if r.selectedRepo == nil {
		return ""
	}
	truncate := lipgloss.NewStyle().MaxWidth(r.common.Width)
	name := r.common.Styles.Repo.HeaderName.Render(r.selectedRepo.Title())
	desc := r.selectedRepo.Description()
	if desc == "" {
		desc = name
		name = ""
	} else {
		desc = r.common.Styles.Repo.HeaderDesc.Render(desc)
	}
	urlStyle := r.common.Styles.URLStyle.Copy().
		Width(r.common.Width - lipgloss.Width(desc) - 1).
		Align(lipgloss.Right)
	url := r.selectedRepo.ID()

	url = common.TruncateString(url, r.common.Width-lipgloss.Width(desc)-1)
	url = r.common.Zone.Mark(
		fmt.Sprintf("%s-url", r.selectedRepo.ID()),
		urlStyle.Render(url),
	)
	style := r.common.Styles.Repo.Header.Copy().Width(r.common.Width)
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

func (r *Inventory) updateStatusBarCmd() tea.Msg {
	/*
		if r.selectedRepo == nil {
			return nil
		}
	*/

	value := r.panes[r.activeTab].(statusbar.Model).StatusBarValue()
	info := r.panes[r.activeTab].(statusbar.Model).StatusBarInfo()
	branch := r.panes[r.activeTab].(statusbar.Model).StatusBarBranch()

	return statusbar.StatusBarMsg{
		Key:    r.selectedRepo.ID(),
		Value:  value,
		Info:   info,
		Branch: branch,
	}
}

func (r *Inventory) updateModels(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 0)
	for i, b := range r.panes {
		m, cmd := b.Update(msg)
		r.panes[i] = m.(common.Component)
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
