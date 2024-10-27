package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
	"github.com/pb33f/libopenapi/renderer"
	"libopenapi-poc.com/m/microcks"
)

// Function to load OpenAPI file, generate examples, and construct APIExamples struct
func generateAPIExamples(ctx context.Context, openAPIPath string) (*microcks.APIExamples, error) {
	// Load the OpenAPI spec file
	specData, err := os.ReadFile(openAPIPath)
	if err != nil {
		return nil, fmt.Errorf("error reading OpenAPI file: %w", err)
	}

	// Parse OpenAPI 3 document
	doc, docErr := libopenapi.NewDocument(specData)
	if docErr != nil {
		return nil, fmt.Errorf("error parsing OpenAPI file: %w", docErr)
	}

	// because we know this is a v3 spec, we can build a ready to go model from it.
	docModel, errors := doc.BuildV3Model()

	// if anything went wrong when building the v3 model, a slice of errors will be returned
	if len(errors) > 0 {
		for i := range errors {
			fmt.Printf("error: %e\n", errors[i])
		}
		panic(fmt.Sprintf("cannot create v3 model from document: %d errors reported",
			len(errors)))
	}

	// Initialize mock generator
	mockGen := renderer.NewMockGenerator(renderer.JSON)
	mockGen.SetPretty()

	// Construct the APIExamples struct
	apiExamples := &microcks.APIExamples{
		APIVersion: "mocks.microcks.io/v1alpha1",
		Kind:       "APIExamples",
		Metadata: microcks.Metadata{
			Name:    docModel.Model.Info.Title,
			Version: docModel.Model.Info.Version,
		},
		Operations: make(map[microcks.OperationName]microcks.ExampleItem),
	}

	// Iterate through paths and operations to generate examples
	for p := range orderedmap.Iterate(ctx, docModel.Model.Paths.PathItems) {
		path, pathItem := p.Key(), p.Value()

		exampleItem := microcks.ExampleItem{}

		for op := range orderedmap.Iterate(ctx, pathItem.GetOperations()) {
			verb, operation := op.Key(), op.Value()

			// Generate request example
			reqMock, reqErr := GenerateMockRequest(ctx, operation, mockGen)
			if reqErr == nil {
				exampleItem.Request = reqMock
			}

			// Generate response examples
			for r := range orderedmap.Iterate(ctx, operation.Responses.Codes) {
				status, response := r.Key(), r.Value()
				resMock, resErr := GenerateMockResponse(ctx, response, mockGen, status)
				if resErr == nil {
					exampleItem.Response = resMock
				}
			}

			apiExamples.Operations[microcks.OperationName{Verb: verb, Path: path}] = exampleItem
		}
	}
	return apiExamples, nil
}

// GenerateMockRequest generates a mock request example using the operation parameters and request body schema.
func GenerateMockRequest(ctx context.Context, operation *v3.Operation, mockGen *renderer.MockGenerator) (microcks.Request, error) {
	mockRequest := microcks.Request{
		Parameters: make(map[string]interface{}),
		Headers:    make(map[string]interface{}),
	}

	// Generate mock data for parameters
	for _, param := range operation.Parameters {
		if param.Schema != nil {
			schema, err := param.Schema.BuildSchema()
			if err != nil {
				return microcks.Request{}, err
			}
			paramExample, err := mockGen.GenerateMock(schema, "")
			if err != nil {
				return mockRequest, fmt.Errorf("error generating mock parameter: %w", err)
			}
			mockRequest.Parameters[param.Name] = string(paramExample)
		}
	}

	// Generate mock data for request body if defined
	if operation.RequestBody != nil && operation.RequestBody.Content != nil {
		for mediaType := range orderedmap.Iterate(ctx, operation.RequestBody.Content) {
			mediaTypeValue := mediaType.Value()
			if mediaTypeValue.Schema != nil {
				schema, err := mediaTypeValue.Schema.BuildSchema()
				if err != nil {
					return microcks.Request{}, err
				}
				bodyExample, err := mockGen.GenerateMock(schema, "")
				if err != nil {
					return mockRequest, fmt.Errorf("error generating mock request body: %w", err)
				}
				// Set the body in the mock request
				mockRequest.Body = string(bodyExample)
				break // Assume one media type for simplicity
			}
		}
	}

	return mockRequest, nil
}

