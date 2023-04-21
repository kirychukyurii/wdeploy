package vars

var Tmpl = `---
inventory: production

ansible_any_errors_fatal: true
ansible_ignore_unreachable: true
ansible_python_interpreter: /usr/bin/python3

# This setting is always appended to the default ssh command line.
ansible_ssh_extra_args: '-o StrictHostKeyChecking=no'

# The username to use when connecting to the host
# ansible_user: admin

# The connection port number, if not the default (22 for ssh)
# ansible_port: 2222

# Private key file used by ssh. Useful if using multiple keys and you don’t want to use SSH agent.
# ansible_ssh_private_key_file: /home/webitel/.ssh/rsa.key

# The password to use to authenticate to the host
# ansible_ssh_pass: "pAssw0rd"

webitel_version: "23.02"
webitel_repository_user: "{{ .WebitelRepositoryUser }}"
webitel_repository_password: "{{ .WebitelRepositoryPassword }}"

rtpengine_mode: "global"

freeswitch_signalwire_key: ""

opensips_version: "3.2"
opensips_fail2ban: true

nginx_letsencrypt: false
nginx_site_name: "webitel.example.com"
nginx_mail_address: "cloud@example.com"

grafana_enable: true
grafana_basic_dashboards: false
grafana_basic_dashboards_language: "en"

# Generate additional locales
locales_gen:
  - "en_US.UTF-8"
`
