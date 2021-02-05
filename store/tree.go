package store

import (
	"fmt"

	"github.com/rivo/tview"
)

func (s *Store) Tree() *tview.TreeNode {
	root := tview.NewTreeNode("∎")
	it := s.Certs.Iterator()
	for it.Next() {
		k, v := it.Key(), it.Value()
		key := k.(string)
		value := v.(*Cert)

		node := tview.NewTreeNode(key)

		for _, version := range value.Versions {
			node.SetChildren(append(
				node.GetChildren(),
				tview.NewTreeNode(
					fmt.Sprintf("id:%s-%s", version.Id, version.Certificate.AuthorityKeyId),
				)))
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

// func (c *certsCache) addToTree(root []*tview.TreeNode, names []string) []*tview.TreeNode {
//	if len(names) > 0 {
//		var i int
//		for i = 0; i < len(root); i++ {
//			if root[i].GetText() == names[0] { //already in tree
//				break
//			}
//		}
//		if i == len(root) {
//			root = append(root, tview.NewTreeNode(names[0]))
//		}
//		root[i].SetChildren(addToTree(root[i].GetChildren(), names[1:]))
//	}
//	return root
// }
