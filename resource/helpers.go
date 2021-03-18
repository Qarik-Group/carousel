package resource

import (
	"log"
	"os"
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
	logger   *log.Logger
	credhub  ccredhub.CredHub
	director cbosh.Director
	state    State
)

func initializeFromSource(source oc.Source) {
	logger = log.New(os.Stderr, "", 0)

	bosh_cfg := source["bosh"].(map[string]string)
	credhub_cfg := source["credhub"].(map[string]string)
	cfg := config.Config{
		Bosh: &config.Bosh{
			Environment:  bosh_cfg["environment"],
			Client:       bosh_cfg["client"],
			ClientSecret: bosh_cfg["client_secret"],
			CaCert:       bosh_cfg["ca_cert"],
		},
		Credhub: &config.Credhub{
			Server: credhub_cfg["server"],
			Client: credhub_cfg["client"],
			Secret: credhub_cfg["secret"],
			CaCert: credhub_cfg["ca_cert"],
		},
	}

	chcli, err := credhubcli.New(
		cfg.Credhub.Server,
		credhubcli.SkipTLSValidation(true), // TODO use CA
		credhubcli.Auth(auth.UaaClientCredentials(cfg.Credhub.Client, cfg.Credhub.Secret)),
	)
	if err != nil {
		logger.Fatalf("failed to connect to Credhub: %s", err)
	}

	credhub = ccredhub.NewCredHub(chcli)

	director, err = cbosh.NewDirector(cfg.Bosh)
	if err != nil {
		logger.Fatalf("failed to connect to BOSH Director: %s", err)
	}

	state = NewState()
}

func refresh() error {
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
			logger.Fatalf("failed to load credentials from Credhub: %s", err)
		}

	}(&wg, &credentials)

	wg.Add(1)
	go func(wg *sync.WaitGroup, variables *[]*cbosh.Variable) {
		defer wg.Done()
		*variables, err = director.GetVariables()
		if err != nil {
			logger.Fatalf("failed to load variables from BOSH Director: %s", err)
		}

	}(&wg, &variables)

	wg.Wait()

	return state.Update(credentials, variables)
}
