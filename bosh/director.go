package bosh

import (
	"github.com/starkandwayne/carousel/config"

	boshdir "github.com/cloudfoundry/bosh-cli/director"
	boshuaa "github.com/cloudfoundry/bosh-cli/uaa"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

type Director interface {
	GetVariables() ([]*Variable, error)
}

func NewDirector(cfg *config.Bosh) (Director, error) {
	dc, err := buildDirector(cfg)
	if err != nil {
		return nil, err
	}

	return &director{dc}, nil
}

type director struct {
	client boshdir.Director
}

func buildDirector(cfg *config.Bosh) (boshdir.Director, error) {
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

	noAuthDir, err := factory.New(factoryConfig, boshdir.NewNoopTaskReporter(), boshdir.NewNoopFileReporter())
	if err != nil {
		return nil, err
	}

	info, err := noAuthDir.Info()
	if err != nil {
		return nil, err
	}

	uaa, err := buildUAA(cfg, info.Auth.Options["url"].(string))
	if err != nil {
		return nil, err
	}

	// Allow Director to fetch UAA tokens when necessary.
	factoryConfig.TokenFunc = boshuaa.NewClientTokenSession(uaa).TokenFunc

	return factory.New(factoryConfig, boshdir.NewNoopTaskReporter(), boshdir.NewNoopFileReporter())
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
