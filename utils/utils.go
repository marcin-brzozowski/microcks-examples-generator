package utils

import (
	"fmt"
	"strings"
)

// BodyType is a custom type that implements yaml.Marshaler.
type BodyType string

// MarshalYAML implements the yaml.Marshaler interface.
func (b BodyType) MarshalYAML() (interface{}, error) {
	// Convert the BodyType to a multiline string formatted with |-
	return fmt.Sprintf("|\n%s", strings.ReplaceAll(string(b), "\n", "\n  ")), nil
}
