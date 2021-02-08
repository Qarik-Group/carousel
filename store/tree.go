package store

import (
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

		node := tview.NewTreeNode(cert.Name).
			SetReference(cert).
			Collapse()

		for _, version := range cert.Versions {
			node.SetChildren(append(
				node.GetChildren(),
				tview.NewTreeNode(version.Id).
					SetReference(version).
					SetChildren(addToTree(version.Signs)),
			))
		}

		root.SetChildren(append(root.GetChildren(), node))
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
		certVersionNode := tview.NewTreeNode(certVersion.Id).
			SetReference(certVersion)
		certNode.AddChild(certVersionNode)
		certVersionNode.SetChildren(addToTree(certVersion.Signs))
		if !exists {
			out = append(out, certNode)
		}
	}
	return out
}
