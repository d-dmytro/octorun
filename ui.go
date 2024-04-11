package main

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type UI struct {
	app           *tview.Application
	navStatus     *tview.TextView
	nav           *tview.List
	logStatus     *tview.TextView
	logPages      *tview.Pages
	input         *tview.InputField
	scrollEnabled bool
}

func NewUI() *UI {
	ui := UI{
		scrollEnabled: true,
	}
	ui.initialize()
	return &ui
}

func (ui *UI) Run() error {
	return ui.app.Run()
}

func (ui *UI) AddLogPage(label string) *tview.TextView {
	textView := tview.NewTextView()
	textView.SetChangedFunc(func() {
		if ui.scrollEnabled {
			textView.ScrollToEnd()
		}
		ui.app.Draw()
	})
	index := ui.logPages.GetPageCount()
	shortcut := rune(strconv.Itoa(index + 1)[0])
	ui.logPages.AddPage(label, textView, true, index == 0)

	if index == 0 {
		ui.logStatus.Clear()
		fmt.Fprint(ui.logStatus, label)
	}

	ui.nav.AddItem(label, "", shortcut, func() {
		ui.logStatus.Clear()
		fmt.Fprint(ui.logStatus, label)
		ui.logPages.SwitchToPage(label)
	})

	return textView
}

func (ui *UI) initialize() {
	ui.app = tview.NewApplication()
	ui.navStatus = tview.NewTextView()
	ui.navStatus.SetChangedFunc(func() {
		ui.app.Draw()
	})
	fmt.Fprint(ui.navStatus, "Commands")
	ui.nav = tview.NewList()
	ui.nav.SetBorderPadding(1, 0, 0, 1)

	ui.logStatus = tview.NewTextView()
	ui.logStatus.SetBorderPadding(0, 0, 1, 0)
	ui.logStatus.SetChangedFunc(func() {
		ui.app.Draw()
	})
	ui.logPages = tview.NewPages()
	ui.logPages.SetBorderPadding(1, 0, 1, 0)

	ui.input = tview.NewInputField()
	ui.input.SetFieldBackgroundColor(tcell.NewRGBColor(117, 80, 123))

	ui.input.SetAcceptanceFunc(func(text string, lastChar rune) bool {
		return text != ":"
	})

	navFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(ui.navStatus, 1, 0, false).
		AddItem(ui.nav, 0, 1, false)

	logFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(ui.logStatus, 1, 0, false).
		AddItem(ui.logPages, 0, 1, false)

	grid := tview.NewGrid().SetRows(0, 1).SetColumns(25, 0).
		AddItem(navFlex, 0, 0, 1, 1, 0, 0, false).
		AddItem(logFlex, 0, 1, 1, 1, 0, 0, false).
		AddItem(ui.input, 1, 0, 1, 2, 0, 0, false)

	ui.app.SetRoot(grid, true).SetFocus(ui.nav)

	var prevFocusedEl tview.Primitive

	ui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlH:
			ui.scrollEnabled = true
			ui.app.SetFocus(ui.nav)
		case tcell.KeyCtrlL:
			ui.scrollEnabled = false
			if _, log := ui.logPages.GetFrontPage(); log != nil {
				ui.app.SetFocus(log)
			}
		case tcell.KeyRune:
			if event.Rune() == ColonRune {
				prevFocusedEl = ui.app.GetFocus()
				ui.app.SetFocus(ui.input)
			}
		case tcell.KeyEsc:
			if ui.app.GetFocus() == ui.input {
				ui.input.SetText("")
				ui.app.SetFocus(prevFocusedEl)
			}
		}

		return event
	})
}
