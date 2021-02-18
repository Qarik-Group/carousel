package store

import (
	"github.com/emirpasic/gods/maps/treebidimap"
	"github.com/starkandwayne/carousel/credhub"

	boshdir "github.com/cloudfoundry/bosh-cli/director"
)

type Store struct {
	deployments         *treebidimap.Map
	paths               *treebidimap.Map
	credentials         *treebidimap.Map
	variableDefinitions *treebidimap.Map
	credhub             credhub.CredHub
	directorClient      boshdir.Director
}

func (s *Store) GetPath(path string) (*Path, bool) {
	p, found := s.paths.Get(path)
	if !found {
		return nil, found
	}
	return p.(*Path), true
}

func (s *Store) EachPath(fn func(*Path)) {
	s.paths.Each(func(_, v interface{}) {
		fn(v.(*Path))
	})
}

func (s *Store) GetCredential(id string) (*Credential, bool) {
	c, found := s.credentials.Get(id)
	if !found {
		return nil, found
	}
	return c.(*Credential), true
}

func (s *Store) EachCredential(fn func(v *Credential)) {
	s.credentials.Each(func(_, v interface{}) {
		fn(v.(*Credential))
	})
}

func (s *Store) Certificates() []*Credential {
	certs := s.credentials.Select(func(_, v interface{}) bool {
		return v.(*Credential).Type == credhub.Certificate
	})
	out := make([]*Credential, 0)
	for _, cert := range certs.Values() {
		out = append(out, cert.(*Credential))
	}
	return out
}

func pathByName(a, b interface{}) int {
	// Type assertion, program will panic if this is not respected
	c1 := a.(*Path)
	c2 := b.(*Path)

	switch {
	case c1.Name > c2.Name:
		return 1
	case c1.Name < c2.Name:
		return -1
	default:
		return 0
	}
}

func credentialById(a, b interface{}) int {
	// Type assertion, program will panic if this is not respected
	c1 := a.(*Credential)
	c2 := b.(*Credential)

	switch {
	case c1.ID > c2.ID:
		return 1
	case c1.ID < c2.ID:
		return -1
	default:
		return 0
	}
}

func deploymentByName(a, b interface{}) int {
	// Type assertion, program will panic if this is not respected
	c1 := a.(*Deployment)
	c2 := b.(*Deployment)

	switch {
	case c1.Name > c2.Name:
		return 1
	case c1.Name < c2.Name:
		return -1
	default:
		return 0
	}
}

func veriableDefinitionByName(a, b interface{}) int {
	// Type assertion, program will panic if this is not respected
	c1 := a.(*VariableDefinition)
	c2 := b.(*VariableDefinition)

	switch {
	case c1.Name > c2.Name:
		return 1
	case c1.Name < c2.Name:
		return -1
	default:
		return 0
	}
}
