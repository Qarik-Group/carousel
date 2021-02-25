package cmd

import (
	"time"

	"github.com/karrick/tparse"
	"github.com/spf13/pflag"
	ccredhub "github.com/starkandwayne/carousel/credhub"
	. "github.com/starkandwayne/carousel/state"
)

type credentialFilters struct {
	deployment    string
	name          string
	types         []string
	expiresWithin string
	latest        bool
	signing       bool
	signedBy      string
	ca            bool
	leaf          bool
}

var filters = credentialFilters{}

func (f credentialFilters) Filters() []Filter {
	out := make([]Filter, 0)
	if f.deployment != "" {
		out = append(out, DeploymentFilter(f.deployment))
	}
	if f.name != "" {
		out = append(out, NameFilter(f.name))
	}
	if f.expiresWithin != "" {
		ew, err := tparse.AddDuration(time.Now(), "+"+f.expiresWithin)
		if err != nil {
			logger.Fatalf("failed to parse duration: %s, got: %s",
				f.expiresWithin, err)
		}
		out = append(out, ExpiresBeforeFilter(ew))

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
	if f.latest {
		out = append(out, LatestFilter())
	}
	if f.signing {
		out = append(out, SigningFilter())
	}
	if f.ca {
		out = append(out, CertificateAuthorityFilter(true))
	}
	if f.leaf {
		out = append(out,
			TypeFilter(ccredhub.Certificate),
			CertificateAuthorityFilter(false),
		)
	}
	return out
}

func addFilterFlags(set *pflag.FlagSet) {
	set.StringSliceVarP(&filters.types, "types", "t", ccredhub.CredentialTypeStringValues(),
		"filter by credential type (comma sperated)")
	set.BoolVar(&filters.signing, "signing", false,
		"only show Certificates used to sign")
}

func addDeploymentFlag(set *pflag.FlagSet) {
	set.StringVarP(&filters.deployment, "deployment", "d", "",
		"filter by deployment name")

}

func addNameFlag(set *pflag.FlagSet) {
	set.StringVar(&filters.name, "name", "",
		"only credential with name")

}

func addExpiresWithinFlag(set *pflag.FlagSet) {
	set.StringVar(&filters.expiresWithin, "expires-within", "",
		"filter certificates by expiry window. Valid units are d for days, w for weeks, m for months, and y for years.")
}

func addSignedByFlag(set *pflag.FlagSet) {
	set.StringVar(&filters.signedBy, "sign-by", "",
		"filter certificates signed by a specific CA.")
}
