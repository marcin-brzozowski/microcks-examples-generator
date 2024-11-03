package cmd

import (
	"github.com/spf13/cobra"
)

var GenerateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"g", "gen"},
	Short:   "Generates examples from OpenAPI spec.",
	Long:    "Generates examples in different output formats from OpenAPI spec.",
}

func init() {
	GenerateCmd.AddCommand(MicrocksApiExamplesCmd)
}
