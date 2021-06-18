package state

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/starkandwayne/carousel/bosh"
	"github.com/starkandwayne/carousel/credhub"
)

func (c *Credential) MarshalJSON() ([]byte, error) {
	deployments := make([]string, 0)
	for _, d := range c.Deployments {
		deployments = append(deployments, d.Name)
	}

	updateMode := bosh.NoOverwrite
	if c.Path.VariableDefinition != nil {
		updateMode = c.Path.VariableDefinition.UpdateMode
	}

	c.RawValue = nil // don't leak raw value
	type Alias Credential
	return json.Marshal(&struct {
		*Alias
		DeploymentsList []string        `json:"deployments"`
		UpdateMode      bosh.UpdateMode `json:"update_mode"`
	}{
		Alias:           (*Alias)(c),
		DeploymentsList: deployments,
		UpdateMode:      updateMode,
	})
}

func (c *Credential) PathVersion() string {
	return fmt.Sprintf("%s@%s", c.Name, c.ID)
}

func (c *Credential) PendingDeploys() Deployments {
	out := make(Deployments, 0)
	if !c.Latest {
		return out
	}
	for _, d := range c.Path.Deployments {
		if !c.Deployments.Includes(d) {
			out = append(out, d)
		}
	}
	return out
}

func (c *Credential) Active() bool {
	if len(c.Deployments) != 0 {
		return true
	}
	for _, cred := range c.ReferencedBy {
		if cred == c {
			// Handle circular references
			continue
		}
		if cred.Active() {
			return true
		}
	}
	return false
}

func (c *Credential) LatestDeployedTo(deployment string) *Credential {
	for _, version := range c.Path.Versions {
		if version.Deployments.IncludesName(deployment) {
			return version
		}
	}
	return nil
}

func (c *Credential) PrintCreatedAt() string {
	return fmt.Sprintf("%s (%s)",
		c.VersionCreatedAt.Format(time.RFC3339),
		humanize.RelTime(*c.VersionCreatedAt, time.Now(), "ago", "from now"),
	)
}

func (c *Credential) PrintExpiry() string {
	if c.ExpiryDate == nil {
		return ""
	}
	return fmt.Sprintf("%s (%s)",
		c.ExpiryDate.Format(time.RFC3339),
		humanize.RelTime(*c.ExpiryDate, time.Now(), "ago", "from now"),
	)
}

func (c *Credential) Summary() string {
	switch c.Type {
	case credhub.Certificate:
		return fmt.Sprintf("type: %s | created at: %s | expiry: %s",
			c.Type.String(), c.PrintCreatedAt(), c.PrintExpiry())
	default:
		return fmt.Sprintf("type: %s | created at: %s",
			c.Type.String(), c.PrintCreatedAt())
	}
}
