package store

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"code.cloudfoundry.org/credhub-cli/credhub"
	"github.com/emirpasic/gods/maps/treebidimap"
	"github.com/emirpasic/gods/utils"
)

func NewStore(ch *credhub.CredHub) (*Store, error) {
	certs, err := ch.GetAllCertificatesMetadata()
	if err != nil {
		return nil, err
	}

	certsStore := treebidimap.NewWith(utils.StringComparator, certByName)
	certVersionsStore := treebidimap.NewWith(utils.StringComparator, certVersionById)
	for _, certMeta := range certs {
		cert := Cert{
			Id:   certMeta.Id,
			Name: certMeta.Name,
		}

		versions := make([]*CertVersion, 0)
		for _, certMetaVersion := range certMeta.Versions {
			cv := CertVersion{
				Id:   certMetaVersion.Id,
				Cert: &cert,
			}
			versions = append(versions, &cv)
			certVersionsStore.Put(cv.Id, &cv)
		}
		cert.Versions = versions
		certsStore.Put(certMeta.Name, &cert)
	}

	for _, certMeta := range certs {
		credentials, err := ch.GetAllVersions(certMeta.Name)
		if err != nil {
			return nil, err
		}

		for _, c := range credentials {
			if c.Base.Type == "certificate" {
				raw := c.Value.(map[string]interface{})
				rawCa := raw["ca"].(string)
				rawCert := raw["certificate"].(string)

				certBlock, _ := pem.Decode([]byte(rawCert))
				certificate, err := x509.ParseCertificate(certBlock.Bytes)
				if err != nil {
					return nil, fmt.Errorf("failed to parse certificate: %s", err)
				}

				caBlock, _ := pem.Decode([]byte(rawCa))
				ca, err := x509.ParseCertificate(caBlock.Bytes)
				if err != nil {
					return nil, fmt.Errorf("failed to parse ca: %s", err)
				}

				cv, _ := certVersionsStore.Get(c.Base.Id)
				certVersion := cv.(*CertVersion)
				certVersion.Certificate = certificate
				certVersion.Ca = ca
			}
		}

		//	cert := certsStore.Get(certMeta.Name)
	}

	return &Store{
		Certs:         certsStore,
		CertVersions:  certVersionsStore,
		credhubClient: ch,
	}, nil
}
