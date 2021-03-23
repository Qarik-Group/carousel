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
	"github.com/cloudboss/ofcourse/ofcourse"
	"github.com/spf13/cobra"
	"github.com/starkandwayne/carousel/resource"
)

var concourseCmd = &cobra.Command{
	Use:    "concourse",
	Short:  "Subcommand to execute concourse resource operations",
	Hidden: true,
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "The check command for the concourse resource",
	Run: func(cmd *cobra.Command, args []string) {
		ofcourse.Check(&resource.Resource{})
	},
}
var inCmd = &cobra.Command{
	Use:   "ci_in",
	Short: "The in command for the concourse resource",
	Run: func(cmd *cobra.Command, args []string) {
		ofcourse.In(&resource.Resource{})
	},
}

var outCmd = &cobra.Command{
	Use:   "ci_out",
	Short: "The out command for the concourse resource",
	Run: func(cmd *cobra.Command, args []string) {
		ofcourse.Out(&resource.Resource{})
	},
}

func init() {
	rootCmd.AddCommand(concourseCmd)
	concourseCmd.AddCommand(checkCmd)
	concourseCmd.AddCommand(inCmd)
	concourseCmd.AddCommand(outCmd)
}
