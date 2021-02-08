package main

import (
	"log"
	"os"

	"code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/auth"
	"github.com/starkandwayne/carousel/store"

	"github.com/rivo/tview"
)

func main() {
	logger := log.New(os.Stderr, "", 0)
	cfg, err := loadConfig()
	if err != nil {
		logger.Fatalf("failed to load environment configuration: %s", err)
	}

	ch, err := credhub.New(
		cfg.Credhub.Server,
		credhub.SkipTLSValidation(true), // TODO use CA
		credhub.Auth(auth.UaaClientCredentials(cfg.Credhub.Client, cfg.Credhub.Secret)),
	)
	if err != nil {
		logger.Fatalf("failed to connect to Credhub: %s", err)
	}

	s, err := store.NewStore(ch)
	if err != nil {
		logger.Fatalf("failed to load data: %s", err)
	}

	root := s.Tree()

	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	details := tview.NewFlex().
		AddItem(tview.NewBox().SetBorder(true).SetTitle("welcome"), 0, 1, false)

	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		// TODO switch detailed view on switching to a node (not only when hitting enter)
		details.Clear().AddItem(s.ShowDetails(node), 0, 1, false)
		node.SetExpanded(!node.IsExpanded())
	})

	flex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			//	AddItem(tview.NewBox().SetBorder(true).SetTitle("Controls"), 0, 1, false).
			AddItem(tview.NewFlex().
				AddItem(tree, 0, 1, true).
				AddItem(details, 0, 1, false),
				0, 5, true), //.
			//			AddItem(tview.NewBox().SetBorder(true).SetTitle("More Controls"), 0, 1, false),
			0, 1, false)

	if err := tview.NewApplication().SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
