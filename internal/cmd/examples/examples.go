package examples

import (
	"github.com/spf13/cobra"
)

// Cmd to api observability
var Cmd = &cobra.Command{
	Use:   "api-examples",
	Short: "Manage Microcks ApiExamples files.",
	Long:  "Enables users to generate Microcks ApiExamples files from OpenAPI specs.",
}

func init() {
	Cmd.AddCommand(GenerateCmd)
}
