// Package resource is an implementation of a Concourse resource.
package resource

import (
	"crypto/sha256"
	oc "github.com/cloudboss/ofcourse/ofcourse"

	cstate "github.com/starkandwayne/carousel/state"
)

type Resource struct{}

// Check implements the ofcourse.Resource Check method, corresponding to the /opt/resource/check command.
// This is called when Concourse does its resource checks, or when the `fly check-resource` command is run.
func (r *Resource) Check(source oc.Source, version oc.Version, env oc.Environment,
	logger *oc.Logger) ([]oc.Version, error) {
	// Returned `versions` should be all of the versions since the one given in the `version`
	// argument. If `version` is nil, then return the first available version. In many cases there
	// will be only one version to return, depending on the type of resource being implemented.
	// For example, a git resource would return a list of commits since the one given in the
	// `version` argument, whereas that would not make sense for resources which do not have any
	// kind of linear versioning.
	initializeFromSource(source)

	deployment := source["deployment"].(string)
	// filters.latest = true
	if deployment == "" {
		logger.Errorf("deployment flag must be set")
	}

	logger.Infof("Refreshing state")
	refresh()
	logger.Infof(" done\n\n")

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
		hash := sha256.Sum256([]byte(allVersions))
		versions = append(versions, oc.Version{"hash": string(hash[:])})
	}

	return versions, nil
}

// In implements the ofcourse.Resource In method, corresponding to the /opt/resource/in command.
// This is called when a Concourse job does `get` on the resource.
func (r *Resource) In(outputDirectory string, source oc.Source, params oc.Params, version oc.Version,
	env oc.Environment, logger *oc.Logger) (oc.Version, oc.Metadata, error) {
	logger.Errorf("Unimplemented")
	return nil, nil, nil
}

// Out implements the ofcourse.Resource Out method, corresponding to the /opt/resource/out command.
// This is called when a Concourse job does a `put` on the resource.
func (r *Resource) Out(inputDirectory string, source oc.Source, params oc.Params,
	env oc.Environment, logger *oc.Logger) (oc.Version, oc.Metadata, error) {
	logger.Errorf("Unimplemented")
	return nil, nil, nil
}
