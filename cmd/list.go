/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Display a list of credentials",
	Long: `List CredHub credentials augmented with information from the BOSH director:
* update_mode: looked up from runtime configs and deployment manifest 'variables:' sections
* deployments: list of deployment names which use this version of the credential`,
	Run: func(cmd *cobra.Command, args []string) {
		initialize()
		refresh()

		out, err := json.Marshal(state.Credentials(filters.Filters()...))
		if err != nil {
			logger.Fatalf("failed to mashal: %s", err)
		}

		fmt.Println(string(out))
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	addFilterFlags(listCmd.Flags())
	addNameFlag(listCmd.Flags())
	addDeploymentFlag(listCmd.Flags())
}
