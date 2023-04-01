package models

type Inventory struct {
	Inventory Hosts `mapstructure:"all"`
}

type Hosts struct {
	Hosts map[string]Host `mapstructure:"hosts"`
}

type Host struct {
	AnsibleGlobalVariables
	AnsibleHost     string   `mapstructure:"ansible_host"`
	WebitelServices []string `mapstructure:"webitel_services"`
}
