package store

import (
	"fmt"

	"github.com/rivo/tview"
)

func (s *Store) ShowDetails(node *tview.TreeNode) *tview.Frame {
	switch v := node.GetReference().(type) {
	case *Cert:
		//fmt.Println(v.Name)
		return tview.NewFrame(tview.NewBox().SetBorder(true).SetTitle(fmt.Sprintf("id: %s", v.Name)))
	case *CertVersion:
		//		fmt.Println(v.Id)
		return tview.NewFrame(tview.NewBox().SetBorder(true).SetTitle(fmt.Sprintf("id: %s", v.Id)))
		// default:
		//	return tview.NewFrame(tview.NewBox()).AddText(
		//		"select a node in the tree using arrow keys and the enter key. Alternatlivly use your mouse.")
	default:
		return tview.NewFrame(tview.NewBox().SetBorder(true).SetTitle(fmt.Sprintf("no type for: %+v", v)))
	}
}
