package state

import (
	"encoding/json"

	"github.com/starkandwayne/carousel/bosh"
	"github.com/starkandwayne/carousel/credhub"
)

type Path struct {
	Name               string                   `json:"name"`
	Versions           Credentials              `json:"-"`
	VariableDefinition *bosh.VariableDefinition `json:"variable_definition"`
}

type Deployment struct {
	Versions Credentials `json:"-"`
	Name     string      `json:"name"`
}

type Credential struct {
	*credhub.Credential
	Deployments []*Deployment `json:"-"`
	SignedBy    *Credential   `json:"-"`
	Signs       Credentials   `json:"-"`
	Latest      bool          `json:"latest"`
	Signing     *bool         `json:"signing,omitempty"`
	Path        *Path         `json:"-"`
}

func (c *Credential) MarshalJSON() ([]byte, error) {
	deployments := make([]string, 0)
	for _, d := range c.Deployments {
		deployments = append(deployments, d.Name)
	}

	updateMode := bosh.NoOverwrite
	if c.Path.VariableDefinition != nil {
		updateMode = c.Path.VariableDefinition.UpdateMode
	}

	c.RawValue = nil // don't leak raw value
	type Alias Credential
	return json.Marshal(&struct {
		*Alias
		DeploymentsList []string        `json:"deployments"`
		UpdateMode      bosh.UpdateMode `json:"update_mode"`
	}{
		Alias:           (*Alias)(c),
		DeploymentsList: deployments,
		UpdateMode:      updateMode,
	})
}

type Credentials []*Credential

func (c Credentials) Len() int {
	return len(c)
}

func (c Credentials) Less(i, j int) bool {
	return c[i].VersionCreatedAt.After(*c[j].VersionCreatedAt)
}

func (c Credentials) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
