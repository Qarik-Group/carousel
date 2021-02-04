package main

import "github.com/rivo/tview"

func addToTree(root []*tview.TreeNode, names []string) []*tview.TreeNode {
	if len(names) > 0 {
		var i int
		for i = 0; i < len(root); i++ {
			if root[i].GetText() == names[0] { //already in tree
				break
			}
		}
		if i == len(root) {
			root = append(root, tview.NewTreeNode(names[0]))
		}
		root[i].SetChildren(addToTree(root[i].GetChildren(), names[1:]))
	}
	return root
}
