package state

import (
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

func (d Deployments) IncludesName(name string) bool {
	for _, deployment := range d {
		if deployment.Name == name {
			return true
		}
	}
	return false
}
