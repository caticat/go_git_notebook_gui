package main

import (
	"path"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/caticat/go_game_server/phelp/ppath"
)

type PFyneTreePath struct {
	*ppath.PPath
}

func NewPFyneTreePath(basePath string) *PFyneTreePath {
	return &PFyneTreePath{
		ppath.NewPPath(basePath),
	}
}

func (p *PFyneTreePath) ChildUIDs(id widget.TreeNodeID) []widget.TreeNodeID {
	return p.Keys(id)
}

func (p *PFyneTreePath) IsBranch(id widget.TreeNodeID) bool {
	isDir, ok := p.IsDir(id)
	return isDir && ok
}

func (p *PFyneTreePath) Create(branch bool) fyne.CanvasObject {
	guiIco := widget.NewIcon(theme.MoreHorizontalIcon())
	guiLab := widget.NewLabel("")
	return container.NewHBox(guiIco, guiLab)
}

func (p *PFyneTreePath) Update(id widget.TreeNodeID, branch bool, o fyne.CanvasObject) {
	if p.IsBranch(id) {
		o.(*fyne.Container).Objects[0].(*widget.Icon).SetResource(theme.FolderIcon())
	} else {
		o.(*fyne.Container).Objects[0].(*widget.Icon).SetResource(theme.FileIcon())
	}
	o.(*fyne.Container).Objects[1].(*widget.Label).SetText(path.Base(id))
}
