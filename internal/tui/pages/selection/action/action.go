package action

type Action interface {
	ID() string
	Title() string
	Description() string
}

type ActionItem struct {
	Command string
	Name    string
	Action  string
}

type ActionItems []ActionItem

func (a ActionItem) Description() string {
	return a.Action
}

func (a ActionItem) ID() string {
	return a.Command
}

func (a ActionItem) FilterValue() string {
	return a.Name
}

func (a ActionItem) Title() string {
	return a.Name
}
