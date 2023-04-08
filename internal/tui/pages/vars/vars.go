package vars

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kirychukyurii/wdeploy/internal/lib"
	"github.com/kirychukyurii/wdeploy/internal/tui/app"
	"github.com/kirychukyurii/wdeploy/internal/tui/common"
	"github.com/kirychukyurii/wdeploy/internal/tui/components/footer"
	"github.com/kirychukyurii/wdeploy/internal/tui/components/statusbar"
)

type state int

const (
	loadingState state = iota
	loadedState
)

type tab int

const (
	readmeTab tab = iota
	testTab
	lastTab
)

func (t tab) String() string {
	return []string{
		"Readme",
		"Test",
	}[t]
}

// ResetURLMsg is a message to reset the URL string.
type ResetURLMsg struct{}

// UpdateStatusBarMsg updates the status bar.
type UpdateStatusBarMsg struct{}

// RepoMsg is a message that contains a git.Repository.
type RepoMsg app.Action

// BackMsg is a message to go back to the previous view.
type BackMsg struct{}

// Repo is a view for a git repository.
type Repo struct {
	common       common.Common
	selectedRepo app.Action
	statusbar    *statusbar.StatusBar
	logger       lib.Logger
}

// New returns a new Repo.
func New(c common.Common, logger lib.Logger) *Repo {
	sb := statusbar.New(c)

	r := &Repo{
		common:    c,
		statusbar: sb,
		logger:    logger,
	}
	return r
}

// SetSize implements common.Component.
func (r *Repo) SetSize(width, height int) {
	r.common.SetSize(width, height)
	hm := r.common.Styles.Repo.Body.GetVerticalFrameSize() +
		r.common.Styles.Repo.Header.GetHeight() +
		r.common.Styles.Repo.Header.GetVerticalFrameSize() +
		r.common.Styles.StatusBar.GetHeight()

	r.statusbar.SetSize(width, height-hm)

}

func (r *Repo) commonHelp() []key.Binding {
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
func (r *Repo) ShortHelp() []key.Binding {
	b := r.commonHelp()
	return b
}

// FullHelp implements help.KeyMap.
func (r *Repo) FullHelp() [][]key.Binding {
	b := make([][]key.Binding, 0)
	b = append(b, r.commonHelp())
	return b
}

// Init implements tea.View.
func (r *Repo) Init() tea.Cmd {
	//fmt.Println("vars.go: Init()")
	return tea.Batch(
		r.statusbar.Init(),
	)
}

// Update implements tea.Model.
func (r *Repo) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case RepoMsg:
		r.selectedRepo = app.Action(msg) //git.GitRepo(msg)
		cmds = append(cmds,
			r.updateModels(msg),
		)

	case tea.KeyMsg, tea.MouseMsg:

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

	return r, tea.Batch(cmds...)
}

// View implements tea.Model.
func (r *Repo) View() string {
	//fmt.Println("test")

	s := r.common.Styles.Repo.Base.Copy().
		Width(r.common.Width).
		Height(r.common.Height)
	view := lipgloss.JoinVertical(lipgloss.Top,
		r.headerView(),
		r.statusbar.View(),
	)
	return s.Render(view)
}

func (r *Repo) headerView() string {
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

func (r *Repo) updateStatusBarCmd() tea.Msg {
	if r.selectedRepo == nil {
		return nil
	}
	ref := ""

	return statusbar.StatusBarMsg{
		Key:    r.selectedRepo.ID(),
		Branch: fmt.Sprintf("* %s", ref),
	}
}

func (r *Repo) updateModels(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 0)

	return tea.Batch(cmds...)
}

func backCmd() tea.Msg {
	return BackMsg{}
}
