package bosh

import (
	"fmt"
	"path"
	"strconv"

	"gopkg.in/yaml.v2"
)

func (d *director) GetVariables() ([]*Variable, error) {
	directorInfo, err := d.client.Info()
	if err != nil {
		return nil, err
	}

	deployments, err := d.client.Deployments()
	if err != nil {
		return nil, err
	}

	variables := make(map[string]*Variable, 0)
	for _, deployment := range deployments {
		vars, err := deployment.Variables()
		if err != nil {
			return nil, err
		}
		for _, v := range vars {
			variables[v.Name] = &Variable{
				ID:         v.ID,
				Name:       v.Name,
				Deployment: deployment.Name(),
			}
		}
	}

	for _, deployment := range deployments {
		rawDeploymentManifest, err := deployment.Manifest()
		if err != nil {
			return nil, err
		}

		if err := addVariableDefinition(variables, rawDeploymentManifest, func(n string) string {
			return path.Join("/", directorInfo.Name, deployment.Name(), n)
		}); err != nil {
			return nil, err
		}

		configs, err := d.client.ListDeploymentConfigs(deployment.Name())
		if err != nil {
			return nil, err
		}

		for _, conf := range configs.GetConfigs() {
			if conf.Type == "runtime" {
				c, err := d.client.LatestConfigByID(strconv.Itoa(conf.Id))
				if err != nil {
					return nil, err
				}

				if err := addVariableDefinition(variables, c.Content, func(n string) string {
					return n
				}); err != nil {
					return nil, err
				}
			}
		}
	}

	out := make([]*Variable, 0, len(variables))
	for _, v := range variables {
		out = append(out, v)
	}

	return out, nil
}

func addVariableDefinition(variables map[string]*Variable, raw string, nameFn func(name string) string) error {
	tmpl := manifest{}

	err := yaml.Unmarshal([]byte(raw), &tmpl)
	if err != nil {
		return err
	}

	for _, varDef := range tmpl.Variables {
		name := nameFn(varDef.Name)
		v, found := variables[name]
		if !found {
			return fmt.Errorf("failed to lookup path for variable definiton with name: %s", name)
		}
		v.Definition = varDef
	}
	return nil
}

type manifest struct {
	Variables []*VariableDefinition `yaml:"variables"`
}
