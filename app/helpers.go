package app

import (
	"time"

	"github.com/starkandwayne/carousel/state"
)

func toStatus(c *state.Credential) string {
	status := "active"
	if c.ExpiryDate != nil && c.ExpiryDate.Sub(time.Now()) < time.Hour*24*30 {
		status = "notice"
	}
	if c.VersionCreatedAt.Sub(time.Now()) > time.Hour*24*365 {
		status = "notice"
	}
	if len(c.Deployments) == 0 {
		status = "unused"
	}
	return status
}
