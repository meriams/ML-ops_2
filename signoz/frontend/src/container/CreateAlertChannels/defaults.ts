import { OpsgenieChannel, PagerChannel } from './config';

export const PagerInitialConfig: Partial<PagerChannel> = {
	description: `[{{ .Status | toUpper }}{{ if eq .Status "firing" }}:{{ .Alerts.Firing | len }}{{ end }}] {{ .CommonLabels.alertname }} for {{ .CommonLabels.job }}
	{{- if gt (len .CommonLabels) (len .GroupLabels) -}}
	  {{" "}}(
	  {{- with .CommonLabels.Remove .GroupLabels.Names }}
		{{- range $index, $label := .SortedPairs -}}
		  {{ if $index }}, {{ end }}
		  {{- $label.Name }}="{{ $label.Value -}}"
		{{- end }}
	  {{- end -}}
	  )
	{{- end }}`,
	severity: '{{ (index .Alerts 0).Labels.severity }}',
	client: 'SigNoz Alert Manager',
	client_url: 'https://enter-signoz-host-n-port-here/alerts',
	details: JSON.stringify({
		firing: `{{ template "pagerduty.default.instances" .Alerts.Firing }}`,
		resolved: `{{ template "pagerduty.default.instances" .Alerts.Resolved }}`,
		num_firing: '{{ .Alerts.Firing | len }}',
		num_resolved: '{{ .Alerts.Resolved | len }}',
	}),
};

export const OpsgenieInitialConfig: Partial<OpsgenieChannel> = {
	message: '{{ .CommonLabels.alertname }}',
	description: `{{ if gt (len .Alerts.Firing) 0 -}}
	Alerts Firing:
	{{ range .Alerts.Firing }}
	 - Message: {{ .Annotations.description }}
	Labels:
	{{ range .Labels.SortedPairs }}   - {{ .Name }} = {{ .Value }}
	{{ end }}   Annotations:
	{{ range .Annotations.SortedPairs }}   - {{ .Name }} = {{ .Value }}
	{{ end }}   Source: {{ .GeneratorURL }}
	{{ end }}
	{{- end }}
	{{ if gt (len .Alerts.Resolved) 0 -}}
	Alerts Resolved:
	{{ range .Alerts.Resolved }}
	 - Message: {{ .Annotations.description }}
	Labels:
	{{ range .Labels.SortedPairs }}   - {{ .Name }} = {{ .Value }}
	{{ end }}   Annotations:
	{{ range .Annotations.SortedPairs }}   - {{ .Name }} = {{ .Value }}
	{{ end }}   Source: {{ .GeneratorURL }}
	{{ end }}
	{{- end }}`,
	priority:
		'{{ if eq (index .Alerts 0).Labels.severity "critical" }}P1{{ else if eq (index .Alerts 0).Labels.severity "warning" }}P2{{ else if eq (index .Alerts 0).Labels.severity "info" }}P3{{ else }}P4{{ end }}',
};
