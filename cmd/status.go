/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless reqbrowsered by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	s "github.com/starkandwayne/carousel/state"
)

// statusCmd represents the browse command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get status of all credentials",
	Long:  `Display status of credentials.`,
	Run: func(cmd *cobra.Command, args []string) {
		initialize()
		refresh()

		fmt.Printf("%+v", filters)
		for _, path := range state.Paths() {
			pathState := PathState{
				path:             path,
				deployments:      s.Deployments{},
				needsRedeploying: s.Deployments{},
			}

			for _, filter := range filters.Filters() {
				if filter(path.Versions[0]) {
					pathState.needsRegenerating = true
				}
			}
			latestDeployments := path.Versions[0].Deployments

			for _, version := range path.Versions {
				for _, deployment := range version.Deployments {
					pathState.deployments = append(pathState.deployments, deployment)
					if !contains(latestDeployments, deployment) {
						pathState.needsRedeploying = append(pathState.needsRedeploying, deployment)
					}
				}
			}
			pathState.printState()
		}
		logger.Printf("%+v", filters)
	},
}

type PathState struct {
	path              *s.Path
	needsRegenerating bool
	needsRedeploying  s.Deployments
	deployments       s.Deployments
}

func (ps *PathState) printState() {
	fmt.Printf("PATH: %v\n", ps.path.Name)
	fmt.Printf("Needs Regenerating: %v\n", ps.needsRegenerating)
	fmt.Printf("All Deployments:\n")
	fmt.Printf("%v\n", ps.deployments.String())
	fmt.Printf("Needs Redeploying:\n")
	fmt.Printf("%v\n", ps.needsRedeploying.String())
}

func contains(s s.Deployments, e *s.Deployment) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func init() {
	rootCmd.AddCommand(statusCmd)
	addOlderThanFlag(statusCmd.Flags())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// browseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// browseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
