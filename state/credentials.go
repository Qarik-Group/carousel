package state

import "github.com/starkandwayne/carousel/credhub"

func TypeFilter(types ...credhub.CredentialType) Filter {
	return func(c *Credential) bool {
		match := false
		for _, t := range types {
			if c.Type == t {
				match = true
			}
		}
		return match
	}
}

func (s *state) Credentials(filters ...Filter) []*Credential {
	certs := s.credentials.Select(func(_, v interface{}) bool {
		for _, fn := range filters {
			if !fn(v.(*Credential)) {
				return false
			}
		}
		return true
	})
	out := make([]*Credential, 0, certs.Size())
	for _, cert := range certs.Values() {
		out = append(out, cert.(*Credential))
	}
	return out
}
