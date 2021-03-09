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

	if cred.Latest && (len(cred.Deployments) != len(cred.Path.Deployments)) {
		for _, d := range cred.Path.Deployments {
			if !cred.Deployments.Includes(d) {
				actions = append(actions, &boshDeployAction{
					subject:    cred,
					deployment: d.Name,
				})
			}
		}
		return
	}

	if !cred.Active() {
		actions = append(actions, &cleanUpAction{subject: cred})
		return
	}

	actions = append(actions, &noOpAction{subject: cred})
	return
}
