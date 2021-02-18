package store

import (
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
	Deployments []*Deployment `json:"deployments"`
	SignedBy    *Credential   `json:"-"`
	Signs       []*Credential `json:"-"`
	Path        *Path         `json:"path"`
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
	Name    string                 `yaml:"name" json:"name"`
	Type    string                 `yaml:"type" json:"type"`
	Options map[string]interface{} `yaml:"options,omitempty" json:"options,omitempty"`
}
