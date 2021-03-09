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
			path = &Path{
				Name:        cred.Name,
				Deployments: make(Deployments, 0),
			}
			s.paths.Put(cred.Name, path)
		}

		c := Credential{
			Credential:  cred,
			Deployments: make(Deployments, 0),
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
		}
		for _, ca := range cert.Ca {
			ca, found := s.getCredentialBySubjectKeyId(ca.SubjectKeyId)
			if found {
				cert.References = append(cert.References, ca)
				ca.ReferencedBy = append(ca.ReferencedBy, cert)
			}
		}
	}

	// Sort Credentials
	s.eachPath(func(p *Path) {
		p.Versions.SortByCreatedAt()
	})
	s.eachDeployment(func(d *Deployment) {
		d.Versions.SortByCreatedAt()
	})
	s.eachCredential(func(c *Credential) {
		c.Signs.SortByCreatedAt()
	})

	// Mark last Credential per Path as Latest
	s.eachPath(func(p *Path) {
		// There can never be a path without at least one version
		// slice already sorted above
		p.Versions[0].Latest = true
		for _, c := range p.Versions {
			if len(c.Signs) != 0 && !c.Transitional {
				signing := true
				c.Signing = &signing
				break
			}
		}
	})

	for _, variable := range variables {
		d := s.getOrCreateDeployment(variable.Deployment)
		credential, found := s.getCredential(variable.ID)
		if !found {
			return fmt.Errorf(`credential not found for bosh variable id: %s
run bosh deploy for '%s' so the BOSH director can converge it's variable references`,
				variable.ID, variable.Deployment)
		}
		credential.Deployments = append(credential.Deployments, d)
		d.Versions = append(d.Versions, credential)

		path, found := s.getPath(variable.Name)
		if !found {
			return fmt.Errorf(`path not found for bosh variable name: %s
run bosh deploy for '%s' so the BOSH director can converge it's variable references`,
				variable.Name, variable.Deployment)
		}

		path.VariableDefinition = variable.Definition

		if !path.Deployments.Includes(d) {
			path.Deployments = append(path.Deployments, d)
		}
	}

	return nil
}
