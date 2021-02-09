package store

import (
	"bytes"
	"crypto/x509"
	"time"

	//	"code.cloudfoundry.org/credhub-cli/credhub/credentials"
	"code.cloudfoundry.org/credhub-cli/credhub"
	"github.com/emirpasic/gods/maps/treebidimap"

	boshdir "github.com/cloudfoundry/bosh-cli/director"
)

type Store struct {
	deployments    *treebidimap.Map
	certVersions   *treebidimap.Map
	certs          *treebidimap.Map
	credhubClient  *credhub.CredHub
	directorClient boshdir.Director
}

type Cert struct {
	Id       string
	Name     string
	Versions []*CertVersion
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
	//	Ca                   *x509.Certificate
}

type Deployment struct {
	Versions []*CertVersion
	Name     string
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

func (s *Store) GetCertVersionBySubjectKeyId(keyId []byte) (*CertVersion, bool) {
	_, foundValue := s.certVersions.Find(func(index interface{}, value interface{}) bool {
		return bytes.Compare(value.(*CertVersion).Certificate.SubjectKeyId, keyId) == 0

	})
	if foundValue != nil {
		return foundValue.(*CertVersion), true
	}
	return nil, false
}
