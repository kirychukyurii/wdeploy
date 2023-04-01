package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kirychukyurii/wdeploy/internal/tui/common"
	"github.com/kirychukyurii/wdeploy/internal/tui/components/footer"
	"github.com/kirychukyurii/wdeploy/internal/tui/components/header"
	"github.com/kirychukyurii/wdeploy/internal/tui/pages/selection"
)

type page int

const (
	selectionPage page = iota
	varsPage
	hostsPage
	deployPage
)

type sessionState int

const (
	startState sessionState = iota
	errorState
	loadedState
)

// UI is the main UI model.
type UI struct {
	common     common.Common
	pages      []common.Component
	activePage page
	state      sessionState
	header     *header.Header
	footer     *footer.Footer
	showFooter bool
	error      error
}

// New returns a new UI model.
func New(c common.Common) *UI {
	h := header.New(c, "wdeploy")

	ui := &UI{
		common:     c,
		pages:      make([]common.Component, 4), // pages
		activePage: selectionPage,
		state:      startState,
		header:     h,
		showFooter: true,
	}
	ui.footer = footer.New(c, ui)
	return ui
}

func (ui *UI) getMargins() (wm, hm int) {
	style := ui.common.Styles.App.Copy()
	switch ui.activePage {
	case selectionPage:
		hm += ui.common.Styles.ServerName.GetHeight() +
			ui.common.Styles.ServerName.GetVerticalFrameSize()
	case varsPage:
	case hostsPage:
	case deployPage:
	}
	wm += style.GetHorizontalFrameSize()
	hm += style.GetVerticalFrameSize()
	if ui.showFooter {
		// NOTE: we don't use the footer's style to determine the margins
		// because footer.Height() is the height of the footer after applying
		// the styles.
		hm += ui.footer.Height()
	}
	return
}

// ShortHelp implements help.KeyMap.
func (ui *UI) ShortHelp() []key.Binding {
	b := make([]key.Binding, 0)
	switch ui.state {
	case errorState:
		b = append(b, ui.common.KeyMap.Back)
	case loadedState:
		b = append(b, ui.pages[ui.activePage].ShortHelp()...)
	}
	/*
		if !ui.IsFiltering() {
			b = append(b, ui.common.KeyMap.Quit)
		}
	*/
	b = append(b, ui.common.KeyMap.Help)
	return b
}

// FullHelp implements help.KeyMap.
func (ui *UI) FullHelp() [][]key.Binding {
	b := make([][]key.Binding, 0)
	switch ui.state {
	case errorState:
		b = append(b, []key.Binding{ui.common.KeyMap.Back})
	case loadedState:
		b = append(b, ui.pages[ui.activePage].FullHelp()...)
	}
	h := []key.Binding{
		ui.common.KeyMap.Help,
	}
	/*
		if !ui.IsFiltering() {
			h = append(h, ui.common.KeyMap.Quit)
		}
	*/
	b = append(b, h)
	return b
}

// SetSize implements common.Component.
func (ui *UI) SetSize(width, height int) {
	ui.common.SetSize(width, height)
	wm, hm := ui.getMargins()
	ui.header.SetSize(width-wm, height-hm)
	ui.footer.SetSize(width-wm, height-hm)
	for _, p := range ui.pages {
		if p != nil {
			p.SetSize(width-wm, height-hm)
		}
	}
}

// Init implements tea.Model.
func (ui *UI) Init() tea.Cmd {
	ui.pages[selectionPage] = selection.New(
		ui.common,
	)
	/*
		ui.pages[varsPage] = vars.New(
			ui.common,
		)
		ui.pages[hostsPage] = hosts.New(
			ui.common,
		)
		ui.pages[deployPage] = deploy.New(
			ui.common,
		)
	*/
	ui.SetSize(ui.common.Width, ui.common.Height)
	cmds := make([]tea.Cmd, 0)
	cmds = append(cmds,
		ui.pages[selectionPage].Init(),
		/*
			ui.pages[varsPage].Init(),
			ui.pages[hostsPage].Init(),
			ui.pages[deployPage].Init(),
		*/
	)

	ui.state = loadedState
	ui.SetSize(ui.common.Width, ui.common.Height)
	return tea.Batch(cmds...)
}

