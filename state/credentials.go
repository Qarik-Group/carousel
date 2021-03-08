package state

import (
	"sort"
)

func (s *state) Credentials(filters ...Filter) Credentials {
	creds := s.credentials.Select(func(_, v interface{}) bool {
		for _, fn := range filters {
			if !fn(v.(*Credential)) {
				return false
			}
		}
		return true
	})
	out := make(Credentials, 0, creds.Size())
	for _, cert := range creds.Values() {
		out = append(out, cert.(*Credential))
	}
	return out
}

func (s *state) Paths() []*Path {
	paths := s.paths.Select(func(_, v interface{}) bool {
		return true
	})
	out := make([]*Path, 0, paths.Size())
	for _, path := range paths.Values() {
		out = append(out, path.(*Path))
	}
	return out
}

func (creds Credentials) Collect(fn Collector) Credentials {
	out := make(Credentials, 0)
	for _, cred := range creds {
		if c := fn(cred); c != nil {
			out = append(out, c...)
		}
	}
	return out
}

func (creds Credentials) Select(filters ...Filter) Credentials {
	out := make(Credentials, 0)

OUTER:
	for _, cred := range creds {
		for _, fn := range filters {
			if !fn(cred) {
				continue OUTER
			}
		}
		out = append(out, cred)
	}
	return out
}

func (creds Credentials) Find(filters ...Filter) (cred *Credential, found bool) {
OUTER:
	for _, cred := range creds {
		for _, fn := range filters {
			if !fn(cred) {
				continue OUTER
			}
		}
		return cred, true
	}
	return nil, false
}

func (creds Credentials) Unique() Credentials {
	out := make(Credentials, 0)
	for _, cred := range creds {
		if !out.Includes(cred) {
			out = append(out, cred)
		}
	}
	return out
}

func (ceds Credentials) Includes(this *Credential) bool {
	for _, cred := range ceds {
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

func (creds Credentials) Any() bool {
	return len(creds) != 0
}
