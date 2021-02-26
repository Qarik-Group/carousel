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
	Short: "Regenerate credentials",
	Long:  `Regenerates CredHub-generated credentials`,
	Run: func(cmd *cobra.Command, args []string) {
		initialize()
		refresh()

		filters.latest = true
		var typeSingular string

		switch cmd.Parent() {
		case caCertificatesCmd:
			filters.ca = true
			typeSingular = ccredhub.Certificate.String()
		case leafCertificatesCmd:
			filters.leaf = true
			typeSingular = ccredhub.Certificate.String()
		case sshKeyPairsCmd:
			typeSingular = ccredhub.SSH.String()
		case rsaKeyPairsCmd:
			typeSingular = ccredhub.RSA.String()
		case usersCmd:
			typeSingular = ccredhub.User.String()
		case passwordsCmd:
			typeSingular = ccredhub.Password.String()
		}

		filters.types = []string{typeSingular}

		credentials := state.Credentials(filters.Filters()...)

		if len(credentials) == 0 {
			cmd.Printf("No %s credentials match criteria, nothing to do\n", typeSingular)
			os.Exit(0)
		}

		cmd.Printf("Regenerating %s credentials:\n", typeSingular)
		for _, cred := range credentials {
			cmd.Printf("- %s\n", cred.Name)
		}
		askForConfirmation()

		for _, cred := range state.Credentials(filters.Filters()...) {
			if regenerateForceFlag || cred.Generated {
				cmd.Printf("- %s", cred.Name)
				err := credhub.ReGenerate(cred.Credential)
				if err != nil {
					cmd.Printf(" got error: %s", cred.Name)
					os.Exit(1)
				}
				cmd.Print(" done\n")
			} else {
				cmd.Println("skipping: %s, since it was not genearted by")
			}
		}
		cmd.Println("Finished")
	},
}

func init() {
	// Common
	addDeploymentFlag(regenerateCmd.Flags())
	addNameFlag(regenerateCmd.Flags())
	regenerateCmd.Flags().BoolVar(&regenerateForceFlag, "force", false,
		"Regenerate both CredHub-generated and manually set credentials.")

	// SSH, RSA, Users, Passwords
	var defaultRegenerateCmd = *regenerateCmd
	addOlderThanFlag(defaultRegenerateCmd.LocalFlags())

	var sshKeyPairsRegenerateCmd = defaultRegenerateCmd
	sshKeyPairsCmd.AddCommand(&sshKeyPairsRegenerateCmd)

	var rsaKeyPairsRegenerateCmd = defaultRegenerateCmd
	rsaKeyPairsCmd.AddCommand(&rsaKeyPairsRegenerateCmd)

	var usersRegenerateCmd = defaultRegenerateCmd
	usersCmd.AddCommand(&usersRegenerateCmd)

	var passwordsRegenerateCmd = defaultRegenerateCmd
	passwordsCmd.AddCommand(&passwordsRegenerateCmd)

	// Certificates
	var certificatesRegenerateCmd = *regenerateCmd
	addExpiresWithinFlag(certificatesRegenerateCmd.LocalFlags())
	addSignedByFlag(certificatesRegenerateCmd.LocalFlags())

	var caCertificatesRegenerateCmd = certificatesRegenerateCmd
	caCertificatesCmd.AddCommand(&caCertificatesRegenerateCmd)

	var leafCertificatesRegenerateCmd = certificatesRegenerateCmd
	leafCertificatesCmd.AddCommand(&leafCertificatesRegenerateCmd)
}
