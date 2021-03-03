package app

import (
	"fmt"
	"os"
	"time"

	"github.com/rivo/tview"
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
	if !c.Active() {
		status = "unused"
	}
	return status
}

func (a *Application) statusModal(status string) {
	a.SetRoot(tview.NewModal().SetText(status), true)
	a.ForceDraw()
}

func (a *Application) renderModalAction(body, status string, fn func() error) {
	modal := tview.NewModal().
		SetText(body).AddButtons([]string{"Continue", "Cancel"})
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Continue" {
			a.statusModal(status)
			err := fn()
			if err != nil {
				a.fatalf("Failed got error: %s", err)
			}
			a.refreshWithStatusModal()
		}
		a.renderHome()
	})

	a.SetRoot(modal, true)
}

func (a *Application) renderHome() {
	a.SetRoot(a.layout.main, true)
	a.SetFocus(a.layout.tree)
}

func (a *Application) refreshWithStatusModal() {
	a.statusModal("Refreshing State...")
	err := a.refresh()
	if err != nil {
		a.fatalf("Failed to refresh got error: %s", err)
	}
	a.updateTree()
	a.renderHome()
}

func (a *Application) fatalf(msg string, args ...interface{}) {
	a.SetRoot(tview.NewBox(), true)
	a.ForceDraw()
	a.Stop()
	fmt.Printf(msg, args...)
	fmt.Println("")
	os.Exit(1)

}
