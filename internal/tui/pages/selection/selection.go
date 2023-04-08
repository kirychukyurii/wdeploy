package selection

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kirychukyurii/wdeploy/internal/lib"
	"github.com/kirychukyurii/wdeploy/internal/tui/app"
	"github.com/kirychukyurii/wdeploy/internal/tui/common"
	"github.com/kirychukyurii/wdeploy/internal/tui/components/selector"
	"github.com/kirychukyurii/wdeploy/internal/tui/components/tabs"
)

type pane int

const (
	selectorPane pane = iota
	readmePane
	lastPane
)

func (p pane) String() string {
	return []string{
		"Actions",
		"About",
	}[p]
}

// Selection is the model for the selection screen/page.
type Selection struct {
	common common.Common

	selector   *selector.Selector
	activePane pane
	tabs       *tabs.Tabs
	logger     lib.Logger
}

// New creates a new selection model.
func New(common common.Common, logger lib.Logger) *Selection {
	ts := make([]string, lastPane)
	for i, b := range []pane{selectorPane, readmePane} {
		ts[i] = b.String()
	}
	t := tabs.New(common, ts)
	t.TabSeparator = lipgloss.NewStyle()
	t.TabInactive = common.Styles.TopLevelNormalTab.Copy()
	t.TabActive = common.Styles.TopLevelActiveTab.Copy()
	t.TabDot = common.Styles.TopLevelActiveTabDot.Copy()
	t.UseDot = true
	sel := &Selection{
		common:     common,
		activePane: selectorPane, // start with the selector focused
		tabs:       t,
		logger:     logger,
	}

	selector := selector.New(common,
		[]selector.IdentifiableItem{},
		ItemDelegate{&common, &sel.activePane}, logger)
	selector.SetShowTitle(false)
	selector.SetShowHelp(false)
	selector.SetShowStatusBar(false)
	selector.DisableQuitKeybindings()
	sel.selector = selector
	return sel
}

func (s *Selection) getMargins() (wm, hm int) {
	wm = 0
	hm = s.common.Styles.Tabs.GetVerticalFrameSize() +
		s.common.Styles.Tabs.GetHeight()
	if s.activePane == selectorPane && s.IsFiltering() {
		// hide tabs when filtering
		hm = 0
	}
	return
}

// FilterState returns the current filter state.
func (s *Selection) FilterState() list.FilterState {
	return s.selector.FilterState()
}

// SetSize implements common.Component.
func (s *Selection) SetSize(width, height int) {
	s.common.SetSize(width, height)
	wm, hm := s.getMargins()
	s.tabs.SetSize(width, height-hm)
	s.selector.SetSize(width-wm, height-hm)
	// s.readme.SetSize(width-wm, height-hm-1) // -1 for readme status line
}

// IsFiltering returns true if the selector is currently filtering.
func (s *Selection) IsFiltering() bool {
	return s.FilterState() == list.Filtering
}

// ShortHelp implements help.KeyMap.
func (s *Selection) ShortHelp() []key.Binding {
	k := s.selector.KeyMap
	kb := make([]key.Binding, 0)
	kb = append(kb,
		s.common.KeyMap.UpDown,
		s.common.KeyMap.Section,
	)
	if s.activePane == selectorPane {
		copyKey := s.common.KeyMap.Copy
		copyKey.SetHelp("c", "copy command")
		kb = append(kb,
			s.common.KeyMap.Select,
			k.Filter,
			k.ClearFilter,
			copyKey,
		)
	}
	return kb
}

// FullHelp implements help.KeyMap.
func (s *Selection) FullHelp() [][]key.Binding {
	b := [][]key.Binding{
		{
			s.common.KeyMap.Section,
		},
	}
	switch s.activePane {
	case selectorPane:
		copyKey := s.common.KeyMap.Copy
		copyKey.SetHelp("c", "copy command")
		k := s.selector.KeyMap
		if !s.IsFiltering() {
			b[0] = append(b[0],
				s.common.KeyMap.Select,
				copyKey,
			)
		}
		b = append(b, []key.Binding{
			k.CursorUp,
			k.CursorDown,
		})
		b = append(b, []key.Binding{
			k.NextPage,
			k.PrevPage,
			k.GoToStart,
			k.GoToEnd,
		})
		b = append(b, []key.Binding{
			k.Filter,
			k.ClearFilter,
			k.CancelWhileFiltering,
			k.AcceptWhileFiltering,
		})
	}
	return b
}

// Init implements tea.Model.
func (s *Selection) Init() tea.Cmd {
	items := make([]selector.IdentifiableItem, 0)

	//actions := app.ActionItems.New()

	actions := app.ActionItems{
		app.ActionItem{
			Command: "vars",
			Name:    "Variables",
			Action:  "Variables are needed for setup version of Webitel services or whether install Grafana Dashboards, Fail2ban, LetsEncrypt certificate",
		},

		app.ActionItem{
			Command: "hosts",
			Name:    "Hosts credentials",
			Action:  "For deploying Webitel services you need specify server(s) credentials for connection",
		},
		app.ActionItem{
			Command: "deploy",
			Name:    "Deploy Webitel",
			Action:  "You are one step closer to deploy Webitel services! Choose this and go on",
		},
	}

	for _, a := range actions {
		items = append(items, Item{
			action: a,
			cmd:    a.Command,
		})
	}

	return tea.Batch(
		s.selector.Init(),
		s.selector.SetItems(items),
	)
}

// Update implements tea.Model.
func (s *Selection) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m, cmd := s.selector.Update(msg)
		s.selector = m.(*selector.Selector)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	case tea.KeyMsg, tea.MouseMsg:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, s.common.KeyMap.Back):
				cmds = append(cmds, s.selector.Init())
			}
		}
		t, cmd := s.tabs.Update(msg)
		s.tabs = t.(*tabs.Tabs)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	case tabs.ActiveTabMsg:
		s.activePane = pane(msg)
	}
	switch s.activePane {
	case selectorPane:
		m, cmd := s.selector.Update(msg)
		s.selector = m.(*selector.Selector)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return s, tea.Batch(cmds...)
}

// View implements tea.Model.
func (s *Selection) View() string {
	var view string
	wm, hm := s.getMargins()
	switch s.activePane {
	case selectorPane:
		ss := lipgloss.NewStyle().
			Width(s.common.Width - wm).
			Height(s.common.Height - hm)
		view = ss.Render(s.selector.View())
	}
	if s.activePane != selectorPane || s.FilterState() != list.Filtering {
		tabs := s.common.Styles.Tabs.Render(s.tabs.View())
		view = lipgloss.JoinVertical(lipgloss.Left,
			tabs,
			view,
		)
	}
	return lipgloss.JoinVertical(
		lipgloss.Left,
		view,
	)
}
