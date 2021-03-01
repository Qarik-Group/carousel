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
	"fmt"
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

		for _, cred := range credentials {
			cmd.Printf("- %s", cred.Name)
			if regenerateForceFlag || cred.Generated {
				err := rotateCredential(cred)
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

func rotateCredential(c *cstate.Credential) error {
	switch c.Credential.Type {
	case ccredhub.Certificate:
		if c.Certificate.IsCA {
			return rotateCa(c)
		} else {
			return rotateLeaf(c)
		}

	default:
		return credhub.ReGenerate(c.Credential)
	}
}

func rotateCa(c *cstate.Credential) error {
	return credhub.ReGenerate(c.Credential)
}

func rotateLeaf(c *cstate.Credential) error {
	latestCA := c.SignedBy.Path.Versions.LatestVersion()
	if latestCA.Transitional {
		activeLeafsNotSignedByLatestCA := make(cstate.Credentials, 0)
		for _, sibling := range c.SignedBy.Signs {
			for _, activeSiblingVersion := range sibling.Path.Versions.ActiveVersions() {
				if !activeSiblingVersion.CAs.Includes(latestCA) &&
					!activeLeafsNotSignedByLatestCA.Includes(activeSiblingVersion) {
					activeLeafsNotSignedByLatestCA = append(
						activeLeafsNotSignedByLatestCA,
						activeSiblingVersion)
				}
			}
		}

		if len(activeLeafsNotSignedByLatestCA) != 0 {
			outErr := "CA check failed, got following failures:"
			for _, leaf := range activeLeafsNotSignedByLatestCA {
				outErr = fmt.Sprintf("%s\n- %s@%s (used by: %s) .ca field did not referenced: %s@%s",
					outErr, leaf.Path.Name, leaf.ID, leaf.Deployments.String(), latestCA.Path.Name, latestCA.ID)
			}
			return fmt.Errorf(outErr)
		}
	}

	err := credhub.UpdateTransitional(c.SignedBy.Path.Versions.SigningVersion().Credential)
	if err != nil {
		return err
	}
	return credhub.ReGenerate(c.Credential)
}
