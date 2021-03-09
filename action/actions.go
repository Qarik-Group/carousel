package action

import (
	"fmt"

	"github.com/starkandwayne/carousel/state"
)

type noOpAction struct {
	subject *state.Credential
}

func (a *noOpAction) Name() string {
	return fmt.Sprintf("%s - up-to-date", a.subject.PathVersion())
}

func (a *noOpAction) Description() string {
	return fmt.Sprintf("no action required on credential %s",
		a.subject.PathVersion())
}

type boshDeployAction struct {
	subject    *state.Credential
	deployment string
	reason     string
}

func (a *boshDeployAction) Name() string {
	return fmt.Sprintf("%s - deploy(%s)", a.subject.PathVersion(), a.deployment)
}

func (a *boshDeployAction) Description() string {
	return fmt.Sprintf("waiting for credential version %s to be used by deployment %s",
		a.subject.PathVersion(), a.deployment)
}

type regenerateAction struct {
	subject *state.Credential
}

func (a *regenerateAction) Name() string {
	return fmt.Sprintf("%s - regenerate", a.subject.PathVersion())
}

func (a *regenerateAction) Description() string {
	return fmt.Sprintf("regenerate %s", a.subject.PathVersion())
}

type cleanUpAction struct {
	subject *state.Credential
}

func (a *cleanUpAction) Name() string {
	return fmt.Sprintf("%s - clean-up", a.subject.PathVersion())
}

func (a *cleanUpAction) Description() string {
	return fmt.Sprintf("clean-up credential version %s since it is no longer used", a.subject.PathVersion())
}
