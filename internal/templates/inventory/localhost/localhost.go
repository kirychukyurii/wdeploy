package localhost

var Tmpl = `---
all:
  hosts:
    node1:
      ansible_host: localhost
      ansible_connection: local
      # ansible_user: admin
      # ansible_port: 2222
      # ansible_ssh_pass: "pAssw0rd"
      # ansible_ssh_private_key_file: /home/webitel/.ssh/rsa.key
      webitel_services:
        - consul
        - rabbitmq
        - postgresql
        - postgresql_main
        - grafana
        - freeswitch
        - rtpengine
        - opensips
        - nginx
        - webitel_core
        - webitel_engine
        - webitel_call_center
        - webitel_flow_manager
        - webitel_storage
        - webitel_messages
`
