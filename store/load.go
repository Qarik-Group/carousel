package store

import (
	"bytes"
	"fmt"
	"path"
	"strconv"

	"github.com/emirpasic/gods/maps/treebidimap"
	"github.com/emirpasic/gods/utils"
	"github.com/starkandwayne/carousel/credhub"
	"gopkg.in/yaml.v2"

	boshdir "github.com/cloudfoundry/bosh-cli/director"
)

func NewStore(ch credhub.CredHub, directorClient boshdir.Director) *Store {
	return &Store{
		paths:               treebidimap.NewWith(utils.StringComparator, pathByName),
		credentials:         treebidimap.NewWith(utils.StringComparator, credentialById),
		deployments:         treebidimap.NewWith(utils.StringComparator, deploymentByName),
		variableDefinitions: treebidimap.NewWith(utils.StringComparator, veriableDefinitionByName),
		credhub:             ch,
		directorClient:      directorClient,
	}
}

func (s *Store) Refresh() error {
	s.paths.Clear()
	s.credentials.Clear()
	s.deployments.Clear()
	s.variableDefinitions.Clear()

	creds, err := s.credhub.FindAll()
	if err != nil {
		return err
	}

	for _, cred := range creds {
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

		path.AppendVersion(&c)
		s.credentials.Put(cred.ID, &c)
	}

	// Lookup Ca for each cert
	for _, cert := range s.Certificates() {
		authorityKeyID := cert.Certificate.AuthorityKeyId
		if cert.SelfSigned {
			continue
		}
		ca, found := s.getCertVersionBySubjectKeyId(authorityKeyID)
		if found {
			ca.Signs = append(ca.Signs, cert)
			cert.SignedBy = ca
		} else {
			return fmt.Errorf("failed to lookup ca Credential with id: %s", cert.ID)
		}
	}

	directorInfo, err := s.directorClient.Info()
	if err != nil {
		return err
	}

	deployments, err := s.directorClient.Deployments()
	if err != nil {
		return err
	}
	for _, deployment := range deployments {
		d := Deployment{
			Name:     deployment.Name(),
			Versions: make([]*Credential, 0),
		}
		s.deployments.Put(d.Name, &d)
		variables, err := deployment.Variables()
		if err != nil {
			return err
		}
		for _, variable := range variables {
			credential, found := s.GetCredential(variable.ID)
			if !found {
				return fmt.Errorf("failed to lookup credential for bosh variable with id: %s",
					variable.ID)
			}
			credential.Deployments = append(credential.Deployments, &d)
			d.Versions = append(d.Versions, credential)
		}

		rawDeploymentManifest, err := deployment.Manifest()
		if err != nil {
			return err
		}

		varDefs, err := rawManifestToVariableDefinitions(rawDeploymentManifest)
		if err != nil {
			return err
		}

		for _, varDef := range varDefs {
			name := path.Join("/", directorInfo.Name, deployment.Name(), varDef.Name)
			s.variableDefinitions.Put(name, varDef)

			p, found := s.GetPath(name)
			if !found {
				return fmt.Errorf("failed to lookup path for variable definiton with name: %s",
					name)
			}
			p.VariableDefinition = varDef
		}

		configs, err := s.directorClient.ListDeploymentConfigs(d.Name)
		if err != nil {
			return err
		}

		for _, conf := range configs.GetConfigs() {
			if conf.Type == "runtime" {
				c, err := s.directorClient.LatestConfigByID(strconv.Itoa(conf.Id))
				if err != nil {
					return err
				}

				varDefs, err := rawManifestToVariableDefinitions(c.Content)
				if err != nil {
					return err
				}

				for _, varDef := range varDefs {
					s.variableDefinitions.Put(varDef.Name, varDef)

					p, found := s.GetPath(varDef.Name)
					if !found {
						return fmt.Errorf("failed to lookup path for variable definiton with name: %s",
							varDef.Name)

					}
					p.VariableDefinition = varDef
				}

			}
		}
	}

	return nil
}

func rawManifestToVariableDefinitions(raw string) ([]*VariableDefinition, error) {
	tmpl := manifest{}

	err := yaml.Unmarshal([]byte(raw), &tmpl)
	if err != nil {
		return nil, err
	}

	return tmpl.Variables, nil
}

type manifest struct {
	Variables []*VariableDefinition `yaml:"variables"`
}

func (s *Store) getCertVersionBySubjectKeyId(keyId []byte) (*Credential, bool) {
	_, foundValue := s.credentials.Find(func(index interface{}, value interface{}) bool {
		if value.(*Credential).Certificate != nil {
			return bytes.Compare(value.(*Credential).Certificate.SubjectKeyId, keyId) == 0
		}
		return false
	})
	if foundValue != nil {
		return foundValue.(*Credential), true
	}
	return nil, false
}
