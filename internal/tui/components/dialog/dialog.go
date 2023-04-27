package dialog

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kirychukyurii/wdeploy/internal/tui/common"
	"strings"
)

// SelectDialogButtonMsg is a message that contains the index of the tab to select.
type SelectDialogButtonMsg int

// ActiveDialogButtonMsg is a message that contains the index of the current active tab.
type ActiveDialogButtonMsg int

type Dialog struct {
	common       common.Common
	activeButton int
	question     string
	buttons      []string
}

func New(c common.Common, question string, buttons []string) *Dialog {
	d := &Dialog{
		common:   c,
		question: question,
		buttons:  buttons,
	}

	return d
}

// SetSize implements common.Component.
func (d *Dialog) SetSize(width, height int) {
	d.common.SetSize(width, height)
}

// Init implements tea.Model.
func (d *Dialog) Init() tea.Cmd {
	d.activeButton = 0
	return nil
}

// Update implements tea.Model.
func (d *Dialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Recalculate content width and line wrap.
		cmds = append(cmds, d.Init())
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, d.common.KeyMap.Left):
			d.activeButton--
			if d.activeButton < 0 {
				d.activeButton = 0
			}
			cmds = append(cmds, d.activeButtonCmd)
		case key.Matches(msg, d.common.KeyMap.Right):
			d.activeButton++
			if d.activeButton >= len(d.buttons) {
				d.activeButton = len(d.buttons) - 1
			}
			cmds = append(cmds, d.activeButtonCmd)
		case key.Matches(msg, d.common.KeyMap.Select):
			cmds = append(cmds, d.selectButtonCmd)
		}
	case SelectDialogButtonMsg:
		button := int(msg)
		if button >= 0 && button < len(d.buttons) {
			d.activeButton = int(msg)
		}
	}

	return d, tea.Batch(cmds...)
}

// View implements tea.Model.
func (d *Dialog) View() string {
	s := strings.Builder{}
	dialogBoxStyle := d.common.Styles.Dialog.Box.Copy()

	for i, button := range d.buttons {
		buttonStyle := d.common.Styles.Dialog.Normal.Button.Copy()
		if i == d.activeButton {
			buttonStyle = d.common.Styles.Dialog.Active.Button.Copy()
		}
		//s += buttonStyle.Render(button)
		s.WriteString(buttonStyle.Render(button))
	}

	// question := lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render(d.question)
	question := d.common.Styles.Dialog.Question.Copy().Render(d.question)
	buttons := lipgloss.JoinHorizontal(lipgloss.Top, s.String())
	ui := lipgloss.JoinVertical(lipgloss.Center, question, buttons)

	dialog := lipgloss.Place(d.common.Width, d.common.Height,
		lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Render(ui),
	)

	return dialog
}

func (d *Dialog) selectButtonCmd() tea.Msg {
	return SelectDialogButtonMsg(d.activeButton)
}

func (d *Dialog) activeButtonCmd() tea.Msg {
	return ActiveDialogButtonMsg(d.activeButton)
}
