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

	"github.com/karrick/tparse"
	"github.com/spf13/cobra"
	"time"

	cstate "github.com/starkandwayne/carousel/state"
)

// statusCmd represents the status command
var rotateCmd = &cobra.Command{
	Use:   "rotate",
	Short: "Rotate all credentials needing rotation",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		initialize()
		refresh()

		ew, err := tparse.AddDuration(time.Now(), "+"+expiresWithin)
		if err != nil {
			logger.Fatalf("failed to parse duration: %s, got: %s",
				expiresWithin, err)
		}

		ot, err := tparse.AddDuration(time.Now(), "-"+olderThan)
		if err != nil {
			logger.Fatalf("failed to parse duration: %s, got: %s",
				olderThan, err)
		}

		criteria := cstate.RegenerationCriteria{
			OlderThan:     ot,
			ExpiresBefore: ew,
		}

		credentials := state.Credentials()
		credentials.SortByNameAndCreatedAt()

		credentialsToRotate := []*cstate.Credential{}

		cmd.Printf("Rotating credentials:\n")

		for _, cred := range credentials {
			switch action := cred.NextAction(criteria); {
			case action == cstate.Regenerate:
				cmd.Printf("%s, ", cred.PathVersion())
				credentialsToRotate = append(credentialsToRotate, cred)
				continue
			}
		}

		askForConfirmation()

		cmd.Printf("\nPerforming credential rotation:\n")

		for _, cred := range credentialsToRotate {
			cmd.Printf("- %s", cred.Name)
			err := credhub.ReGenerate(cred.Credential)
			if err != nil {
				cmd.Printf(" got error: %s\n", err)
				os.Exit(1)
			}
			cmd.Print(" done\n")
		}

		cmd.Println("Finished")
	},
}

func init() {
	rootCmd.AddCommand(rotateCmd)
	rotateCmd.Flags().StringVar(&expiresWithin, "expires-within", "3m",
		"filter certificates by expiry window (suffixes: d day, w week, y year)")

	rotateCmd.Flags().StringVar(&olderThan, "older-than", "1y",
		"filter credentials by age (suffixes: d day, w week, y year)")
}
