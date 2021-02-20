package store

import (
	"encoding/json"
	"time"

	"github.com/starkandwayne/carousel/credhub"
)

type Path struct {
	Name               string              `json:"name"`
	Versions           []*Credential       `json:"-"`
	VariableDefinition *VariableDefinition `json:"variable_definition"`
}

func (p *Path) AppendVersion(c *Credential) {
	p.Versions = append(p.Versions, c)
}

type Credential struct {
	*credhub.Credential
	Deployments []*Deployment `json:"-"`
	SignedBy    *Credential   `json:"-"`
	Signs       []*Credential `json:"-"`
	Path        *Path         `json:"-"`
}

func (c *Credential) MarshalJSON() ([]byte, error) {
	deployments := make([]string, 0)
	for _, d := range c.Deployments {
		deployments = append(deployments, d.Name)
	}

	updateMode := NoOverwrite
	if c.Path.VariableDefinition != nil {
		updateMode = c.Path.VariableDefinition.UpdateMode
	}

	type Alias Credential
	return json.Marshal(&struct {
		*Alias
		DeploymentsList []string   `json:"deployments"`
		UpdateMode      UpdateMode `json:"update_mode"`
	}{
		Alias:           (*Alias)(c),
		DeploymentsList: deployments,
		UpdateMode:      updateMode,
	})
}

func (c *Credential) Status() string {
	status := "active"
	if c.ExpiryDate != nil && c.ExpiryDate.Sub(time.Now()) < time.Hour*24*30 {
		status = "notice"
	}
	if c.VersionCreatedAt.Sub(time.Now()) > time.Hour*24*365 {
		status = "notice"
	}
	if len(c.Deployments) == 0 {
		status = "unused"
	}
	return status
}

type Deployment struct {
	Versions []*Credential `json:"-"`
	Name     string        `json:"name"`
}

type VariableDefinition struct {
	Name       string                 `yaml:"name" json:"name"`
	Type       string                 `yaml:"type" json:"type"`
	UpdateMode UpdateMode             `yaml:"update_mode,omitempty" json:"update_mode,omitempty"`
	Options    map[string]interface{} `yaml:"options,omitempty" json:"options,omitempty"`
}

type UpdateMode string

const (
	NoOverwrite, Overwrite, Converge UpdateMode = "no-overwrite", "overwrite", "converge"
)

func (v *VariableDefinition) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// update_mode [String, optional]: Update mode to use when generating credentials.
	// Currently supported update modes are no-overwrite, overwrite, and converge. Defaults to no-overwrite
	// https://bosh.io/docs/manifest-v2/#variables

	type VariableDefinitionDefaulted VariableDefinition
	var defaults = VariableDefinitionDefaulted{
		UpdateMode: NoOverwrite,
	}

	out := defaults
	err := unmarshal(&out)
	*v = VariableDefinition(out)
	return err
}
