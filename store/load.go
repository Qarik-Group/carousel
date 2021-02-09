package store

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"code.cloudfoundry.org/credhub-cli/credhub"
	"github.com/emirpasic/gods/maps/treebidimap"
	"github.com/emirpasic/gods/utils"
)

func NewStore(ch *credhub.CredHub) (*Store, error) {
	certs, err := ch.GetAllCertificatesMetadata()
	if err != nil {
		return nil, err
	}

	store := Store{
		certs:         treebidimap.NewWith(utils.StringComparator, certByName),
		certVersions:  treebidimap.NewWith(utils.StringComparator, certVersionById),
		credhubClient: ch,
	}

	for _, certMeta := range certs {
		cert := Cert{
			Id:   certMeta.Id,
			Name: certMeta.Name,
		}

		versions := make([]*CertVersion, 0)
		for _, certMetaVersion := range certMeta.Versions {
			expiry, err := time.Parse(time.RFC3339, certMetaVersion.ExpiryDate)
			if err != nil {
				return nil, fmt.Errorf("failed to parse expiry date: %s for cert version: %s",
					certMetaVersion.ExpiryDate, certMetaVersion.Id)
			}
			cv := CertVersion{
				Id:                   certMetaVersion.Id,
				Cert:                 &cert,
				Transitional:         certMetaVersion.Transitional,
				CertificateAuthority: certMetaVersion.CertificateAuthority,
				SelfSigned:           certMetaVersion.SelfSigned,
				Expiry:               expiry,
			}
			versions = append(versions, &cv)
			store.certVersions.Put(cv.Id, &cv)
		}
		cert.Versions = versions
		store.certs.Put(certMeta.Name, &cert)
	}

	// for each certMeta fetch raw cert + ca and decode with x509
	for _, certMeta := range certs {
		credentials, err := ch.GetAllVersions(certMeta.Name)
		if err != nil {
			return nil, err
		}

		for _, c := range credentials {
			if c.Base.Type == "certificate" {
				raw := c.Value.(map[string]interface{})
				// rawCa := raw["ca"].(string)
				rawCert := raw["certificate"].(string)

				certBlock, _ := pem.Decode([]byte(rawCert))
				certificate, err := x509.ParseCertificate(certBlock.Bytes)
				if err != nil {
					return nil, fmt.Errorf("failed to parse certificate: %s", err)
				}

				// caBlock, _ := pem.Decode([]byte(rawCa))
				// ca, err := x509.ParseCertificate(caBlock.Bytes)
				// if err != nil {
				//	return nil, fmt.Errorf("failed to parse ca: %s", err)
				// }

				cv, _ := store.certVersions.Get(c.Base.Id)
				certVersion := cv.(*CertVersion)
				certVersion.Certificate = certificate
				//				certVersion.Ca = ca
			}
		}
	}

	// Lookup Ca for each cert
	it := store.certVersions.Iterator()
	for it.End(); it.Prev(); {
		_, value := it.Key(), it.Value()
		v := value.(*CertVersion)
		authorityKeyID := v.Certificate.AuthorityKeyId
		if v.SelfSigned {
			continue
		}
		ca, found := store.GetCertVersionBySubjectKeyId(authorityKeyID)
		if found {
			ca.Signs = append(ca.Signs, v)
			v.SignedBy = ca
		} else {
			return nil, fmt.Errorf("failed to lookup ca CertVersion with id: %s", v.Id)
		}
	}

	return &store, nil
}
