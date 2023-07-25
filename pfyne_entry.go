package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type PCSEntry struct {
	widget.Entry
	m_funSave func()
}

func newPCSEntry(funSave func()) *PCSEntry {
	p := &PCSEntry{
		m_funSave: funSave,
	}
	p.MultiLine = true
	p.Wrapping = fyne.TextTruncate
	p.ExtendBaseWidget(p)

	return p
}

func (p *PCSEntry) TypedShortcut(s fyne.Shortcut) {
	cs, ok := s.(*desktop.CustomShortcut)
	if !ok {
		p.Entry.TypedShortcut(s)
		return
	}

	if cs.Modifier == fyne.KeyModifierControl && cs.Key() == fyne.KeyS { // ctrl + s
		p.m_funSave()
	}
}
