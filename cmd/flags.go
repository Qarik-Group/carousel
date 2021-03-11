package cmd

import (
	"fmt"
	"time"

	"github.com/karrick/tparse"
	"github.com/spf13/pflag"
	ccredhub "github.com/starkandwayne/carousel/credhub"
	. "github.com/starkandwayne/carousel/state"
	cstate "github.com/starkandwayne/carousel/state"
)

type actionCriteria struct {
	expiresWithin    string
	olderThan        string
	ignoreUpdateMode bool
}

var criteria = actionCriteria{}

func (c actionCriteria) RegenerationCriteria() (cstate.RegenerationCriteria, error) {
	ew, err := tparse.AddDuration(time.Now(), "+"+c.expiresWithin)
	if err != nil {
		return cstate.RegenerationCriteria{},
			fmt.Errorf("failed to parse --expires-within flag into duration: %s, got: %s",
				c.expiresWithin, err)
	}

	ot, err := tparse.AddDuration(time.Now(), "-"+c.olderThan)
	if err != nil {
		return cstate.RegenerationCriteria{},
			fmt.Errorf("failed to parse --older-than flag into duration: %s, got: %s",
				c.olderThan, err)
	}

	return cstate.RegenerationCriteria{
		OlderThan:        ot,
		ExpiresBefore:    ew,
		IgnoreUpdateMode: c.ignoreUpdateMode,
	}, nil
}

type credentialFilters struct {
	deployment    string
	name          string
	types         []string
	unused        bool
	expiresWithin string
	olderThan     string
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
	if f.unused {
		out = append(out, AndFilter(
			NotFilter(ActiveFilter()),
			NotFilter(TransitionalFilter()),
			NotFilter(AnyFilter(SignsCollector())),
		))
	}
	if f.expiresWithin != "" {
		ew, err := tparse.AddDuration(time.Now(), "+"+f.expiresWithin)
		if err != nil {
			logger.Fatalf("failed to parse duration: %s, got: %s",
				f.expiresWithin, err)
		}
		out = append(out, ExpiresBeforeFilter(ew))

	}
	if f.olderThan != "" {
		ot, err := tparse.AddDuration(time.Now(), "-"+f.olderThan)
		if err != nil {
			logger.Fatalf("failed to parse duration: %s, got: %s",
				f.expiresWithin, err)
		}
		out = append(out, OlderThanFilter(ot))

	}
	if f.latest {
		out = append(out, LatestFilter())
	}
	if f.signing {
		out = append(out, SigningFilter())
	}
	if f.signedBy != "" {
		out = append(out, SignedByFilter(f.signedBy))
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

func addTypesFlag(set *pflag.FlagSet) {
	set.StringSliceVarP(&filters.types, "types", "t", ccredhub.CredentialTypeStringValues(),
		"filter by credential type (comma sperated)")
}

func addDeploymentFlag(set *pflag.FlagSet) {
	set.StringVarP(&filters.deployment, "deployment", "d", "",
		"filter by deployment name")

}

func addNameFlag(set *pflag.FlagSet) {
	set.StringVar(&filters.name, "name", "",
		"only credential with name")

}

func addSignedByFlag(set *pflag.FlagSet) {
	set.StringVar(&filters.signedBy, "signed-by", "",
		"filter certificates signed by a specific CA")
}

func addExpiresWithinFlag(set *pflag.FlagSet) {
	set.StringVar(&filters.expiresWithin, "expires-within", "",
		"filter certificates by expiry window (suffixes: d day, w week, y year)")
}

func addOlderThanFlag(set *pflag.FlagSet) {
	set.StringVar(&filters.olderThan, "older-than", "",
		"filter credentials by age (suffixes: d day, w week, y year)")
}

func addExpiresWithinCriteriaFlag(set *pflag.FlagSet) {
	set.StringVar(&criteria.expiresWithin, "expires-within", "3m",
		"regenerate certificates by expiry window (suffixes: d day, w week, y year)")
}

func addOlderThanCireteriaFlag(set *pflag.FlagSet) {
	set.StringVar(&criteria.olderThan, "older-than", "1y",
		"regenerate credentials by age (suffixes: d day, w week, y year)")
}

func addIgnoreUpdateModeCireteriaFlag(set *pflag.FlagSet) {
	set.BoolVar(&criteria.ignoreUpdateMode, "ignore-update-mode", false,
		"ignore the value of BOSH /variables/.../update_mode")
}
