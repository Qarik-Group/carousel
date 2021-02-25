package cmd

import (
	"github.com/vito/go-interact/interact"
)

func askForConfirmation() {
	if nonInteractive {
		return
	}

	falseByDefault := false

	err := interact.NewInteraction("Continue?").Resolve(&falseByDefault)
	if err != nil {
		logger.Fatalf("Asking for confirmation: %s", err)
	}

	if falseByDefault == false {
		logger.Fatal("Stopped")
	}
}
