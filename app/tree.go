package app

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/starkandwayne/carousel/store"
)

const treePanel = "TreePanel"

func (a *Application) viewTree() *tview.TreeView {
	return tview.NewTreeView()
}

func (a *Application) renderTree() {
	root := tview.NewTreeNode("∎")
	a.store.EachCert(func(cert *store.Cert) {
		// only interested in top level certVersions
		if len(cert.Versions) != 0 && cert.Versions[0].SignedBy != nil {
			return
		}
		root.SetChildren(append(root.GetChildren(), addToTree(cert.Versions)...))
	})

	var currentNode *tview.TreeNode

	if a.selectedID == "" {
		currentNode = root
	}

	root.Walk(func(node, parent *tview.TreeNode) bool {
		if currentNode != nil {
			return false
		}
		if refToID(node.GetReference()) == a.selectedID {
			currentNode = node
			a.actionShowDetails(currentNode.GetReference())
			return false
		}
		return true
	})

	a.layout.tree.SetRoot(root).SetCurrentNode(currentNode)

	a.layout.tree.SetChangedFunc(func(node *tview.TreeNode) {
		a.selectedID = refToID(node.GetReference())
		a.actionShowDetails(node.GetReference())
	})

	a.layout.tree.SetSelectedFunc(func(node *tview.TreeNode) {
		node.SetExpanded(!node.IsExpanded())
	})
}

func refToID(ref interface{}) string {
	switch v := ref.(type) {
	case *store.Cert:
		return v.Id
	case *store.CertVersion:
		return v.Id
	default:
		return ""
	}
}

func addToTree(certVersions []*store.CertVersion) []*tview.TreeNode {
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
		lbl := fmt.Sprintf("%s (%s)", certVersion.Id, certVersion.Status())
		if certVersion.Transitional {
			lbl = lbl + " (transitional)"
		}
		certVersionNode := tview.NewTreeNode(lbl).SetReference(certVersion)
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
