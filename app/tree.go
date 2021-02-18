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
	a.store.EachPath(func(path *store.Path) {
		// only interested in top level certVersions
		if len(path.Versions) != 0 && path.Versions[0].SignedBy != nil {
			return
		}
		root.SetChildren(append(root.GetChildren(), addToTree(path.Versions)...))
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
	case *store.Credential:
		return v.ID
	case *store.Path:
		return v.Name
	default:
		return ""
	}
}

func addToTree(creds []*store.Credential) []*tview.TreeNode {
	out := make([]*tview.TreeNode, 0)
	for _, cred := range creds {
		pathNode := tview.NewTreeNode(cred.Path.Name).
			SetReference(cred.Path)

		var exists bool
		for _, n := range out {
			if refToID(n.GetReference()) == refToID(cred.Path) {
				exists = true
				pathNode = n
			}
		}

		lbl := fmt.Sprintf("%s (%s)", cred.ID, cred.Status())
		if cred.Transitional {
			lbl = lbl + " (transitional)"
		}
		credNode := tview.NewTreeNode(lbl).SetReference(cred)
		switch cred.Status() {
		case "unused":
			credNode.SetColor(tcell.Color102)
		case "notice":
			credNode.SetColor(tcell.ColorDarkGoldenrod)
		}
		pathNode.AddChild(credNode)
		credNode.SetChildren(addToTree(cred.Signs))
		if !exists {
			out = append(out, pathNode)
		}
	}
	return out
}
