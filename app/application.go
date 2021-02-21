package app

import (
	"github.com/starkandwayne/carousel/credhub"
	"github.com/starkandwayne/carousel/state"

	"github.com/gdamore/tcell/v2"

	"github.com/rivo/tview"
)

type Application struct {
	*tview.Application
	state       state.State
	credhub     credhub.CredHub
	layout      *Layout
	keyBindings map[tcell.Key]func()
	selectedID  string
	refresh     func()
}

type Layout struct {
	main    *tview.Flex
	tree    *tview.TreeView
	details *tview.Flex
}

func NewApplication(state state.State, ch credhub.CredHub, refresh func()) *Application {
	return &Application{
		Application: tview.NewApplication(),
		state:       state,
		keyBindings: make(map[tcell.Key]func(), 0),
		credhub:     ch,
		refresh:     refresh,
	}
}

func (a *Application) Init() *Application {
	a.layout = &Layout{
		tree:    a.viewTree(),
		details: a.viewDetails(),
	}

	flex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewFlex().
				AddItem(a.layout.tree, 0, 1, false).
				AddItem(a.layout.details, 0, 1, false),
				0, 5, true),
			0, 1, false)

	a.layout.main = flex

	a.SetRoot(flex, true)
	a.SetFocus(a.layout.tree)
	a.EnableMouse(false)

	a.renderTree()
	a.actionShowDetails(nil)

	a.initGlobalKeyInputCaputreHandler()

	return a
}

func (a *Application) nextFocusInputCaptureHandler(p tview.Primitive) func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			a.SetFocus(p)
		}
		return event
	}
}

func (a *Application) initGlobalKeyInputCaputreHandler() {
	a.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		for k, fn := range a.keyBindings {
			if event.Key() == k {
				fn()
				return nil
			}
		}
		return event
	})
}

func (a *Application) statusModal(status string) {
	a.SetRoot(tview.NewModal().SetText(status), true)
	a.ForceDraw()
}
