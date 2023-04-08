package app

type Action interface {
	ID() string
	Title() string
	Description() string
}
