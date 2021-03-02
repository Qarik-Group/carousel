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
	"github.com/spf13/pflag"

	ccredhub "github.com/starkandwayne/carousel/credhub"
)

var cleanUpForceFlag bool

// cleanUpCmd represents the clean-up command
var cleanUpCmd = &cobra.Command{
	Use:   "clean-up",
	Short: "Clean-up unused credentials",
	Long:  `Clean-up credentials not used by any BOSH deployment`,
	Run: func(cmd *cobra.Command, args []string) {
		initialize()
		refresh()

		filters.unused = true

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

		cmd.Printf("Cleaning up %s credentials:\n", typeSingular)
		for _, cred := range credentials {
			cmd.Printf("- %s@%s\n", cred.Name, cred.ID)
		}
		askForConfirmation()

		for _, cred := range credentials {
			cmd.Printf("- %s@%s", cred.Name, cred.ID)
			err := credhub.Delete(cred.Credential)
			if err != nil {
				cmd.Printf(" got error: %s\n", err)
				os.Exit(1)
			}
			cmd.Print(" removed\n")
		}
		cmd.Println("Finished")
	},
}

func init() {
	addCommonFlags := func(f *pflag.FlagSet) {
		addNameFlag(f)
	}

	// Currently only supported by CredHub to delete versions of certificates
	// SSH, RSA, Users, Passwords
	// var defaultCleanUpCmd = *cleanUpCmd
	// addCommonFlags(defaultCleanUpCmd.Flags())
	// addOlderThanFlag(defaultCleanUpCmd.Flags())

	// var sshKeyPairsCleanUpCmd = defaultCleanUpCmd
	// sshKeyPairsCmd.AddCommand(&sshKeyPairsCleanUpCmd)

	// var rsaKeyPairsCleanUpCmd = defaultCleanUpCmd
	// rsaKeyPairsCmd.AddCommand(&rsaKeyPairsCleanUpCmd)

	// var usersCleanUpCmd = defaultCleanUpCmd
	// usersCmd.AddCommand(&usersCleanUpCmd)

	// var passwordsCleanUpCmd = defaultCleanUpCmd
	// passwordsCmd.AddCommand(&passwordsCleanUpCmd)

	// Certificates
	var certificatesCleanUpCmd = *cleanUpCmd
	addCommonFlags(certificatesCleanUpCmd.Flags())
	addSignedByFlag(certificatesCleanUpCmd.Flags())

	var caCertificatesCleanUpCmd = certificatesCleanUpCmd
	caCertificatesCmd.AddCommand(&caCertificatesCleanUpCmd)

	var leafCertificatesCleanUpCmd = certificatesCleanUpCmd
	leafCertificatesCmd.AddCommand(&leafCertificatesCleanUpCmd)
}
