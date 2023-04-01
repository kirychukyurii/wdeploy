package app

type Action interface {
	ID() string
	Title() string
	Description() string
}

type Actions interface {
	GetAction(string) (Action, error)
	NewActionItem() []Action
}
