package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/marcin-brzozowski/openapi-examples/microcks"
	"github.com/spf13/cobra"
)

var MicrocksApiExamplesCmd = &cobra.Command{
	Use:     "microcks-api-examples",
	Short:   "Generates Microcks APIExamples files from OpenAPI specs",
	Long:    "Generates Microcks APIExamples files from OpenAPI specs",
	Example: "openapi-examples generate microcks-api-examples openapi.yaml",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		openApiContent, err := getOpenApiContent(args)
		if err != nil {
			return fmt.Errorf("error reading OpenAPI content: %v", err)
		}
		apiExamples, err := microcks.GenerateAPIExamples(cmd.Context(), openApiContent)
		if err != nil {
			return fmt.Errorf("error generating API examples: %v", err)
		}

		var output bytes.Buffer
		if err := microcks.RenderApiExamples(apiExamples, &output); err != nil {
			return fmt.Errorf("failed to render API examples: %v", err)
		}

		fmt.Fprintln(cmd.OutOrStdout(), output.String())

		return nil
	},
}

func init() {
}

func getOpenApiContent(args []string) ([]byte, error) {
	// Read OpenAPI content from file
	if len(args) > 0 {
		filePath := args[0]
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return nil, fmt.Errorf("provided file path does not exist: %v", err)
		}
		return os.ReadFile(filePath)
	}

	// Read OpenAPI content from stdin when no file path is provided
	stat, err := os.Stdin.Stat()
	if err != nil {
		return nil, fmt.Errorf("error stating stdin: %v", err)
	}

	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return nil, fmt.Errorf("failed to read from standard input")
	}

	return io.ReadAll(os.Stdin)
}
