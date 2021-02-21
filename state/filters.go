package state

import (
	"github.com/starkandwayne/carousel/credhub"
)

type Filter func(*Credential) bool

func OrFilter(fns ...Filter) Filter {
	return func(c *Credential) bool {
		for _, fn := range fns {
			if fn(c) {
				return true
			}
		}
		return false
	}
}

func AndFilter(fns ...Filter) Filter {
	return func(c *Credential) bool {
		for _, fn := range fns {
			if !fn(c) {
				return false
			}
		}
		return true
	}
}

func SelfSignedFilter() Filter {
	return func(c *Credential) bool {
		return c.SignedBy == nil
	}
}

func ActiveFilter() Filter {
	return func(c *Credential) bool {
		return len(c.Deployments) != 0
	}
}

func LatestFilter() Filter {
	return func(c *Credential) bool {
		return c.Latest
	}
}

func SigningFilter() Filter {
	return func(c *Credential) bool {
		if c.Signing != nil && *c.Signing {
			return true
		}
		return false
	}
}

func TransitionalFilter() Filter {
	return func(c *Credential) bool {
		return c.Transitional
	}
}

func TypeFilter(types ...credhub.CredentialType) Filter {
	return func(c *Credential) bool {
		for _, t := range types {
			if c.Type == t {
				return true
			}
		}
		return false
	}
}

func DeploymentFilter(deployments ...string) Filter {
	return func(c *Credential) bool {
		for _, name := range deployments {
			for _, d := range c.Deployments {
				if d.Name == name {
					return true
				}
			}
		}
		return false
	}
}
