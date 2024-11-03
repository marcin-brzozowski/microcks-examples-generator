package cmd

import "github.com/spf13/cobra"

// Cmd to api observability
var ApiExamplesCmd = &cobra.Command{
	Use:   "api-examples",
	Short: "Manage Microcks ApiExamples files.",
	Long:  "Enables users to generate Microcks ApiExamples files from OpenAPI specs.",
}

func init() {
	ApiExamplesCmd.AddCommand(GenerateCmd)
}
