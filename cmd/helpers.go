package cmd

import (
	"log"
	"os"

	credhubcli "code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/auth"
	"github.com/starkandwayne/carousel/config"
	ccredhub "github.com/starkandwayne/carousel/credhub"
	cstore "github.com/starkandwayne/carousel/store"

	boshdir "github.com/cloudfoundry/bosh-cli/director"
	boshuaa "github.com/cloudfoundry/bosh-cli/uaa"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

var (
	logger  *log.Logger
	credhub ccredhub.CredHub
	store   *cstore.Store
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

	authURL, err := chcli.AuthURL()
	if err != nil {
		logger.Fatalf("failed to lookup auth url: %s", err)
	}

	uaa, err := buildUAA(cfg.Bosh, authURL)
	if err != nil {
		logger.Fatalf("failed to initialize uaa client: %s", err)
	}

	d, err := buildDirector(cfg.Bosh, uaa)
	if err != nil {
		logger.Fatalf("failed to initialize bosh director client: %s", err)
	}

	credhub = ccredhub.NewCredHub(chcli)

	store = cstore.NewStore(credhub, d)
}

func buildUAA(cfg *config.Bosh, authURL string) (boshuaa.UAA, error) {
	logger := boshlog.NewLogger(boshlog.LevelError)
	factory := boshuaa.NewFactory(logger)

	// Build a UAA config from a URL.
	// HTTPS is required and certificates are always verified.
	config, err := boshuaa.NewConfigFromURL(authURL)
	if err != nil {
		return nil, err
	}

	// Set client credentials for authentication.
	// Machine level access should typically use a client instead of a particular user.
	config.Client = cfg.Client
	config.ClientSecret = cfg.ClientSecret

	// Configure trusted CA certificates.
	// If nothing is provided default system certificates are used.
	config.CACert = cfg.CaCert

	return factory.New(config)
}

func buildDirector(cfg *config.Bosh, uaa boshuaa.UAA) (boshdir.Director, error) {
	logger := boshlog.NewLogger(boshlog.LevelError)
	factory := boshdir.NewFactory(logger)

	// Build a Director config from address-like string.
	// HTTPS is required and certificates are always verified.
	factoryConfig, err := boshdir.NewConfigFromURL(cfg.Environment)
	if err != nil {
		return nil, err
	}

	// Configure custom trusted CA certificates.
	// If nothing is provided default system certificates are used.
	factoryConfig.CACert = cfg.CaCert

	// Allow Director to fetch UAA tokens when necessary.
	factoryConfig.TokenFunc = boshuaa.NewClientTokenSession(uaa).TokenFunc

	return factory.New(factoryConfig, boshdir.NewNoopTaskReporter(), boshdir.NewNoopFileReporter())
}
