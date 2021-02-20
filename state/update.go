package state

import (
	"fmt"

	"github.com/starkandwayne/carousel/bosh"
	"github.com/starkandwayne/carousel/credhub"
)

func (s *state) Update(credentials []*credhub.Credential, variables []*bosh.Variable) error {
	s.clear()

	for _, cred := range credentials {
		var path *Path
		p, found := s.paths.Get(cred.Name)
		if found {
			path = p.(*Path)
		} else {
			path = &Path{Name: cred.Name}
			s.paths.Put(cred.Name, path)
		}

		c := Credential{
			Credential:  cred,
			Deployments: make([]*Deployment, 0),
			Path:        path,
		}

		path.Versions = append(path.Versions, &c)
		s.credentials.Put(cred.ID, &c)
	}

	// Lookup Ca for each cert
	for _, cert := range s.Credentials(TypeFilter(credhub.Certificate)) {
		authorityKeyID := cert.Certificate.AuthorityKeyId
		if cert.SelfSigned {
			continue
		}
		ca, found := s.getCredentialBySubjectKeyId(authorityKeyID)
		if found {
			ca.Signs = append(ca.Signs, cert)
			cert.SignedBy = ca
		} else {
			return fmt.Errorf("failed to lookup ca Credential with id: %s", cert.ID)
		}
	}

	// Mark last Credential per Path as Latest
	s.eachPath(func(p *Path) {
		latest := p.Versions[0] // There can never be a path without at least one version
		var signing *Credential
		for _, c := range p.Versions {
			if latest.VersionCreatedAt.Before(*c.VersionCreatedAt) {
				latest = c
			}
			if c.CertificateAuthority && len(c.Signs) != 0 && len(c.Deployments) != 0 {
				if signing == nil || signing.VersionCreatedAt.Before(*c.VersionCreatedAt) {
					signing = c
				}
			}
		}
		if signing != nil {
			*signing.Signing = true
		}
		latest.Latest = true
	})

	for _, variable := range variables {
		d := s.getOrCreateDeployment(variable.Deployment)
		credential, found := s.getCredential(variable.ID)
		if !found {
			return fmt.Errorf("failed to lookup credential for bosh variable with id: %s",
				variable.ID)
		}
		credential.Deployments = append(credential.Deployments, d)
		d.Versions = append(d.Versions, credential)

		path, found := s.getPath(variable.Name)
		if !found {
			return fmt.Errorf("failed to lookup path for bosh variable with name: %s",
				variable.Name)
		}

		path.VariableDefinition = variable.Definition
	}

	return nil
}
