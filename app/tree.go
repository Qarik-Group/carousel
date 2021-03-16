package app

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/starkandwayne/carousel/credhub"
	"github.com/starkandwayne/carousel/state"
)

const treePanel = "TreePanel"

func (a *Application) viewTree() *tview.TreeView {
	return tview.NewTreeView()
}

func (a *Application) updateTree() {
	root := tview.NewTreeNode("∎")
	root.SetReference("root")
	a.expanded[refToID(root.GetReference())] = true

	for _, credType := range credhub.CredentialTypeValues() {
		credentials := a.state.Credentials(
			state.TypeFilter(credType),
			state.SelfSignedFilter(),
			state.LatestFilter())

		if len(credentials) == 0 {
			continue
		}

		typeNode := tview.NewTreeNode(credType.String()).SetReference(credType.String())

		if expanded, ok := a.expanded[refToID(typeNode.GetReference())]; ok {
			typeNode.SetExpanded(expanded)
		} else {
			typeNode.SetExpanded(false)
		}

		root.AddChild(typeNode)

		for _, credential := range credentials {
			typeNode.SetChildren(append(typeNode.GetChildren(), addToTree(credential.Path.Versions)...))
		}
	}

	var currentNode *tview.TreeNode

	if a.selectedID == "" {
		currentNode = root
	}

	root.Walk(func(node, parent *tview.TreeNode) bool {
		if refToID(node.GetReference()) == a.selectedID {
			currentNode = node
			a.actionShowDetails(currentNode.GetReference())
		}

		if expanded, ok := a.expanded[refToID(node.GetReference())]; ok {
			node.SetExpanded(expanded)
		} else {
			node.SetExpanded(false)
		}

		return true
	})

	if currentNode == nil {
		currentNode = root
	}

	a.layout.tree.SetRoot(root).SetCurrentNode(currentNode)

	a.layout.tree.SetChangedFunc(func(node *tview.TreeNode) {
		a.selectedID = refToID(node.GetReference())
		a.actionShowDetails(node.GetReference())
	})

	a.layout.tree.SetSelectedFunc(func(node *tview.TreeNode) {
		a.expanded[refToID(node.GetReference())] = !node.IsExpanded()
		node.SetExpanded(!node.IsExpanded())
	})
}

func refToID(ref interface{}) string {
	switch v := ref.(type) {
	case *state.Credential:
		return v.ID
	case *state.Path:
		return v.Name
	case string:
		return v
	default:
		return "none"
	}
}

func addToTree(creds []*state.Credential) []*tview.TreeNode {
	out := make([]*tview.TreeNode, 0)
	for _, cred := range creds {
		pathNode := tview.NewTreeNode(cred.Path.Name).
			SetReference(cred.Path).Collapse()

		var exists bool
		for _, n := range out {
			if refToID(n.GetReference()) == refToID(cred.Path) {
				exists = true
				pathNode = n
			}
		}

		lbl := fmt.Sprintf("%s (%s)", cred.ID, toStatus(cred))
		if cred.Transitional {
			lbl = lbl + " (transitional)"
		}
		credNode := tview.NewTreeNode(lbl).SetReference(cred)
		switch toStatus(cred) {
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
