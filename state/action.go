package state

import (
	"time"

	"github.com/starkandwayne/carousel/bosh"
	"github.com/starkandwayne/carousel/credhub"
)

//
//go:generate go run github.com/alvaroloes/enumer -type=Action -json -transform=snake

type Action int

const (
	None Action = iota
	NoOverwrite
	BoshDeploy
	Regenerate
	CleanUp
	MarkTransitional
	UnMarkTransitional
)

type RegenerationCriteria struct {
	OlderThan        time.Time
	ExpiresBefore    time.Time
	IgnoreUpdateMode bool
}

func (cred *Credential) NextAction(r RegenerationCriteria) Action {
	for _, ct := range []credhub.CredentialType{credhub.JSON, credhub.Value} {
		if cred.Type == ct {
			return None
		}
	}

	if !r.IgnoreUpdateMode && cred.Path.VariableDefinition != nil &&
		cred.Path.VariableDefinition.UpdateMode == bosh.NoOverwrite {
		return NoOverwrite
	}

	if cred.Signing != nil && *cred.Signing {
		latest, found := cred.Path.Versions.Find(LatestFilter())
		if found && latest.Transitional && latest.Active() && len(latest.PendingDeploys()) == 0 {
			return MarkTransitional
		}
	}

	if cred.Latest && cred.ExpiryDate != nil &&
		cred.ExpiryDate.Before(r.ExpiresBefore) {
		if cred.SignedBy == nil {
			return Regenerate
		} else if !cred.SignedBy.ExpiryDate.Before(r.ExpiresBefore) {
			return Regenerate
		} else {
			return None
		}
	}

	if cred.Latest && cred.VersionCreatedAt.Before(r.OlderThan) {
		return Regenerate
	}

	if cred.Latest && len(cred.PendingDeploys()) != 0 {
		return BoshDeploy
	}

	if cred.Transitional {
		signing, found := cred.Path.Versions.Find(SigningFilter())
		if found && len(signing.PendingDeploys()) == 0 {
			return UnMarkTransitional
		}
	}

	if !cred.Active() {
		return CleanUp
	}

	return None
}
