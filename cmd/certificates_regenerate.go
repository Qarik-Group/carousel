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

	ccredhub "github.com/starkandwayne/carousel/credhub"
)

var regenerateForceFlag bool

// regenerateCmd represents the regenerate command
var regenerateCmd = &cobra.Command{
	Use:   "regenerate",
	Short: "Regenerate certificates",
	Long: `Regenerates CredHub-generated certificates. 
By default, certificates that have been manually set in CredHub are not regenerated.`,
	Run: func(cmd *cobra.Command, args []string) {
		initialize()
		refresh()

		filters.latest = true
		filters.types = []string{ccredhub.Certificate.String()}

		switch cmd.Parent() {
		case caCertificatesCmd:
			filters.ca = true
		case leafCertificatesCmd:
			filters.leaf = true
		}

		credentials := state.Credentials(filters.Filters()...)

		if len(credentials) == 0 {
			cmd.Println("No Certificates match criteria, nothing to do")
			os.Exit(0)
		}

		cmd.Println("Regenerating Certificates:")
		for _, cred := range credentials {
			cmd.Printf("- %s\n", cred.Name)
		}
		askForConfirmation()

		for _, cred := range state.Credentials(filters.Filters()...) {
			credhub.ReGenerate(cred.Credential, regenerateForceFlag)
		}
		cmd.Println("Done")
	},
}

func init() {
	addDeploymentFlag(regenerateCmd.Flags())
	addNameFlag(regenerateCmd.Flags())
	addExpiresWithinFlag(regenerateCmd.Flags())
	addSignedByFlag(regenerateCmd.Flags())

	regenerateCmd.Flags().BoolVar(&regenerateForceFlag, "force", false,
		"Regenerate both CredHub-generated and manually set Certificates.")

	var caCertificatesRegenerateCmd = *regenerateCmd
	var leafCertificatesRegenerateCmd = *regenerateCmd

	caCertificatesCmd.AddCommand(&caCertificatesRegenerateCmd)
	leafCertificatesCmd.AddCommand(&leafCertificatesRegenerateCmd)
}
