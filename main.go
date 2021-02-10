package main

import (
	"log"
	"os"

	"code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/auth"
	"github.com/starkandwayne/carousel/app"
	"github.com/starkandwayne/carousel/store"

	boshdir "github.com/cloudfoundry/bosh-cli/director"
	boshuaa "github.com/cloudfoundry/bosh-cli/uaa"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

func main() {
	logger := log.New(os.Stderr, "", 0)
	cfg, err := loadConfig()
	if err != nil {
		logger.Fatalf("failed to load environment configuration: %s", err)
	}

	ch, err := credhub.New(
		cfg.Credhub.Server,
		credhub.SkipTLSValidation(true), // TODO use CA
		credhub.Auth(auth.UaaClientCredentials(cfg.Credhub.Client, cfg.Credhub.Secret)),
	)
	if err != nil {
		logger.Fatalf("failed to connect to Credhub: %s", err)
	}

	authURL, err := ch.AuthURL()
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

	s, err := store.NewStore(ch, d)
	if err != nil {
		logger.Fatalf("failed to load data: %s", err)
	}

	app := app.NewApplication(s).Init()

	if err := app.Run(); err != nil {
		panic(err)
	}
}

func buildUAA(cfg *Bosh, authURL string) (boshuaa.UAA, error) {
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

func buildDirector(cfg *Bosh, uaa boshuaa.UAA) (boshdir.Director, error) {
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
