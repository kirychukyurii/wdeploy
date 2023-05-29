package custom

var Tmpl = `---
all:
  hosts:
    node1:
      ansible_host: 1.1.1.1
      # ansible_user: admin
      # ansible_port: 2222
      # ansible_ssh_pass: "pAssw0rd"
      # ansible_ssh_private_key_file: /home/webitel/.ssh/rsa.key
      webitel_services:
        - opensips
        - rtpengine
        - nginx
        - webitel_core
        - webitel_engine
        - webitel_call_center
        - webitel_messages

    node2:
      ansible_host: 2.2.2.2
      # ansible_user: admin
      # ansible_port: 2222
      # ansible_ssh_pass: "pAssw0rd"
      # ansible_ssh_private_key_file: /home/webitel/.ssh/rsa.key
      webitel_services:
        - postgresql
        - postgresql_main
        - grafana
        - rabbitmq
        - consul
        - webitel_storage

    node3:
      ansible_host: 3.3.3.3
      # ansible_user: admin
      # ansible_port: 2222
      # ansible_ssh_pass: "pAssw0rd"
      # ansible_ssh_private_key_file: /home/webitel/.ssh/rsa.key
      webitel_services:
        - freeswitch
        - webitel_flow_manager
`
