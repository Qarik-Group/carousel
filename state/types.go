package state

import (
	"encoding/json"

	"github.com/starkandwayne/carousel/bosh"
	"github.com/starkandwayne/carousel/credhub"
)

type Path struct {
	Name               string                   `json:"name"`
	Versions           []*Credential            `json:"-"`
	VariableDefinition *bosh.VariableDefinition `json:"variable_definition"`
}

type Deployment struct {
	Versions []*Credential `json:"-"`
	Name     string        `json:"name"`
}

type Credential struct {
	*credhub.Credential
	Deployments []*Deployment `json:"-"`
	SignedBy    *Credential   `json:"-"`
	Signs       []*Credential `json:"-"`
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
