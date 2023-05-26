package deploy

import (
	"bufio"
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kirychukyurii/wdeploy/internal/config"
	"github.com/kirychukyurii/wdeploy/internal/lib/ansible"
	"github.com/kirychukyurii/wdeploy/internal/lib/logger"
	"github.com/kirychukyurii/wdeploy/internal/tui/common"
	"github.com/kirychukyurii/wdeploy/internal/tui/components/action"
	"github.com/kirychukyurii/wdeploy/internal/tui/components/dialog"
	"github.com/kirychukyurii/wdeploy/internal/tui/components/footer"
	"github.com/kirychukyurii/wdeploy/internal/tui/components/statusbar"
	"github.com/kirychukyurii/wdeploy/internal/tui/components/tabs"
	"io"
)

type tab int

const (
	viewTab tab = iota
	logTab

	lastTab
)

func (t tab) String() string {
	return []string{
		"Deploy",
		"Log",
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

// Deploy is a view for a git repository.
type Deploy struct {
	common       common.Common
	selectedRepo action.Action
	statusbar    *statusbar.StatusBar

	activeTab tab
	tabs      *tabs.Tabs
	panes     []common.Component
	sub       chan string

	cfg    config.Config
	logger logger.Logger
}

// New returns a new Repo.
func New(c common.Common, cfg config.Config, logger logger.Logger) *Deploy {
	sb := statusbar.New(c)
	ts := make([]string, lastTab)

	// Tabs must match the order of tab constants above.
	for i, t := range []tab{viewTab, logTab} {
		ts[i] = t.String()
	}
	tb := tabs.New(c, ts)

	view := NewView(c, cfg, logger)
	log := NewLog(c, cfg, logger)

	// Make sure the order matches the order of tab constants above.
	panes := []common.Component{
		view,
		log,
	}

	d := &Deploy{
		common:    c,
		statusbar: sb,
		tabs:      tb,
		panes:     panes,
		sub:       make(chan string),
		cfg:       cfg,
		logger:    logger,
	}
	return d
}

// SetSize implements common.Component.
func (d *Deploy) SetSize(width, height int) {
	d.common.SetSize(width, height)
	hm := d.common.Styles.Repo.Body.GetVerticalFrameSize() +
		d.common.Styles.Repo.Header.GetHeight() +
		d.common.Styles.Repo.Header.GetVerticalFrameSize() +
		d.common.Styles.StatusBar.GetHeight()
	d.tabs.SetSize(width, height-hm)
	d.statusbar.SetSize(width, height-hm)

	for _, p := range d.panes {
		p.SetSize(width, height-hm)
	}
}

func (d *Deploy) commonHelp() []key.Binding {
	b := make([]key.Binding, 0)
	back := d.common.KeyMap.Back
	back.SetHelp("esc", "back to menu")
	tab := d.common.KeyMap.Section
	tab.SetHelp("tab", "switch tab")
	b = append(b, back)
	b = append(b, tab)

	return b
}

// ShortHelp implements help.KeyMap.
func (d *Deploy) ShortHelp() []key.Binding {
	b := d.commonHelp()
	b = append(b, d.panes[d.activeTab].(help.KeyMap).ShortHelp()...)

	return b
}

// FullHelp implements help.KeyMap.
func (d *Deploy) FullHelp() [][]key.Binding {
	b := make([][]key.Binding, 0)
	b = append(b, d.commonHelp())
	b = append(b, d.panes[d.activeTab].(help.KeyMap).FullHelp()...)

	return b
}

// Init implements tea.Log.
func (d *Deploy) Init() tea.Cmd {
	return tea.Batch(
		d.tabs.Init(),
		d.statusbar.Init(),
		waitForActivity(d.sub),
	)
}

// Update implements tea.Model.
func (d *Deploy) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case dialog.SelectDialogButtonMsg:
		if msg == 0 {
			d.activeTab = logTab
			cmd := tabs.SelectTabCmd(int(logTab))
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			//r.tabs.Update()

			// TODO
			_ = d.deploy(d.sub)
		}
	case RepoMsg:
		d.activeTab = 0
		d.selectedRepo = action.Action(msg) //git.GitRepo(msg)
		cmds = append(cmds,
			d.tabs.Init(),
			d.updateStatusBarCmd,
			d.updateModels(msg),
		)

	case tabs.SelectTabMsg:
		d.activeTab = tab(msg)
		t, cmd := d.tabs.Update(msg)
		d.tabs = t.(*tabs.Tabs)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	case tabs.ActiveTabMsg:
		d.activeTab = tab(msg)
		cmds = append(cmds,
			d.updateStatusBarCmd,
		)
	case tea.KeyMsg, tea.MouseMsg:
		t, cmd := d.tabs.Update(msg)
		d.tabs = t.(*tabs.Tabs)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		cmds = append(cmds, d.updateStatusBarCmd)
		switch msg := msg.(type) {
		case tea.MouseMsg:
			switch msg.Type {
			case tea.MouseLeft:
				switch {
				case d.common.Zone.Get("repo-help").InBounds(msg):
					cmds = append(cmds, footer.ToggleFooterCmd)
				}
			case tea.MouseRight:
				switch {
				case d.common.Zone.Get("repo-main").InBounds(msg):
					cmds = append(cmds, backCmd)
				}
			}
		}
	// The Log bubble is the only bubble that uses a spinner, so this is fine
	// for now. We need to pass the TickMsg to the Log bubble when the Log is
	// loading but not the current selected tab so that the spinner works.
	case UpdateStatusBarMsg:
		cmds = append(cmds, d.updateStatusBarCmd)
	case tea.WindowSizeMsg:
		cmds = append(cmds, d.updateModels(msg))
	}
	s, cmd := d.statusbar.Update(msg)
	d.statusbar = s.(*statusbar.StatusBar)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	m, cmd := d.panes[d.activeTab].Update(msg)
	d.panes[d.activeTab] = m.(common.Component)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return d, tea.Batch(cmds...)
}