// Update implements tea.Model.
func (ui *UI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		ui.SetSize(msg.Width, msg.Height)
		/*
			for i, p := range ui.pages {
				m, cmd := p.Update(msg)
				ui.pages[i] = m.(common.Component)
				if cmd != nil {
					cmds = append(cmds, cmd)
				}
			}

		*/
	case tea.KeyMsg, tea.MouseMsg:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, ui.common.KeyMap.Back) && ui.error != nil:
				ui.error = nil
				ui.state = loadedState
				// Always show the footer on error.
				ui.showFooter = ui.footer.ShowAll()
			case key.Matches(msg, ui.common.KeyMap.Help):
				cmds = append(cmds, footer.ToggleFooterCmd)
			case key.Matches(msg, ui.common.KeyMap.Quit):
				/*if !ui.IsFiltering() {
					// Stop bubblezone background workers.
					ui.common.Zone.Close()
					return ui, tea.Quit
				}
				*/
				ui.common.Zone.Close()
				return ui, tea.Quit
			case ui.activePage == varsPage && key.Matches(msg, ui.common.KeyMap.Back):
				ui.activePage = selectionPage
				// Always show the footer on selection page.
				ui.showFooter = true
			}
		case tea.MouseMsg:
			switch msg.Type {
			case tea.MouseLeft:
				switch {
				case ui.common.Zone.Get("footer").InBounds(msg):
					cmds = append(cmds, footer.ToggleFooterCmd)
				}
			}
		}
	case footer.ToggleFooterMsg:
		ui.footer.SetShowAll(!ui.footer.ShowAll())
		// Show the footer when on repo page and shot all help.
		if ui.error == nil && ui.activePage == varsPage {
			ui.showFooter = !ui.showFooter
		}
		/*
			case app.Action.RepoMsg:
				ui.activePage = varsPage
				// Show the footer on repo page if show all is set.
				ui.showFooter = ui.footer.ShowAll()
		*/
	case common.ErrorMsg:
		ui.error = msg
		ui.state = errorState
		ui.showFooter = true
		return ui, nil
		/*
			case selector.SelectMsg:
				switch msg.IdentifiableItem.(type) {
				case selection.Item:
					if ui.activePage == selectionPage {
						cmds = append(cmds, ui.setActionCmd(msg.ID()))
					}
				}
		*/
	}
	h, cmd := ui.header.Update(msg)
	ui.header = h.(*header.Header)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	f, cmd := ui.footer.Update(msg)
	ui.footer = f.(*footer.Footer)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	if ui.state == loadedState {
		_, cmd := ui.pages[ui.activePage].Update(msg)
		//ui.pages[ui.activePage] = m.(common.Component)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	// This fixes determining the height margin of the footer.
	ui.SetSize(ui.common.Width, ui.common.Height)
	return ui, tea.Batch(cmds...)
}

// View implements tea.Model.
func (ui *UI) View() string {
	var view string
	wm, hm := ui.getMargins()
	switch ui.state {
	case startState:
		view = "Loading..."
	case errorState:
		err := ui.common.Styles.ErrorTitle.Render("Bummer")
		err += ui.common.Styles.ErrorBody.Render(ui.error.Error())
		view = ui.common.Styles.Error.Copy().
			Width(ui.common.Width -
				wm -
				ui.common.Styles.ErrorBody.GetHorizontalFrameSize()).
			Height(ui.common.Height -
				hm -
				ui.common.Styles.Error.GetVerticalFrameSize()).
			Render(err)
	case loadedState:
		view = ui.pages[ui.activePage].View()
	default:
		view = "Unknown state :/ this is a bug!"
	}
	if ui.activePage == selectionPage {
		view = lipgloss.JoinVertical(lipgloss.Left, ui.header.View(), view)
	}
	if ui.showFooter {
		view = lipgloss.JoinVertical(lipgloss.Left, view, ui.footer.View())
	}
	return ui.common.Zone.Scan(
		ui.common.Styles.App.Render(view),
	)
}
