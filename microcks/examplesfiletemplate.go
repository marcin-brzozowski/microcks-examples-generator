package microcks

import (
	"io"
	"text/template"

	"github.com/marcin-brzozowski/openapi-examples/utils"
)

// Define the template
// TODO - currently only a single example is generated, and the "Example" property is hardcoded
// Implement generating multiple examples, one for each status code returned by the endpoint
const microcksApiExamplesTemplate = `
apiVersion: {{ .APIVersion }}
kind: {{ .Kind }}
metadata:
  name: {{ .Metadata.Name }}
  version: {{ .Metadata.Version }}
operations:
{{- range $key, $value := .Operations }}
  '{{$key}}':
    Example:
      request:
        {{- if $value.Request.Parameters }}
        parameters:
          {{- range $paramKey, $paramValue := $value.Request.Parameters }}
          {{$paramKey}}: {{$paramValue}}
          {{- end }}
        {{- end }}
        {{- if $value.Request.Headers }}
        headers:
          {{- range $headerKey, $headerValue := $value.Request.Headers }}
          {{$headerKey}}: {{$headerValue}}
          {{- end }}
        {{- end }}
        {{- if $value.Request.Body }}
        body: |-
{{ SafeJSON $value.Request.Body | indent 10 }}
        {{- end }}
      response:
        {{- if $value.Response.Headers }}
        headers:
          {{- range $headerKey, $headerValue := $value.Response.Headers }}
          {{$headerKey}}: {{$headerValue}}
          {{- end }}
        {{- end }}
        {{- if $value.Response.MediaType }}
        mediaType: {{$value.Response.MediaType}}
        {{- end }}
        code: {{$value.Response.Code}}
        {{- if $value.Response.Body }}
        body: |-
{{ SafeJSON $value.Response.Body | indent 10 }}
        {{- end }}
{{- end }}
`

// Function to execute the template
func RenderApiExamples(apiExamples *APIExamples, w io.Writer) error {
	return getTemplate().Execute(w, apiExamples)
}

func getTemplate() *template.Template {
	t := template.Must(template.New("apiExamples").Funcs(template.FuncMap{
		"SafeJSON": utils.SafeJSON,
		"indent":   utils.Indent,
	}).Parse(microcksApiExamplesTemplate))

	return t
}
