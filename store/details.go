package store

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/grantae/certinfo"

	humanize "github.com/dustin/go-humanize"
)

func (s *Store) ShowDetails(node *tview.TreeNode) *tview.Frame {
	switch v := node.GetReference().(type) {
	case *Cert:
		return renderCertDetail(v)
	case *CertVersion:
		return renderCertVersionDetail(v)
	default:
		return renderHelp()
	}
}

func renderCertDetail(c *Cert) *tview.Frame {
	t := tview.NewTable()
	addSimpleRow(t, "ID", c.Id)
	addSimpleRow(t, "Name", c.Name)
	return tview.NewFrame(t)
}

func renderCertVersionDetail(cv *CertVersion) *tview.Frame {
	t := tview.NewTable()
	t.SetBorder(true)
	t.SetTitle("Credhub & BOSH")

	addSimpleRow(t, "ID", cv.Id)
	addSimpleRow(t, "Expiry", fmt.Sprintf("%s (%s)",
		cv.Expiry.Format(time.RFC3339),
		humanize.RelTime(cv.Expiry, time.Now(), "ago", "from now")))
	addSimpleRow(t, "Transitional", strconv.FormatBool(cv.Transitional))
	addSimpleRow(t, "Certificate Authority", strconv.FormatBool(cv.CertificateAuthority))
	addSimpleRow(t, "Self Signed", strconv.FormatBool(cv.SelfSigned))

	addSimpleRow(t, "Deployments", renderDeployments(cv.Deployments))
	// skid := make([]byte, hex.DecodedLen(len(cv.Certificate.SubjectKeyId)))
	// n, _ := hex.Decode(skid, cv.Certificate.SubjectKeyId)
	// addSimpleRow(t, "SKID", string(skid[:n]))
	// addSimpleRow(t, "DNS Names", strings.Join(cv.Certificate.DNSNames, ", "))

	i, err := certinfo.CertificateText(cv.Certificate)
	if err != nil {
		panic(err)
	}

	info := tview.NewTextView().SetText(i).
		SetTextColor(tcell.Color102)

	info.SetBorder(true)
	info.SetTitle("Raw Certificate")

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(t, 8, 1, false).
		AddItem(info, 0, 1, false)

	return tview.NewFrame(flex)
}

func renderHelp() *tview.Frame {
	return tview.NewFrame(tview.NewBox().SetBorder(true).SetTitle("help"))
}

func addSimpleRow(t *tview.Table, lbl, val string) {
	if val == "" {
		return
	}
	row := t.GetRowCount()
	t.SetCell(row, 0, tview.NewTableCell(lbl).SetStyle(tcell.Style{}.Bold(true)))
	t.SetCellSimple(row, 1, val)
}

func renderDeployments(deployments []*Deployment) string {
	tmp := make([]string, 0)
	for _, d := range deployments {
		tmp = append(tmp, d.Name)
	}

	return strings.Join(tmp, ", ")
}
