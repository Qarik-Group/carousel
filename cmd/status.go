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
	"github.com/spf13/cobra"

	cstate "github.com/starkandwayne/carousel/state"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Display status of credentials",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		initialize()

		regenerationCriteria, err := criteria.RegenerationCriteria()
		if err != nil {
			logger.Fatal(err)
		}

		refresh()

		credentials := state.Credentials(filters.Filters()...)
		credentials.SortByNameAndCreatedAt()
		for _, cred := range credentials {
			switch action := cred.NextAction(regenerationCriteria); {
			case action == cstate.None:
				continue
			case action == cstate.BoshDeploy:
				cmd.Printf("- %s(%s) %s\n",
					action.String(), cred.PendingDeploys().String(), cred.PathVersion())
			default:
				cmd.Printf("- %s %s\n",
					action.String(), cred.PathVersion())
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	addExpiresWithinCriteriaFlag(statusCmd.Flags())
	addOlderThanCireteriaFlag(statusCmd.Flags())
	addIgnoreUpdateModeCireteriaFlag(statusCmd.Flags())
	addNameFlag(statusCmd.Flags())
	addSignedByFlag(statusCmd.Flags())
}
