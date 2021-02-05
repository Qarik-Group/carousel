package store

import (
	"code.cloudfoundry.org/credhub-cli/credhub/credentials"
	"github.com/emirpasic/gods/maps/treebidimap"
	"github.com/emirpasic/gods/utils"
)

func NewStore(certs []credentials.CertificateMetadata) *Store {
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

	// for _, certMeta := range certs; {
	//	cert := certsStore.Get(certMeta.Name)
	// }

	return &Store{
		Certs:        certsStore,
		CertVersions: certVersionsStore,
	}
}
