{{- range . }}
{{- if .Comment }}
// {{ .Comment }}
{{- end }}
replace {{ .Dependency }} => {{ .Replacement }}
{{ end -}}

