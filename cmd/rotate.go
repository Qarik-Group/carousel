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

		var credentialsToDeploy, credentialsToAction cstate.Credentials

		for {
			cmd.Printf("Refreshing state")
			refresh()
			cmd.Printf(" done\n\n")

			credentialsToAction = make(cstate.Credentials, 0)
			credentialsToDeploy = make(cstate.Credentials, 0)

			credentials := state.Credentials(filters.Filters()...)
			credentials.SortByNameAndCreatedAt()

			for _, cred := range credentials {
				switch action := cred.NextAction(regenerationCriteria); {
				case action == cstate.BoshDeploy:
					credentialsToDeploy = append(credentialsToDeploy, cred)
				case action == cstate.NoOverwrite:
					continue
				case action == cstate.None:
					continue
				default:
					credentialsToAction = append(credentialsToAction, cred)
				}
			}

			if len(credentialsToAction) == 0 {
				cmd.Printf("No further actions to perform\n\n")
				break
			} else {
				cmd.Printf("Perform actions:\n")

				for _, cred := range credentialsToAction {
					cmd.Printf("- %s %s\n  L %s\n",
						cred.NextAction(regenerationCriteria).String(), cred.PathVersion(), cred.Summary())
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
				cmd.Println("")
			}
		}

		if len(credentialsToDeploy) != 0 {
			cmd.Printf("Found credential(s) pending a bosh deploy:\n")
			for _, cred := range credentialsToDeploy {
				cmd.Printf("- bosh_deploy(%s) %s\n  L %s\n",
					cred.PendingDeploys().String(), cred.PathVersion(), cred.Summary())
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(rotateCmd)

	addExpiresWithinCriteriaFlag(rotateCmd.Flags())
	addOlderThanCireteriaFlag(rotateCmd.Flags())
	addIgnoreUpdateModeCireteriaFlag(rotateCmd.Flags())
	addNameFlag(rotateCmd.Flags())
	addDeploymentFlag(rotateCmd.Flags())
}
