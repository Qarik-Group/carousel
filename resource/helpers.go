package resource

import (
	"sync"

	credhubcli "code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/auth"
	oc "github.com/cloudboss/ofcourse/ofcourse"
	cbosh "github.com/starkandwayne/carousel/bosh"
	"github.com/starkandwayne/carousel/config"
	ccredhub "github.com/starkandwayne/carousel/credhub"
	. "github.com/starkandwayne/carousel/state"
)

var (
	credhub  ccredhub.CredHub
	director cbosh.Director
	state    State
)

func initializeFromSource(source oc.Source, logger *oc.Logger) {

	cfg := config.Config{
		Bosh: &config.Bosh{
			Environment:  source["bosh_environment"].(string),
			Client:       source["bosh_client"].(string),
			ClientSecret: source["bosh_client_secret"].(string),
			CaCert:       source["bosh_ca_cert"].(string),
		},
		Credhub: &config.Credhub{
			Server: source["credhub_server"].(string),
			Client: source["credhub_client"].(string),
			Secret: source["credhub_secret"].(string),
			CaCert: source["credhub_ca_cert"].(string),
		},
	}

	chcli, err := credhubcli.New(
		cfg.Credhub.Server,
		credhubcli.SkipTLSValidation(true), // TODO use CA
		credhubcli.Auth(auth.UaaClientCredentials(cfg.Credhub.Client, cfg.Credhub.Secret)),
	)
	if err != nil {
		logger.Errorf("failed to connect to Credhub: %s", err)
	}

	credhub = ccredhub.NewCredHub(chcli)

	director, err = cbosh.NewDirector(cfg.Bosh)
	if err != nil {
		logger.Errorf("failed to connect to BOSH Director: %s", err)
	}

	state = NewState()
}

func refresh(logger *oc.Logger) error {
	var (
		wg          sync.WaitGroup
		err         error
		credentials []*ccredhub.Credential
		variables   []*cbosh.Variable
	)

	wg.Add(1)
	go func(wg *sync.WaitGroup, credentials *[]*ccredhub.Credential) {
		defer wg.Done()
		*credentials, err = credhub.FindAll()
		if err != nil {
			logger.Errorf("failed to load credentials from Credhub: %s", err)
		}

	}(&wg, &credentials)

	wg.Add(1)
	go func(wg *sync.WaitGroup, variables *[]*cbosh.Variable) {
		defer wg.Done()
		*variables, err = director.GetVariables()
		if err != nil {
			logger.Errorf("failed to load variables from BOSH Director: %s", err)
		}

	}(&wg, &variables)

	wg.Wait()

	return state.Update(credentials, variables)
}
