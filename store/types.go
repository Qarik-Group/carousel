package store

import (
	"bytes"
	"crypto/x509"
	"fmt"
	"net/http"
	"time"

	//	"code.cloudfoundry.org/credhub-cli/credhub/credentials"
	"code.cloudfoundry.org/credhub-cli/credhub"
	"github.com/emirpasic/gods/maps/treebidimap"

	boshdir "github.com/cloudfoundry/bosh-cli/director"
)

type Store struct {
	deployments         *treebidimap.Map
	certVersions        *treebidimap.Map
	certs               *treebidimap.Map
	variableDefinitions *treebidimap.Map
	credhubClient       *credhub.CredHub
	directorClient      boshdir.Director
}

type Cert struct {
	Id                 string
	Name               string
	Versions           []*CertVersion
	VariableDefinition *VariableDefinition
}

type CertVersion struct {
	Id                   string
	Expiry               time.Time
	Transitional         bool
	CertificateAuthority bool
	SelfSigned           bool
	Cert                 *Cert
	Deployments          []*Deployment
	SignedBy             *CertVersion
	Signs                []*CertVersion
	Certificate          *x509.Certificate
}

func (s *Store) ToggleTransitional(cv *CertVersion) error {
	path := fmt.Sprintf("/api/v1/certificates/%s/update_transitional_version", cv.Cert.Id)
	body := map[string]interface{}{"version": cv.Id}
	if cv.Transitional == true {
		body["version"] = nil
	}
	_, err := s.credhubClient.Request(http.MethodPut, path, nil, body, true)
	if err != nil {
		return fmt.Errorf("failed request: %s with body: %s got: %s", path, body, err)
	}
	return nil
}

type Deployment struct {
	Versions []*CertVersion
	Name     string
}

type VariableDefinition struct {
	Name    string      `yaml:"name"`
	Type    string      `yaml:"type"`
	Options interface{} `yaml:"options,omitempty"`
}

func (s *Store) EachCert(fn func(*Cert)) {
	s.certs.Each(func(_, v interface{}) {
		fn(v.(*Cert))
	})
}

func (cv *CertVersion) Status() string {
	status := "active"
	if cv.Expiry.Sub(time.Now()) < time.Hour*24*30 {
		status = "notice"
	}
	if len(cv.Deployments) == 0 {
		status = "unused"
	}
	return status
}

func certByName(a, b interface{}) int {
	// Type assertion, program will panic if this is not respected
	c1 := a.(*Cert)
	c2 := b.(*Cert)

	switch {
	case c1.Name > c2.Name:
		return 1
	case c1.Name < c2.Name:
		return -1
	default:
		return 0
	}
}

func certVersionById(a, b interface{}) int {
	// Type assertion, program will panic if this is not respected
	c1 := a.(*CertVersion)
	c2 := b.(*CertVersion)

	switch {
	case c1.Id > c2.Id:
		return 1
	case c1.Id < c2.Id:
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

func (s *Store) getCertVersionBySubjectKeyId(keyId []byte) (*CertVersion, bool) {
	_, foundValue := s.certVersions.Find(func(index interface{}, value interface{}) bool {
		return bytes.Compare(value.(*CertVersion).Certificate.SubjectKeyId, keyId) == 0

	})
	if foundValue != nil {
		return foundValue.(*CertVersion), true
	}
	return nil, false
}
