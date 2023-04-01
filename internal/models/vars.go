package models

type Variables struct {
	AnsibleGlobalVariables
	TelephonyVariables
	WebVariables
	LocalesGen []string `mapstructure:"locales_gen,omitempty"`
}

type AnsibleGlobalVariables struct {
	AnsibleAnyErrorsFatal    bool   `mapstructure:"ansible_any_errors_fatal"`
	AnsibleIgnoreUnreachable bool   `mapstructure:"ansible_ignore_unreachable"`
	AnsiblePythonInterpreter string `mapstructure:"ansible_python_interpreter"`
	AnsibleSSHExtraArgs      string `mapstructure:"ansible_ssh_extra_args"`       // This setting is always appended to the default ssh command line
	AnsibleUser              string `mapstructure:"ansible_user"`                 // The username to use when connecting to the host
	AnsiblePort              int    `mapstructure:"ansible_port"`                 // The connection port number, if not the default (22 for ssh)
	AnsibleSSHPrivateKeyFile string `mapstructure:"ansible_ssh_private_key_file"` // Private key file used by ssh. Useful if using multiple keys and you don’t want to use SSH agent
	AnsibleSSHPass           string `mapstructure:"ansible_ssh_pass"`             // The password to use to authenticate to the host
}

type WebitelVariables struct {
	WebitelVersion            string `mapstructure:"webitel_version"`
	WebitelRepositoryUser     string `mapstructure:"webitel_repository_user"`
	WebitelRepositoryPassword string `mapstructure:"webitel_repository_password"`
}

type TelephonyVariables struct {
	RTPEngineMode           string `mapstructure:"rtp_engine_mode"`
	FreeswitchSignalwireKey string `mapstructure:"freeswitch_signalwire_key"`
	OpensipsVersion         string `mapstructure:"opensips_version"`
	OpensipsFail2ban        bool   `mapstructure:"opensips_fail_2_ban"`
}

type WebVariables struct {
	NginxLetsencrypt               bool   `mapstructure:"nginx_letsencrypt"`
	NginxSiteName                  string `mapstructure:"nginx_site_name"`
	NginxMailAddress               string `mapstructure:"nginx_mail_address"`
	GrafanaEnable                  bool   `mapstructure:"grafana_enable"`
	GrafanaBasicDashboards         bool   `mapstructure:"grafana_basic_dashboards"`
	GrafanaBasicDashboardsLanguage string `mapstructure:"grafana_basic_dashboards_language"`
}
