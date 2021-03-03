package app

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/starkandwayne/carousel/credhub"
	"github.com/starkandwayne/carousel/state"
	"gopkg.in/yaml.v2"

	"github.com/grantae/certinfo"

	humanize "github.com/dustin/go-humanize"
)

func (a *Application) viewDetails() *tview.Flex {
	return tview.NewFlex()
}

func (a *Application) actionShowDetails(ref interface{}) {
	a.layout.details.Clear().AddItem(a.renderDetailsFor(ref), 0, 1, false)
}

func (a *Application) actionToggleTransitional(cred *state.Credential) {
	message := fmt.Sprintf(
		"Set transitional=%s for %s@%s\nWarning manual updates can lead to an inconsitent state",
		strconv.FormatBool(!cred.Transitional), cred.Name, cred.ID)

	a.renderModalAction(message, "Updating Transitional...", func() error {
		return a.credhub.UpdateTransitional(cred.Credential, cred.Transitional)
	})
}

func (a *Application) actionDelete(cred *state.Credential) {
	message := fmt.Sprintf("Delete %s@%s\nWarning manual deletion can lead to an inconsitent state",
		cred.Name, cred.ID)

	a.renderModalAction(message, "Deleting...", func() error {
		a.selectedID = cred.Path.Name
		return a.credhub.Delete(cred.Credential)
	})
}

func (a *Application) actionRegenerate(cred *state.Credential) {
	message := fmt.Sprintf("Regenerate %s@%s\nWarning manual regeneration can lead to an inconsitent state",
		cred.Name, cred.ID)

	a.renderModalAction(message, "Regenerating...", func() error {
		return a.credhub.ReGenerate(cred.Credential)
	})
}

func (a *Application) actionRefresh() {
	a.refreshWithStatusModal()
}

func (a *Application) renderDetailsFor(ref interface{}) tview.Primitive {
	switch v := ref.(type) {
	case *state.Path:
		return a.renderPathDetail(v)
	case *state.Credential:
		return a.renderCredentialDetail(v)
	default:
		return a.renderWelcome()
	}
}

func (a *Application) renderPathDetail(p *state.Path) tview.Primitive {
	t := tview.NewTable()
	t.SetBorder(true)
	t.SetTitle("Credhub & BOSH")

	addSimpleRow(t, "Name", p.Name)

	variableDef, err := yaml.Marshal(p.VariableDefinition)
	if err != nil {
		a.fatalf("Failed to marshal bosh variable definition got: %s", err)
	}

	info := tview.NewTextView().SetText(string(variableDef)).
		SetTextColor(tcell.Color102)

	info.SetBorder(true)
	info.SetTitle("BOSH variable definition")

	a.layout.tree.SetInputCapture(a.nextFocusInputCaptureHandler(t))
	t.SetInputCapture(a.nextFocusInputCaptureHandler(a.layout.tree))

	return tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(t, 3, 1, false).
		AddItem(a.renderPathActions(p), 1, 1, false).
		AddItem(info, 0, 1, true)
}

