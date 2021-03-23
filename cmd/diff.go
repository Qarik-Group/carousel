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
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path"
	"sort"

	"gopkg.in/yaml.v3"

	"github.com/gonvenience/ytbx"
	"github.com/homeport/dyff/pkg/dyff"
	"github.com/spf13/cobra"
	cstate "github.com/starkandwayne/carousel/state"

	boshtpl "github.com/cloudfoundry/bosh-cli/director/template"
)

var doNotInspectCerts bool

// statusCmd represents the status command
var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show diff off what should be deployed",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		initialize()

		if filters.deployment == "" {
			logger.Fatal("deployment flag must be set")
		}

		refresh()

		latest := state.Credentials(append(filters.Filters(), cstate.LatestFilter())...)
		active := state.Credentials(append(filters.Filters(), cstate.ActiveFilter())...)

		manfest, err := director.GetManifest(filters.deployment)
		if err != nil {
			logger.Fatalf("failed to get bosh manifest: %s", err)
		}

		latestYAML, err := renderTemplate(manfest, latest)
		if err != nil {
			logger.Fatalf("failed to build latest yaml: %s", err)
		}
		latestNames := []string{"manifest"}

		activeYAML, err := renderTemplate(manfest, active)
		if err != nil {
			logger.Fatalf("failed to build active yaml: %s", err)
		}
		activeNames := []string{"manifest"}

		appendConfigs(director.GetLatestCloudConfigs, latest, &latestYAML, &latestNames)
		appendConfigs(director.GetActiveCloudConfigs, active, &activeYAML, &activeNames)
		appendConfigs(director.GetLatestRuntimeConfigs, latest, &latestYAML, &latestNames)
		appendConfigs(director.GetActiveRuntimeConfigs, active, &activeYAML, &activeNames)

		report, err := dyff.CompareInputFiles(ytbx.InputFile{
			Documents: activeYAML,
			Names:     activeNames,
		}, ytbx.InputFile{
			Documents: latestYAML,
			Names:     latestNames,
		})

		if len(report.Diffs) == 0 {
			cmd.Println("Nothing to deploy")
		} else {
			reportWriter := &dyff.HumanReport{
				Report:            report,
				DoNotInspectCerts: doNotInspectCerts,
				NoTableStyle:      false,
				OmitHeader:        true,
			}

			var buf bytes.Buffer
			out := bufio.NewWriter(&buf)
			reportWriter.WriteReport(out)
			out.Flush()
			fmt.Print(buf.String())

			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)

	addDeploymentFlag(diffCmd.Flags())
	diffCmd.Flags().BoolVar(&doNotInspectCerts, "do-not-inspect-certs", false,
		"don't show a human readable diff for certificates")
}

type cred struct {
	Type        string `yaml:"type"`
	Version     string `yaml:"version"`
	CreatedAt   string `yaml:"created_at"`
	Certificate string `yaml:"certificate",omitempty`
	Ca          string `yaml:"ca",omitempty`
	Expiry      string `yaml:"expiry",omitempty`
}

func renderTemplate(manifest []byte, creds cstate.Credentials) ([]*yaml.Node, error) {
	tpl := boshtpl.NewTemplate(manifest)
	staticVars := boshtpl.StaticVariables{}

	for _, cred := range creds {
		staticVars[cred.Name] = cred.Credential.ToStaticVariable()
		staticVars[path.Base(cred.Name)] = cred.Credential.ToStaticVariable()
	}

	evalOpts := boshtpl.EvaluateOpts{
		UnescapedMultiline: true,
		ExpectAllKeys:      false,
		ExpectAllVarsUsed:  false,
	}

	bytes, err := tpl.Evaluate(staticVars, nil, evalOpts)
	if err != nil {
		return nil, fmt.Errorf("could not evaluate template: %v", err)
	}

	return ytbx.LoadYAMLDocuments(bytes)
}

func appendConfigs(fn func(deployment string) (map[string][]byte, error),
	creds cstate.Credentials, appendTo *[]*yaml.Node, names *[]string) {
	confs, err := fn(filters.deployment)
	if err != nil {
		logger.Fatalf("failed to get latest cloud configs: %s", err)
	}

	sortedNames := make([]string, 0)
	for name, _ := range confs {
		sortedNames = append(sortedNames, name)
	}
	sort.Strings(sortedNames)

	sortedNamesMap := make(map[string]int, 0)
	for i, name := range sortedNames {
		sortedNamesMap[name] = i
	}

	confYAMLs := make([][]*yaml.Node, len(confs), len(confs))
	for name, conf := range confs {
		confYAML, err := renderTemplate(conf, creds)
		if err != nil {
			logger.Fatalf("failed to interpolate config variables into: %s, got: %s", name, err)
		}
		confYAMLs[sortedNamesMap[name]] = confYAML
	}

	for _, confYAML := range confYAMLs {
		*appendTo = append(*appendTo, confYAML...)
	}
	*names = append(*names, sortedNames...)
}
