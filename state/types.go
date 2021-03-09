package state

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/starkandwayne/carousel/bosh"
	"github.com/starkandwayne/carousel/credhub"
)

type Path struct {
	Name               string                   `json:"name"`
	Versions           Credentials              `json:"-"`
	VariableDefinition *bosh.VariableDefinition `json:"variable_definition"`
	Deployments        Deployments
}

type Deployment struct {
	Versions Credentials `json:"-"`
	Name     string      `json:"name"`
}

type Credential struct {
	*credhub.Credential
	Deployments  Deployments `json:"-"`
	SignedBy     *Credential `json:"-"`
	ReferencedBy Credentials `json:"-"`
	References   Credentials `json:"-"`
	Signs        Credentials `json:"-"`
	Latest       bool        `json:"latest"`
	Signing      *bool       `json:"signing,omitempty"`
	Path         *Path       `json:"-"`
}

type Credentials []*Credential

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

func (c *Credential) PathVersion() string {
	return fmt.Sprintf("%s@%s", c.Name, c.ID)
}

func (c *Credential) PendingDeploys() Deployments {
	out := make(Deployments, 0)
	if !c.Latest {
		return out
	}
	for _, d := range c.Path.Deployments {
		if !c.Deployments.Includes(d) {
			out = append(out, d)
		}
	}
	return out
}

func (c *Credential) Active() bool {
	if len(c.Deployments) != 0 {
		return true
	}
	for _, cred := range c.ReferencedBy {
		if cred.Active() {
			return true
		}
	}
	return false
}

type Deployments []*Deployment

func (d Deployments) String() string {
	tmp := make([]string, 0, len(d))
	for _, deployment := range d {
		tmp = append(tmp, deployment.Name)
	}
	return strings.Join(tmp, ", ")
}

func (d Deployments) Includes(this *Deployment) bool {
	for _, deployment := range d {
		if deployment == this {
			return true
		}
	}
	return false
}
