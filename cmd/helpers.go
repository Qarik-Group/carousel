package cmd

import (
	"log"
	"os"
	"sync"

	credhubcli "code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/auth"
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

func initialize() {
	logger = log.New(os.Stderr, "", 0)
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatalf("failed to load environment configuration: %s", err)
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

func refresh() {
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

	err = state.Update(credentials, variables)
	if err != nil {
		logger.Fatalf("failed to update state got: %s", err)
	}
}
