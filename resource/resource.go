// Package resource is an implementation of a Concourse resource.
package resource

import (
	"crypto/sha256"
	"encoding/hex"

	oc "github.com/cloudboss/ofcourse/ofcourse"

	cstate "github.com/starkandwayne/carousel/state"
)

type Resource struct{}

// Check implements the ofcourse.Resource Check method, corresponding to the /opt/resource/check command.
// This is called when Concourse does its resource checks, or when the `fly check-resource` command is run.
func (r *Resource) Check(source oc.Source, version oc.Version, env oc.Environment,
	logger *oc.Logger) ([]oc.Version, error) {
	initializeFromSource(source, logger)

	deployment := source["deployment"].(string)
	if deployment == "" {
		logger.Errorf("deployment flag must be set")
	}

	logger.Infof("Refreshing state for deployment '%s'", deployment)
	refresh(logger)
	logger.Infof("done\n")

	credentials := state.Credentials(cstate.DeploymentFilter(deployment))
	credentials.SortByNameAndCreatedAt()

	deployNeeded := false
	allVersions := ""
	for _, cred := range credentials {
		if cred.PendingDeploys().IncludesName(deployment) {
			deployNeeded = true
			allVersions += cred.ID
		}
	}
	versions := []oc.Version{}
	if deployNeeded {
		logger.Infof("Found credentials to re-deploy")
		hash := sha256.Sum256([]byte(allVersions))
		versions = append(versions, oc.Version{"hash": hex.EncodeToString(hash[:])})
	} else {
		logger.Infof("Nothing to deploy")
	}

	return versions, nil
}

// In implements the ofcourse.Resource In method, corresponding to the /opt/resource/in command.
// This is called when a Concourse job does `get` on the resource.
func (r *Resource) In(outputDirectory string, source oc.Source, params oc.Params, version oc.Version,
	env oc.Environment, logger *oc.Logger) (oc.Version, oc.Metadata, error) {
	metadata := oc.Metadata{}

	return version, metadata, nil
}

// Out implements the ofcourse.Resource Out method, corresponding to the /opt/resource/out command.
// This is called when a Concourse job does a `put` on the resource.
func (r *Resource) Out(inputDirectory string, source oc.Source, params oc.Params,
	env oc.Environment, logger *oc.Logger) (oc.Version, oc.Metadata, error) {
	logger.Errorf("Unimplemented")
	return nil, nil, nil
}
