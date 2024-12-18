package microcks

import (
	"fmt"
	"strings"
)

type APIExamples struct {
	APIVersion string                        `json:"apiVersion" yaml:"apiVersion"`
	Kind       string                        `json:"kind" yaml:"kind"`
	Metadata   Metadata                      `json:"metadata" yaml:"metadata"`
	Operations map[OperationName]ExampleItem `json:"operations" yaml:"operations"`
}

type OperationName struct {
	Verb string
	Path string
}

// String method implements the Stringer interface
func (o OperationName) String() string {
	return fmt.Sprintf("%s %s", strings.ToUpper(o.Verb), o.Path)
}

type Metadata struct {
	Name    string `json:"name" yaml:"name"`
	Version string `json:"version" yaml:"version"`
}

type ExampleItem struct {
	Request  Request  `json:"request" yaml:"request"`
	Response Response `json:"response" yaml:"response"`
}

type Request struct {
	Parameters map[string]interface{} `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	Headers    map[string]interface{} `json:"headers,omitempty" yaml:"headers,omitempty"`
	Body       interface{}            `json:"body,omitempty" yaml:"body,omitempty"`
}

type Response struct {
	Headers   map[string]interface{} `json:"headers,omitempty" yaml:"headers,omitempty"`
	MediaType string                 `json:"mediaType,omitempty" yaml:"mediaType,omitempty"`
	Code      string                 `json:"code,omitempty" yaml:"code,omitempty"`
	Body      interface{}            `json:"body,omitempty" yaml:"body,omitempty"`
}
