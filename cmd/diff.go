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
	"os"

	"github.com/spf13/cobra"

	cstate "github.com/starkandwayne/carousel/state"
)

// statusCmd represents the status command
var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show diff in what should be deployed",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		initialize()

		filters.latest = true
		if filters.deployment == "" {
			logger.Fatal("deployment flag must be set")
		}

		cmd.Printf("Refreshing state")
		refresh()
		cmd.Printf(" done\n\n")

		credentials := state.Credentials(filters.Filters()...)
		credentials.SortByNameAndCreatedAt()

		var latestDeployed *cstate.Credential
		exitOne := false
		for _, cred := range credentials {
			latestDeployed = cred.LatestDeployedTo(filters.deployment)
			if latestDeployed != cred {
				if !exitOne {
					exitOne = true
					cmd.Printf("Found credential(s) pending a bosh deploy:\n")
				}
				cmd.Printf("%s\n  + version: %s | %s\n", cred.Path.Name, cred.ID, cred.Summary())
				if latestDeployed != nil {
					cmd.Printf("  - version: %s | %s\n", latestDeployed.ID, latestDeployed.Summary())
				}
			}
		}
		if exitOne {
			os.Exit(1)
		} else {
			cmd.Println("Nothing to deploy")
		}
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)

	addDeploymentFlag(diffCmd.Flags())
}
