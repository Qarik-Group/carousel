package app

import (
	"github.com/starkandwayne/carousel/store"

	"github.com/gdamore/tcell/v2"

	"github.com/rivo/tview"
)

type Application struct {
	*tview.Application
	store  *store.Store
	layout *Layout
}

type Layout struct {
	tree    *tview.TreeView
	details *tview.Flex
}

func NewApplication(store *store.Store) *Application {
	return &Application{
		Application: tview.NewApplication(),
		store:       store,
	}
}

func (a *Application) Init() *Application {
	a.layout = &Layout{
		tree:    a.viewTree(),
		details: a.viewDetails(),
	}

	flex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Controls"), 3, 0, false).
			AddItem(tview.NewFlex().
				AddItem(a.layout.tree, 0, 1, false).
				AddItem(a.layout.details, 0, 1, false),
				0, 5, true).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("More Controls"), 3, 0, false),
			0, 1, false)

	a.SetRoot(flex, true)
	a.SetFocus(a.layout.tree)
	a.EnableMouse(false)

	a.actionShowDetails(nil)

	return a
}

func (a *Application) nextFocusIncputCaptureHandler(p tview.Primitive) func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			a.SetFocus(p)
		}
		return event
	}
}
