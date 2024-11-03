package microcks

import (
	"context"
	"fmt"

	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
	"github.com/pb33f/libopenapi/renderer"
)

// Function to load OpenAPI file, generate examples, and construct APIExamples struct
func GenerateAPIExamples(ctx context.Context, spec []byte) (*APIExamples, error) {
	// Parse OpenAPI 3 document
	doc, docErr := libopenapi.NewDocument(spec)
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
	apiExamples := &APIExamples{
		APIVersion: "mocks.microcks.io/v1alpha1",
		Kind:       "APIExamples",
		Metadata: Metadata{
			Name:    docModel.Model.Info.Title,
			Version: docModel.Model.Info.Version,
		},
		Operations: make(map[OperationName]ExampleItem),
	}

	// Iterate through paths and operations to generate examples
	for p := range orderedmap.Iterate(ctx, docModel.Model.Paths.PathItems) {
		path, pathItem := p.Key(), p.Value()

		exampleItem := ExampleItem{}

		for op := range orderedmap.Iterate(ctx, pathItem.GetOperations()) {
			verb, operation := op.Key(), op.Value()

			// Generate request example
			reqMock, reqErr := generateMockRequest(ctx, operation, mockGen)
			if reqErr == nil {
				exampleItem.Request = reqMock
			}

			// Generate response examples
			for r := range orderedmap.Iterate(ctx, operation.Responses.Codes) {
				status, response := r.Key(), r.Value()
				resMock, resErr := generateMockResponse(ctx, response, mockGen, status)
				if resErr == nil {
					exampleItem.Response = resMock
				}
			}

			apiExamples.Operations[OperationName{Verb: verb, Path: path}] = exampleItem
		}
	}
	return apiExamples, nil
}

// GenerateMockRequest generates a mock request example using the operation parameters and request body schema.
func generateMockRequest(ctx context.Context, operation *v3.Operation, mockGen *renderer.MockGenerator) (Request, error) {
	mockRequest := Request{
		Parameters: make(map[string]interface{}),
		Headers:    make(map[string]interface{}),
	}

	// Generate mock data for parameters
	for _, param := range operation.Parameters {
		if param.Schema != nil {
			schema, err := param.Schema.BuildSchema()
			if err != nil {
				return Request{}, err
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
					return Request{}, err
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
func generateMockResponse(ctx context.Context, response *v3.Response, mockGen *renderer.MockGenerator, statusCode string) (Response, error) {
	mockResponse := Response{
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
				return Response{}, err
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
				return Response{}, err
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
