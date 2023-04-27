package spinner

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kirychukyurii/wdeploy/internal/tui/common"
)

type Spinner struct {
	common   common.Common
	spinner  spinner.Model
	quitting bool
	err      error
}

func New(c common.Common, quitting bool) *Spinner {
	s := spinner.New()
	s.Spinner = spinner.Dot

	d := &Spinner{
		common:   c,
		spinner:  s,
		quitting: quitting,
	}

	return d
}

// SetSize implements common.Component.
func (d *Spinner) SetSize(width, height int) {
	d.common.SetSize(width, height)
}

// Init implements tea.Model.
func (d *Spinner) Init() tea.Cmd {
	return d.spinner.Tick
}

// Update implements tea.Model.
func (d *Spinner) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	return d, tea.Batch(cmds...)
}

// View implements tea.Model.
func (d *Spinner) View() string {
	str := fmt.Sprintf("\n\n   %s Loading forever...press q to quit\n\n", d.spinner.View())
	if d.quitting {
		return str + "\n"
	}
	return str
}
