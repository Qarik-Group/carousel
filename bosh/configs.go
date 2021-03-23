package bosh

import (
	"path"
	"strconv"

	boshdir "github.com/cloudfoundry/bosh-cli/director"
)

func (d *director) GetActiveCloudConfigs(deployment string) (map[string][]byte, error) {
	return d.getConfigs(deployment, "cloud", true)
}
func (d *director) GetLatestCloudConfigs(deployment string) (map[string][]byte, error) {
	return d.getConfigs(deployment, "cloud", false)
}
func (d *director) GetActiveRuntimeConfigs(deployment string) (map[string][]byte, error) {
	return d.getConfigs(deployment, "runtime", true)
}
func (d *director) GetLatestRuntimeConfigs(deployment string) (map[string][]byte, error) {
	return d.getConfigs(deployment, "runtime", false)
}

func (d *director) getConfigs(deployment, configType string, active bool) (map[string][]byte, error) {
	configs, err := d.client.ListDeploymentConfigs(deployment)
	if err != nil {
		return nil, err
	}

	out := make(map[string][]byte, 0)
	for _, conf := range configs.GetConfigs() {
		if conf.Type == configType {
			var c boshdir.Config
			if active {
				c, err = d.client.LatestConfigByID(strconv.Itoa(conf.Id))
			} else {
				c, err = d.client.LatestConfig(conf.Type, conf.Name)
			}
			if err != nil {
				return nil, err
			}
			out[path.Join(c.Type, c.Name)] = []byte(c.Content)

		}
	}

	return out, nil
}