// GenerateMockResponse generates a mock response example for a given response schema and status code.
func GenerateMockResponse(ctx context.Context, response *v3.Response, mockGen *renderer.MockGenerator, statusCode string) (microcks.Response, error) {
	mockResponse := microcks.Response{
		Headers: make(map[string]interface{}),
		Code:    statusCode,
	}

	// Generate mock data for response headers
	for header := range orderedmap.Iterate(ctx, response.Headers) {
		headerName := header.Key()
		headerValue := header.Value()

		if headerValue.Schema != nil {
			schema, err := headerValue.Schema.BuildSchema()
			if err != nil {
				return microcks.Response{}, err
			}
			headerExample, err := mockGen.GenerateMock(schema, "")
			if err != nil {
				return mockResponse, fmt.Errorf("error generating mock header: %w", err)
			}
			mockResponse.Headers[headerName] = string(headerExample)
		}
	}

	// Generate mock data for response body if defined
	for mediaType := range orderedmap.Iterate(ctx, response.Content) {
		mediaTypeName, mediaTypeValue := mediaType.Key(), mediaType.Value()

		// TODO: reduce nesting and log schema issues
		if mediaTypeValue.Schema != nil {
			schema, err := mediaTypeValue.Schema.BuildSchema()
			if err != nil {
				return microcks.Response{}, err
			}
			bodyExample, err := mockGen.GenerateMock(schema, "")
			if err != nil {
				return mockResponse, fmt.Errorf("error generating mock response body: %w", err)
			}
			mockResponse.Body = string(bodyExample)
			mockResponse.MediaType = mediaTypeName
			break // Assume one media type for simplicity
		}
	}

	return mockResponse, nil
}

func main() {
	openAPIPath := "assets/specs/petstore.yaml"
	ctx := context.TODO()

	apiExamples, err := generateAPIExamples(ctx, openAPIPath)
	if err != nil {
		log.Fatalf("Error generating API examples: %v", err)
	}

	// Define the template
	// TODO - currently only a single example is generated, and the "Example" property is hardcoded
	// Implement generating multiple examples, one for each status code returned by the endpoint
	tmpl := `
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
	// Create a new template and parse the letter into it
	t := template.Must(template.New("apiExamples").Funcs(template.FuncMap{
		"PrettyJSON": PrettyJSON,
		"SafeJSON":   SafeJSON,
		"indent":     indent,
	}).Parse(tmpl))

	// Create a buffer to hold the output
	var output bytes.Buffer

	// Execute the template with the apiExamples data
	if err := t.Execute(&output, apiExamples); err != nil {
		panic(err)
	}

	// Output the generated YAML
	fmt.Println(output.String())

	// // Convert to YAML and print result
	// yamlData, err := yaml.Marshal(apiExamples)
	// if err != nil {
	// 	log.Fatalf("Error marshalling API examples to YAML: %v", err)
	// }
	// fmt.Println(string(yamlData))
}
func indent(amount int, html template.HTML) template.HTML {
	// Convert template.HTML back to string for processing
	str := string(html)

	pad := strings.Repeat(" ", amount)
	// Indent each line by adding the padding
	indentedStr := pad + strings.ReplaceAll(str, "\n", "\n"+pad)

	// Return the indented string as template.HTML
	return template.HTML(indentedStr)
}

func SafeJSON(input string) (template.HTML, error) {
	return template.HTML(input), nil
}

// This is bullshit, I tried to pretty-print JSON
// It turned out it's as simple as wrapping json in template.HTML to avoid escaping characters
// and then go template will render it as raw string, then we just need the indent function
// to work with the template.HTML input and return it also as the output and we're golden
func PrettyJSON(input interface{}) (template.HTML, error) {
	var formattedJSON bytes.Buffer

	inputStr, ok := input.(string)
	if !ok {
		return "", fmt.Errorf("value is not a string")
	}

	// Declare an interface{} to hold the unmarshalled data
	var result interface{}

	// Unmarshal the JSON string into the interface{}
	err := json.Unmarshal([]byte(inputStr), &result)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return "", fmt.Errorf("value is not a valid JSON string")
	}

	encoder := json.NewEncoder(&formattedJSON)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(result); err != nil {
		return "", err
	}
	return template.HTML(formattedJSON.String()), nil // Mark as HTML safe
}
