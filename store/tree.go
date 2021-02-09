package store

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (s *Store) Tree() *tview.TreeNode {
	root := tview.NewTreeNode("∎")
	it := s.certs.Iterator()
	for it.Next() {
		_, v := it.Key(), it.Value()
		cert := v.(*Cert)

		// only interested in top level certVersions
		if len(cert.Versions) != 0 && cert.Versions[0].SignedBy != nil {
			continue
		}
		root.SetChildren(append(root.GetChildren(), addToTree(cert.Versions)...))
	}

	return root
}

func addToTree(certVersions []*CertVersion) []*tview.TreeNode {
	out := make([]*tview.TreeNode, 0)
	for _, certVersion := range certVersions {
		certNode := tview.NewTreeNode(certVersion.Cert.Name).
			SetReference(certVersion.Cert)

		var exists bool
		for _, n := range out {
			if n.GetText() == certVersion.Cert.Name {
				exists = true
				certNode = n
			}
		}
		certVersionNode := tview.NewTreeNode(fmt.Sprintf("%s (%s)", certVersion.Id, certVersion.Status())).SetReference(certVersion)
		switch certVersion.Status() {
		case "unused":
			certVersionNode.SetColor(tcell.Color102)
		case "notice":
			certVersionNode.SetColor(tcell.ColorDarkGoldenrod)
		}
		certNode.AddChild(certVersionNode)
		certVersionNode.SetChildren(addToTree(certVersion.Signs))
		if !exists {
			out = append(out, certNode)
		}
	}
	return out
}
