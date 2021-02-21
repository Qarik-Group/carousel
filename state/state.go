package state

import (
	"github.com/emirpasic/gods/maps/treebidimap"
	"github.com/emirpasic/gods/utils"
	"github.com/starkandwayne/carousel/bosh"
	"github.com/starkandwayne/carousel/credhub"
)

type State interface {
	Update([]*credhub.Credential, []*bosh.Variable) error
	Credentials(...filter) []*Credential
}

func NewState() State {
	return &state{
		deployments: treebidimap.NewWith(utils.StringComparator, deploymentComparator),
		paths:       treebidimap.NewWith(utils.StringComparator, pathComparator),
		credentials: treebidimap.NewWith(utils.StringComparator, credentialComparator),
	}
}

type state struct {
	deployments *treebidimap.Map
	paths       *treebidimap.Map
	credentials *treebidimap.Map
}

func (s *state) clear() {
	s.paths.Clear()
	s.credentials.Clear()
	s.deployments.Clear()
}
