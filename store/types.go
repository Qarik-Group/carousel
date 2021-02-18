package store

import (
	"time"

	"github.com/starkandwayne/carousel/credhub"
)

type Path struct {
	Name               string
	Versions           []*Credential
	VariableDefinition *VariableDefinition
}

func (p *Path) AppendVersion(c *Credential) {
	p.Versions = append(p.Versions, c)
}

type Credential struct {
	*credhub.Credential
	Deployments []*Deployment
	SignedBy    *Credential
	Signs       []*Credential
	Path        *Path
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
	Versions []*Credential
	Name     string
}

type VariableDefinition struct {
	Name    string      `yaml:"name"`
	Type    string      `yaml:"type"`
	Options interface{} `yaml:"options,omitempty"`
}
