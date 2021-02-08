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

		node := tview.NewTreeNode(cert.Name)
		for _, version := range cert.Versions {
			node.SetChildren(append(
				node.GetChildren(),
				tview.NewTreeNode(version.Id).SetChildren(addToTree(version.Signs)),
			))
		}

		root.SetChildren(append(root.GetChildren(), node))
	}

	// for _, cert := range certs {
	//	if cert.Name == cert.SignedBy {
	//		root.SetChildren(addToTree(root.GetChildren(), []string{cert.Name}))
	//	} else {
	//		root.SetChildren(addToTree(root.GetChildren(), []string{cert.SignedBy, cert.Name}))
	//	}
	// }

	return root
}

// func (c *certsCache) withNames(names []string) []*credentials.CertificateMetadata {
//	out := make([]*credentials.CertificateMetadata, 0)
//	for _, cert := range c.certs {
//		for _, name := range names {
//			if name == cert.Name {
//				out = append(out, &cert)
//			}
//		}
//	}

// }

func addToTree(certVersions []*CertVersion) []*tview.TreeNode {
	out := make([]*tview.TreeNode, 0)
	for _, certVersion := range certVersions {
		certNode := tview.NewTreeNode(certVersion.Cert.Name)
		var exists bool
		for _, n := range out {
			if n.GetText() == certVersion.Cert.Name {
				exists = true
				certNode = n
			}
		}
		certVersionNode := tview.NewTreeNode(certVersion.Id)
		certNode.AddChild(certVersionNode)
		certVersionNode.SetChildren(addToTree(certVersion.Signs))
		if !exists {
			out = append(out, certNode)
		}
	}
	return out
}
