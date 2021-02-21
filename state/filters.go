package state

import (
	"github.com/starkandwayne/carousel/credhub"
)

type filter func(*Credential) bool

func SelfSignedFilter() filter {
	return func(c *Credential) bool {
		return c.SignedBy == nil
	}
}

func LatestFilter() filter {
	return func(c *Credential) bool {
		return c.Latest
	}
}

func TypeFilter(types ...credhub.CredentialType) filter {
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
