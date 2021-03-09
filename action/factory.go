package action

import (
	"time"

	"github.com/starkandwayne/carousel/state"
)

type ActionFactory interface {
	NextAction(*state.Credential) []Action
}

type ConcreteActionFactory struct {
	OlderThan     time.Time
	ExpiresBefore time.Time
}

func (f *ConcreteActionFactory) NextAction(cred *state.Credential) (actions []Action) {
	actions = make([]Action, 0)

	if cred.Latest && cred.VersionCreatedAt.Before(f.OlderThan) {
		actions = append(actions, &regenerateAction{subject: cred})
		return
	}

	if cred.Latest {
		pendingDeploys := cred.PendingDeploys()
		if len(pendingDeploys) != 0 {
			for _, d := range pendingDeploys {
				actions = append(actions, &boshDeployAction{
					subject:    cred,
					deployment: d.Name,
				})
			}
			return
		}
	}

	if !cred.Active() {
		actions = append(actions, &cleanUpAction{subject: cred})
		return
	}

	if cred.Signing != nil && *cred.Signing {
		latest, found := cred.Path.Versions.Find(state.LatestFilter())
		if found && latest.Transitional && latest.Active() && len(latest.PendingDeploys()) == 0 {
			actions = append(actions, &markTransitionalAction{subject: cred})
			return
		}
	}

	if cred.Latest && cred.ExpiryDate != nil &&
		cred.ExpiryDate.Before(f.ExpiresBefore) &&
		!cred.SignedBy.ExpiryDate.Before(f.ExpiresBefore) {

		actions = append(actions, &regenerateAction{subject: cred})
		return
	}

	actions = append(actions, &noOpAction{subject: cred})
	return
}
