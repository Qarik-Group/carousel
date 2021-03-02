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
	cstate "github.com/starkandwayne/carousel/state"
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
		var preRegenerateHook func(*cobra.Command, cstate.Credentials)

		switch cmd.Parent() {
		case caCertificatesCmd:
			filters.ca = true
			typeSingular = ccredhub.Certificate.String()
		case leafCertificatesCmd:
			filters.leaf = true
			typeSingular = ccredhub.Certificate.String()
			preRegenerateHook = preRegenerateLeaf
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

		if preRegenerateHook != nil {
			preRegenerateHook(cmd, credentials)
		}

		cmd.Printf("\nPerforming credential regeneration:\n")
		for _, cred := range credentials {
			cmd.Printf("- %s", cred.Name)
			if regenerateForceFlag || cred.Generated {
				err := credhub.ReGenerate(cred.Credential)
				if err != nil {
					cmd.Printf(" got error: %s\n", err)
					os.Exit(1)
				}
				cmd.Print(" done\n")
			} else {
				cmd.Print(" skipping: since it was not genearted by CredHub\n")
			}
		}
		cmd.Println("Finished")
	},
}

func init() {
	addCommonFlags := func(f *pflag.FlagSet) {
		addDeploymentFlag(f)
		addNameFlag(f)
		f.BoolVar(&regenerateForceFlag, "force", false,
			"Regenerate both CredHub-generated and manually set credentials.")
	}

	// SSH, RSA, Users, Passwords
	var defaultRegenerateCmd = *regenerateCmd
	addCommonFlags(defaultRegenerateCmd.Flags())
	addOlderThanFlag(defaultRegenerateCmd.Flags())

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
	addCommonFlags(certificatesRegenerateCmd.Flags())
	addExpiresWithinFlag(certificatesRegenerateCmd.Flags())
	addSignedByFlag(certificatesRegenerateCmd.Flags())

	var caCertificatesRegenerateCmd = certificatesRegenerateCmd
	caCertificatesCmd.AddCommand(&caCertificatesRegenerateCmd)

	var leafCertificatesRegenerateCmd = certificatesRegenerateCmd
	leafCertificatesCmd.AddCommand(&leafCertificatesRegenerateCmd)
}

func preRegenerateLeaf(cmd *cobra.Command, credentials cstate.Credentials) {
	transitionalCAs := credentials.Collect(
		cstate.SignedByCollector()).
		Collect(cstate.SibilingsCollector()).
		Unique().Select(cstate.TransitionalFilter())

	if transitionalCAs.Any() {
		logger.Print("\nVerifying all transitional CA's have been deployed:\n")

		ok := true
		for _, tca := range transitionalCAs {
			logger.Printf("- %s   %s\n", tca.ID, tca.Name)
			for _, leaf := range tca.Path.Versions.
				Collect(cstate.SignsCollector()).
				Unique().
				Select(cstate.ActiveFilter()) {
				check := "✓"
				if !leaf.References.Includes(tca) {
					ok = false
					check = "×"
				}
				logger.Printf("L %s %s %s deployments [%s]\n", check, leaf.ID, leaf.Name, leaf.Deployments.String())
			}
		}
		if !ok {
			logger.Fatal("Failed: not all transitional CA's have been deployed.\nPlease issue a deploy for the deployments marked above")
		} else {
			logger.Print("All clear\n\n")
		}

		logger.Print("Moving transitional CA flag to signing versions\n")
		for _, signingCA := range transitionalCAs.Collect(cstate.SibilingsCollector()).Select(cstate.SigningFilter()) {
			cmd.Printf("- %s", signingCA.Name)
			err := credhub.UpdateTransitional(signingCA.Credential, false)
			if err != nil {
				cmd.Printf(" got error: %s\n", err)
				os.Exit(1)
			}
			cmd.Print(" done\n")
		}
		cmd.Printf("Finished\n\n")
	}

}
