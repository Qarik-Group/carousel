package store

import (
	"crypto/x509"
	"time"

	//	"code.cloudfoundry.org/credhub-cli/credhub/credentials"
	"code.cloudfoundry.org/credhub-cli/credhub"
	"github.com/emirpasic/gods/maps/treebidimap"
)

type Store struct {
	Deployments   *treebidimap.Map
	CertVersions  *treebidimap.Map
	Certs         *treebidimap.Map
	credhubClient *credhub.CredHub
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
	Signs                []*CertVersion
	Certificate          *x509.Certificate
	Ca                   *x509.Certificate
}

type Deployment struct {
	Versions []*CertVersion
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
