package view

var Tmpl = `# SUMMARY Webitel v{{ .WebitelVersion }}

| Name     | Server      | Services  |
|----------|-------------|-----------|
{{ range $k, $v := .Inventory.Inventory.Hosts -}}
| {{ $k }} | {{ $v.AnsibleHost }} | {{ range $v.WebitelServices }} {{ . }} |
|          |             | {{ end }} |
{{- end }}
`
