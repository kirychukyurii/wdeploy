package view

var Tmpl = `# S U M M A R Y 
## Webitel v{{ .WebitelVersion }}

Hosts: {{ range .Inventory.Inventory.Hosts }}
1. {{ .AnsibleHost }}: {{ range .WebitelServices }} 
   - {{ . }} {{ end }} {{ end }}
`
