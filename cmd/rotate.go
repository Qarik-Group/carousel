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
var rotateCmd = &cobra.Command{
	Use:   "rotate",
	Short: "Rotate all credentials needing rotation",
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

		credentialsToAction := []*cstate.Credential{}

		cmd.Printf("Perform actions:\n")

		for _, cred := range credentials {
			switch action := cred.NextAction(regenerationCriteria); {
			case action == cstate.BoshDeploy:
				continue
			case action == cstate.NoOverwrite:
				continue
			case action == cstate.None:
				continue
			default:
				cmd.Printf("- %s %s\n",
					action.String(), cred.PathVersion())
				credentialsToAction = append(credentialsToAction, cred)
			}
		}

		askForConfirmation()

		cmd.Printf("\nPerforming actions:\n")

		for _, cred := range credentialsToAction {
			action := cred.NextAction(regenerationCriteria)
			cmd.Printf("- %s %s",
				action.String(), cred.PathVersion())
			switch action {
			case cstate.Regenerate:
				err := credhub.ReGenerate(cred.Credential)
				if err != nil {
					cmd.Printf(" got error: %s\n", err)
					os.Exit(1)
				}
			case cstate.MarkTransitional:
				err := credhub.UpdateTransitional(cred.Credential, false)
				if err != nil {
					cmd.Printf(" got error: %s\n", err)
					os.Exit(1)
				}
			case cstate.UnMarkTransitional:
				err := credhub.UpdateTransitional(cred.Credential, true)
				if err != nil {
					cmd.Printf(" got error: %s\n", err)
					os.Exit(1)
				}
			case cstate.CleanUp:
				err := credhub.Delete(cred.Credential)
				if err != nil {
					cmd.Printf(" got error: %s\n", err)
					os.Exit(1)
				}
			}
			cmd.Print(" done\n")
		}

		cmd.Println("Finished")
	},
}

func init() {
	rootCmd.AddCommand(rotateCmd)

	addExpiresWithinCriteriaFlag(rotateCmd.Flags())
	addOlderThanCireteriaFlag(rotateCmd.Flags())
	addIgnoreUpdateModeCireteriaFlag(rotateCmd.Flags())
	addNameFlag(rotateCmd.Flags())
}