func (a *Application) renderCredentialDetail(cred *state.Credential) tview.Primitive {
	t := tview.NewTable()
	t.SetBorder(true)
	t.SetTitle("Credhub & BOSH")

	addSimpleRow(t, "ID", cred.ID)
	addSimpleRow(t, "Created At", fmt.Sprintf("%s (%s)",
		cred.VersionCreatedAt.Format(time.RFC3339),
		humanize.RelTime(*cred.VersionCreatedAt, time.Now(), "ago", "from now")))
	addSimpleRow(t, "Deployments", renderDeployments(cred.Deployments))
	addSimpleRow(t, "Latest", strconv.FormatBool(cred.Latest))

	var info *tview.TextView
	detailRows := 3 + 2 // 2 for top and bottom border
	detailRows = detailRows + len(cred.Deployments)

	switch cred.Type {
	case credhub.Certificate:
		addSimpleRow(t, "Expiry", fmt.Sprintf("%s (%s)",
			cred.ExpiryDate.Format(time.RFC3339),
			humanize.RelTime(*cred.ExpiryDate, time.Now(), "ago", "from now")))
		addSimpleRow(t, "Transitional", strconv.FormatBool(cred.Transitional))
		addSimpleRow(t, "Certificate Authority", strconv.FormatBool(cred.CertificateAuthority))
		addSimpleRow(t, "Self Signed", strconv.FormatBool(cred.SelfSigned))
		addSimpleRow(t, "Referenced by leafs", renderCredentials(cred.ReferencedBy))
		addSimpleRow(t, "Referenced CA's", renderCredentials(cred.References))

		detailRows = detailRows + 6

		i, err := certinfo.CertificateText(cred.Certificate)
		if err != nil {
			a.fatalf("Failed to render certificate to text got: %s", err)
		}

		info = tview.NewTextView().SetText(i).
			SetTextColor(tcell.Color102)
		info.SetBorder(true)
		info.SetTitle("Raw Certificate")
	default:
		info = tview.NewTextView().SetText("TODO").
			SetTextColor(tcell.Color102)
		info.SetBorder(true)
		info.SetTitle("Info")
	}

	a.layout.tree.SetInputCapture(a.nextFocusInputCaptureHandler(t))
	t.SetInputCapture(a.nextFocusInputCaptureHandler(info))
	info.SetInputCapture(a.nextFocusInputCaptureHandler(a.layout.tree))

	return tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(t, detailRows, 1, false).
		AddItem(a.renderCredentialActions(cred), 1, 1, false).
		AddItem(info, 0, 1, true)
}

func (a *Application) renderCredentialActions(cred *state.Credential) tview.Primitive {
	actions := []string{
		"Refresh State",
		"Toggle Transitional",
		"Delete",
	}

	out := []string{}
	for _, lbl := range actions {
		out = append(out, fmt.Sprintf("[yellow]^%s[white] %s",
			string([]rune(lbl)[0]), lbl))
	}

	a.keyBindings[tcell.KeyCtrlR] = a.actionRefresh

	a.keyBindings[tcell.KeyCtrlT] = func() {
		a.actionToggleTransitional(cred)
	}

	a.keyBindings[tcell.KeyCtrlD] = func() {
		a.actionDelete(cred)
	}

	return tview.NewTextView().
		SetDynamicColors(true).
		SetText(" " + strings.Join(out, "  "))
}

func (a *Application) renderPathActions(p *state.Path) tview.Primitive {
	actions := []string{
		"Refresh State",
		"Generate",
	}

	out := []string{}
	for _, lbl := range actions {
		out = append(out, fmt.Sprintf("[yellow]^%s[white] %s",
			string([]rune(lbl)[0]), lbl))
	}

	a.keyBindings[tcell.KeyCtrlR] = a.actionRefresh

	a.keyBindings[tcell.KeyCtrlG] = func() {
		a.actionRegenerate(p.Versions[0])
	}

	return tview.NewTextView().
		SetDynamicColors(true).
		SetText(" " + strings.Join(out, "  "))
}

func (a *Application) renderWelcome() tview.Primitive {
	u := tview.NewTable()
	u.SetBorder(true)
	u.SetTitle("Usage")

	addSimpleRow(u, "Arrow Up/Down", "navigate the tree / scroll text panel")
	addSimpleRow(u, "Enter", "expand/collapse tree node")
	addSimpleRow(u, "Tab", "cycle trough panels")
	addSimpleRow(u, "Control", "modifier for actions ([yellow]^[white])")

	return u
}

func addSimpleRow(t *tview.Table, lbl, val string) {
	if val == "" {
		return
	}
	row := t.GetRowCount()
	t.SetCell(row, 0, tview.NewTableCell(lbl).SetStyle(tcell.Style{}.Bold(true)))
	t.SetCellSimple(row, 1, val)
}

func renderDeployments(deployments []*state.Deployment) string {
	tmp := make([]string, 0)
	for _, d := range deployments {
		tmp = append(tmp, d.Name)
	}

	return strings.Join(tmp, ", ")
}

func renderCredentials(credentials state.Credentials) string {
	tmp := make([]string, 0)
	for _, c := range credentials {
		tmp = append(tmp, c.ID)
	}

	return strings.Join(tmp, ", ")
}
