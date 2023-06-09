package common

import (
	"github.com/kirychukyurii/wdeploy/internal/tui/keymap"
	"github.com/kirychukyurii/wdeploy/internal/tui/styles"
	zone "github.com/lrstanley/bubblezone"
)

// Common is a struct all components should embed.
type Common struct {
	// Copy   *osc52.Output
	Styles *styles.Styles
	KeyMap *keymap.KeyMap
	Width  int
	Height int
	Zone   *zone.Manager
}

// SetSize sets the width and height of the common struct.
func (c *Common) SetSize(width, height int) {
	c.Width = width
	c.Height = height
}
