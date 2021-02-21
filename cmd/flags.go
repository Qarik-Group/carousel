package cmd

import (
	"github.com/spf13/pflag"
	ccredhub "github.com/starkandwayne/carousel/credhub"
	. "github.com/starkandwayne/carousel/state"
)

type credentialFilters struct {
	deployments   []string
	types         []string
	expiresWithin string
	all           bool
	signing       bool
}

var filters = credentialFilters{
	deployments:   make([]string, 0),
	expiresWithin: "",
	types:         make([]string, 0),
	all:           false,
	signing:       false,
}

func (f credentialFilters) Filters() []Filter {
	out := make([]Filter, 0)
	if len(f.deployments) != 0 {
		out = append(out, DeploymentFilter(f.deployments...))
	}
	if len(f.types) != 0 {
		types := make([]ccredhub.CredentialType, 0)
		for _, t := range f.types {
			ct, err := ccredhub.CredentialTypeString(t)
			if err != nil {
				logger.Fatalf("Invalid credential type: %s got: %s", t, err)
			}
			types = append(types, ct)
		}
		out = append(out, TypeFilter(types...))
	}
	if !f.all {
		out = append(out, OrFilter(
			AndFilter(
				LatestFilter(),
				ActiveFilter(),
			),
			SigningFilter(),
			TransitionalFilter(),
		))
	}
	if f.signing {
		out = append(out, SigningFilter())
	}
	return out
}

func addFilterFlags(set *pflag.FlagSet) {
	set.StringSliceVarP(&filters.types, "types", "t", ccredhub.CredentialTypeStringValues(),
		"filter by credential type (comma sperated)")
	set.StringSliceVarP(&filters.deployments, "deployments", "d", nil,
		"filter by deployment names (comma sperated)")
	set.BoolVar(&filters.all, "include-all", false,
		"also show unused credential versions")
	set.BoolVar(&filters.signing, "signing", false,
		"only show Certificates used to sign")
}