// View implements tea.Model.
func (d *Deploy) View() string {
	s := d.common.Styles.Repo.Base.Copy().
		Width(d.common.Width).
		Height(d.common.Height)
	repoBodyStyle := d.common.Styles.Repo.Body.Copy()
	hm := repoBodyStyle.GetVerticalFrameSize() +
		d.common.Styles.Repo.Header.GetHeight() +
		d.common.Styles.Repo.Header.GetVerticalFrameSize() +
		d.common.Styles.StatusBar.GetHeight() +
		d.common.Styles.Tabs.GetHeight() +
		d.common.Styles.Tabs.GetVerticalFrameSize()
	mainStyle := repoBodyStyle.
		Height(d.common.Height - hm)
	main := d.common.Zone.Mark(
		"repo-main",
		mainStyle.Render(d.panes[d.activeTab].View()),
	)
	view := lipgloss.JoinVertical(lipgloss.Top,
		d.headerView(),
		d.tabs.View(),
		main,
		d.statusbar.View(),
	)

	return s.Render(view)
}

func (d *Deploy) headerView() string {
	if d.selectedRepo == nil {
		return ""
	}
	truncate := lipgloss.NewStyle().MaxWidth(d.common.Width)
	name := d.common.Styles.Repo.HeaderName.Render(d.selectedRepo.Title())
	desc := d.selectedRepo.Description()
	if desc == "" {
		desc = name
		name = ""
	} else {
		desc = d.common.Styles.Repo.HeaderDesc.Render(desc)
	}
	urlStyle := d.common.Styles.URLStyle.Copy().
		Width(d.common.Width - lipgloss.Width(desc) - 1).
		Align(lipgloss.Right)
	url := d.selectedRepo.ID()

	url = common.TruncateString(url, d.common.Width-lipgloss.Width(desc)-1)
	url = d.common.Zone.Mark(
		fmt.Sprintf("%s-url", d.selectedRepo.ID()),
		urlStyle.Render(url),
	)
	style := d.common.Styles.Repo.Header.Copy().Width(d.common.Width)

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

func (d *Deploy) updateStatusBarCmd() tea.Msg {
	value := d.panes[d.activeTab].(statusbar.Model).StatusBarValue()
	info := d.panes[d.activeTab].(statusbar.Model).StatusBarInfo()
	branch := d.panes[d.activeTab].(statusbar.Model).StatusBarBranch()

	return statusbar.StatusBarMsg{
		Key:    d.selectedRepo.ID(),
		Value:  value,
		Info:   info,
		Branch: branch,
	}
}

func (d *Deploy) updateModels(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 0)
	for i, b := range d.panes {
		m, cmd := b.Update(msg)
		d.panes[i] = m.(common.Component)
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

func (d *Deploy) deploy(sub chan string) error {
	reader, writer := io.Pipe()

	go func() {
		r := bufio.NewReader(reader)

		for {
			line, err := readLine(r)
			if err != nil {
				if err != io.EOF {
					d.logger.Zap.Error(err)
				}

				break
			}

			d.logger.Zap.Info(line)
			sub <- line
		}
	}()

	go func() {
		executor := ansible.NewExecutor(d.cfg, d.logger, writer)

		_ = executor.RunPlaybook()
	}()

	return nil
}

func readLine(r *bufio.Reader) (string, error) {
	var line []byte
	for {
		l, more, err := r.ReadLine()
		if err != nil {
			return "", err
		}

		// Avoid the copy if the first call produced a full line.
		if line == nil && !more {
			return string(l), nil
		}

		line = append(line, l...)
		if !more {
			break
		}
	}

	return string(line), nil
}

// A command that waits for the activity on a channel.
func waitForActivity(sub chan string) tea.Cmd {
	return func() tea.Msg {
		return LogMsg{
			message: <-sub,
			sub:     sub,
		}
	}
}
