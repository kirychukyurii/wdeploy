package config

type Variables struct {
	AnsibleAnyErrorsFatal    bool   `mapstructure:"ansible_any_errors_fatal" yaml:"ansible_any_errors_fatal"`
	AnsibleIgnoreUnreachable bool   `mapstructure:"ansible_ignore_unreachable" yaml:"ansible_ignore_unreachable"`
	AnsiblePythonInterpreter string `mapstructure:"ansible_python_interpreter" yaml:"ansible_python_interpreter"`
	AnsibleSSHExtraArgs      string `mapstructure:"ansible_ssh_extra_args" yaml:"ansible_ssh_extra_args"`             // This setting is always appended to the default ssh command line
	AnsibleUser              string `mapstructure:"ansible_user" yaml:"ansible_user"`                                 // The username to use when connecting to the host
	AnsiblePort              int    `mapstructure:"ansible_port" yaml:"ansible_port"`                                 // The connection port number, if not the default (22 for ssh)
	AnsibleSSHPrivateKeyFile string `mapstructure:"ansible_ssh_private_key_file" yaml:"ansible_ssh_private_key_file"` // Private key file used by ssh. Useful if using multiple keys and you don’t want to use SSH agent
	AnsibleSSHPass           string `mapstructure:"ansible_ssh_pass" yaml:"ansible_ssh_pass"`                         // The password to use to authenticate to the host

	WebitelVersion            string `mapstructure:"webitel_version" yaml:"webitel_version"`
	WebitelRepositoryUser     string `mapstructure:"webitel_repository_user" yaml:"webitel_repository_user"`
	WebitelRepositoryPassword string `mapstructure:"webitel_repository_password" yaml:"webitel_repository_password"`

	RTPEngineMode           string `mapstructure:"rtp_engine_mode" yaml:"rtp_engine_mode"`
	FreeswitchSignalwireKey string `mapstructure:"freeswitch_signalwire_key" yaml:"freeswitch_signalwire_key"`
	OpensipsVersion         string `mapstructure:"opensips_version" yaml:"opensips_version"`
	OpensipsFail2ban        bool   `mapstructure:"opensips_fail_2_ban" yaml:"opensips_fail_2_ban"`

	NginxLetsencrypt               bool   `mapstructure:"nginx_letsencrypt" yaml:"nginx_letsencrypt"`
	NginxSiteName                  string `mapstructure:"nginx_site_name" yaml:"nginx_site_name"`
	NginxMailAddress               string `mapstructure:"nginx_mail_address" yaml:"nginx_mail_address"`
	GrafanaEnable                  bool   `mapstructure:"grafana_enable" yaml:"grafana_enable"`
	GrafanaBasicDashboards         bool   `mapstructure:"grafana_basic_dashboards" yaml:"grafana_basic_dashboards"`
	GrafanaBasicDashboardsLanguage string `mapstructure:"grafana_basic_dashboards_language" yaml:"grafana_basic_dashboards_language"`

	LocalesGen []string `mapstructure:"locales_gen,omitempty" yaml:"locales_gen,omitempty"`
}

type Inventory struct {
	Inventory Hosts `mapstructure:"all" yaml:"all"`
}

type Hosts struct {
	Hosts map[string]Host `mapstructure:"hosts" yaml:"hosts"`
}

type Host struct {
	AnsibleHost     string   `mapstructure:"ansible_host" yaml:"ansible_host"`
	WebitelServices []string `mapstructure:"webitel_services" yaml:"webitel_services"`
}
