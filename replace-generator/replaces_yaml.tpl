replaces:
{{- range . }}
{{- if .Comment }}
  # {{ .Comment }}
{{- end }}
  - {{ .Dependency }} => {{ .Replacement }}
{{- end }}

