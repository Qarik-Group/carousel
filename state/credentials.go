package state

import (
	"sort"
)

func (s *state) Credentials(filters ...Filter) Credentials {
	certs := s.credentials.Select(func(_, v interface{}) bool {
		for _, fn := range filters {
			if !fn(v.(*Credential)) {
				return false
			}
		}
		return true
	})
	out := make(Credentials, 0, certs.Size())
	for _, cert := range certs.Values() {
		out = append(out, cert.(*Credential))
	}
	return out
}

func (c Credentials) LatestVersion() *Credential {
	for _, cred := range c {
		if cred.Latest {
			return cred
		}
	}
	return nil
}

func (c Credentials) ActiveVersions() Credentials {
	out := make(Credentials, 0)
	for _, cred := range c {
		if len(cred.Deployments) != 0 {
			out = append(out, cred)
		}
	}
	return out
}

func (c Credentials) SigningVersion() *Credential {
	for _, cred := range c {
		if cred.Signing != nil && *cred.Signing {
			return cred
		}
	}
	return nil
}

func (c Credentials) Includes(this *Credential) bool {
	for _, cred := range c {
		if cred == this {
			return true
		}
	}
	return false
}

func (creds Credentials) SortByCreatedAt() {
	sort.Slice(creds, func(i, j int) bool {
		return creds[i].VersionCreatedAt.After(*creds[j].VersionCreatedAt)
	})
}

func (creds Credentials) SortByNameAndCreatedAt() {
	sort.Slice(creds, func(i, j int) bool {
		if creds[i].Name == creds[j].Name {
			return creds[i].VersionCreatedAt.After(*creds[j].VersionCreatedAt)
		}
		return creds[i].Name < creds[j].Name
	})
}
