package examples

import (
	"bytes"
	"examplesgenerator/microcks"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"examplesgenerator"
)

// GenerateCmd
var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates Microcks APIExamples files from OpenAPI specs",
	Long:  "Generates Microcks APIExamples files from OpenAPI specs",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if _, err := os.Stat(openApiFilePath); os.IsNotExist(err) {
			return fmt.Errorf("OpenAPI file does not exist: %v", err)
		}

		apiExamples, err := examplesgenerator.GenerateAPIExamples(cmd.Context(), openApiFilePath)
		if err != nil {
			return fmt.Errorf("error generating API examples: %v", err)
		}

		var output bytes.Buffer
		if err := microcks.RenderApiExamples(apiExamples, &output); err != nil {
			return fmt.Errorf("failed to render API examples: %v", err)
		}

		fmt.Fprintln(cmd.OutOrStdout(), output.String())
		return
	},
}

var (
	openApiFilePath string
)

func init() {
	GenerateCmd.Flags().StringVarP(&openApiFilePath, "file", "f",
		"", "OpenAPI File Path")

	_ = GenerateCmd.MarkFlagRequired("file")
}
