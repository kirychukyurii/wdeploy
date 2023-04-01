package app

type ActionItem struct {
	SetID          string
	SetName        string
	SetDescription string
}

func (a ActionItem) Description() string {
	return a.SetDescription
}

func (a ActionItem) ID() string {
	return a.SetID
}

func (a ActionItem) FilterValue() string {
	return a.SetName
}

func (a ActionItem) Title() string {
	return a.SetName
}

func (a ActionItem) NewActionItem() []ActionItem {
	return []ActionItem{
		{
			SetID:          "vars",
			SetName:        "Configure vars",
			SetDescription: "test var",
		},
		{
			SetID:          "hosts",
			SetName:        "Configure hosts",
			SetDescription: "test host",
		},
		{
			SetID:          "deploy",
			SetName:        "Deploy Webitel",
			SetDescription: "test deploy",
		},
	}
}
