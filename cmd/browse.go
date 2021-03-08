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
	"github.com/spf13/cobra"

	"github.com/starkandwayne/carousel/app"
)

// browseCmd represents the browse command
var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "Browse credentials using an interacive terminal UI",
	Long: `To make it easier to debug credential, and in particular certificate issues, 
carousel provides an interactive terminal UI. Which gives the user an simpel
way of browsing trought certificate signing chains in a tree like fashion.`,
	Run: func(cmd *cobra.Command, args []string) {
		initialize()
		refresh()

		app := app.NewApplication(state, credhub, refresh).Init()

		if err := app.Run(); err != nil {
			logger.Fatalf("the browse encountered an error: %s", err)
		}

	},
}

func init() {
	rootCmd.AddCommand(browseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// browseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// browseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
