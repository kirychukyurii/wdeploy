package view

var Tmpl = `# SUMMARY Webitel v{{ .WebitelVersion }}
Hosts: {{ range .Inventory.Inventory.Hosts }}
1. {{ .AnsibleHost }}:
   - services: {{ range .WebitelServices }}{{ . }} {{ end }} {{ end }}
1
1
1
1
1
1
1
1
1
1
1
1
1
1
1
1
1
1
1
1
1
1
1
1
1
1
1
1
1
1
1
1
1
1
1
1
1
1
1
1
1
`
