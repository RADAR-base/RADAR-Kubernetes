{{- define "radar-seaweedfs.write-config" -}}
{{- $users := list }}
{{- range $username, $fields := .Values.s3Config.identities }}
{{- $credentials := list $fields.credentials }}
{{- $actions := $fields.actions }}
{{- $user := dict "name" $username "credentials" $credentials  "actions" $actions }}
{{- $users = append $users $user }}
{{- end }}
{{ dict "identities" $users | toJson }}
{{- end -}}
