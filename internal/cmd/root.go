// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"cmd/examples"

	"github.com/spf13/cobra"
)

// RootCmd to manage openapi-examples
var RootCmd = &cobra.Command{
	Use:   "openapi-examples",
	Short: "openapi-examples can be used to generate examples for API requests, responses or models.",
	Long: `openapi-examples can be used to generate examples for API requests, responses or models 
from schema defined in OpenAPI specification document.`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		// clilog.Error.Println(err)
	}
}

var (
	printOutput, noOutput bool
)

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().BoolVarP(&printOutput, "print-output", "",
		true, "Control printing of info log statements")

	RootCmd.PersistentFlags().BoolVarP(&noOutput, "no-output", "",
		false, "Disable printing all statements to stdout")

	RootCmd.AddCommand(examples.Cmd)
}

func initConfig() {

}

// GetRootCmd returns the root of the cobra command-tree.
func GetRootCmd() *cobra.Command {
	return RootCmd
}
