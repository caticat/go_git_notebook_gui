package main

import (
	"fmt"
	"io/ioutil"
	"path"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/caticat/go_game_server/phelp"
	"github.com/caticat/go_game_server/plog"
)

func initGUIHome(text string, icon fyne.Resource) *container.TabItem {
	var (
		guiSearch            *widget.SelectEntry
		guiButSearchToggle   *widget.Button
		guiButRefresh        *widget.Button
		guiButSave           *widget.Button
		guiButAdd            *widget.Button
		guiButDel            *widget.Button
		guiPath              *widget.Tree
		guiEditorContent     *widget.Entry
		guiEditorPathContent *widget.Entry
		guiBodyContent       *fyne.Container
		guiBodyPathContent   *container.Split
		binEditor            = binding.NewString()
	)

	funOnEditorChange := func(s string) {
		if err := binEditor.Set(s); err != nil {
			plog.ErrorLn(err)
		}
		onEditorChange(binEditor, guiButSave)
	}
	funOnRefresh := func() {
		if err := sync(); err != nil {
			plog.Error(err)
			return
		}
		files := getFiles()
		if err := files.Refresh(); err != nil {
			plog.Error(err)
			return
		}

		guiPath.Refresh()
		openFileName := getOpenFileName()
		if files.Has(openFileName) {
			guiPath.UnselectAll()
			guiPath.Select(openFileName)
			guiSearch.SetText(openFileName)
		} else {
			guiPath.Unselect(openFileName)
			guiSearch.SetText("")
			if err := binEditor.Set(""); err != nil {
				plog.Error(err)
				return
			}
		}
	}
	setFunRefresh(funOnRefresh)
	initGUISearch(&guiSearch, &guiPath)
	initGUIButSearchToggle(&guiButSearchToggle, &guiBodyContent, &guiBodyPathContent)
	initGUIButRefresh(&guiButRefresh, funOnRefresh)
	initGUIButSave(&guiButSave, binEditor, &guiPath)
	initGUIButAdd(&guiButAdd, &guiPath, funOnRefresh)
	initGUIButDel(&guiButDel, &guiPath, funOnRefresh)
	guiHead := container.NewGridWithColumns(2, guiSearch, container.NewHBox(guiButSearchToggle, guiButRefresh, guiButSave, guiButAdd, guiButDel))

	initGUIPath(&guiPath, binEditor, &guiEditorContent, &guiEditorPathContent)
	initGuiBodyContent(&guiEditorContent, &guiBodyContent, binEditor, funOnEditorChange)
	initGuiBodyPathContent(&guiEditorPathContent, &guiBodyPathContent, binEditor, funOnEditorChange, &guiPath)
	guiBody := container.NewMax(guiBodyContent, guiBodyPathContent)

	guiMain := container.NewBorder(guiHead, nil, nil, nil, guiBody)

	guiHome := container.NewTabItemWithIcon(text, icon, guiMain)

	return guiHome
}

func initGUISearch(pGuiSearch **widget.SelectEntry, pGuiPath **widget.Tree) {
	files := getFiles()

	(*pGuiSearch) = widget.NewSelectEntry(files.KeysAll())

	guiSearch := (*pGuiSearch)
	guiSearch.OnChanged = func(s string) {
		if !files.Has(s) {
			return
		}
		guiSearch.OnSubmitted(s)
	}
	guiSearch.OnSubmitted = func(s string) {
		if !files.Has(s) {
			return
		}
		(*pGuiPath).Select(s)
	}
}

func initGUIButSearchToggle(pGuiButSearchToggle **widget.Button, pGuiBody1 **fyne.Container, pGuiBody2 **container.Split) {
	*pGuiButSearchToggle = widget.NewButtonWithIcon("", theme.ListIcon(), func() {
		guiBody1 := *pGuiBody1
		guiBody2 := *pGuiBody2
		if guiBody1.Visible() {
			guiBody1.Hide()
			guiBody2.Show()
		} else {
			guiBody1.Show()
			guiBody2.Hide()
		}
	})
}

func initGUIButRefresh(pGuiButRefresh **widget.Button, funRefresh func()) {
	*pGuiButRefresh = widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), funRefresh)
}

func initGUIButSave(pGuiButSave **widget.Button, binEditor binding.String, pGuiPath **widget.Tree) {
	*pGuiButSave = widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), func() {
		contentOri := getOpenFileContent()
		contentNew, err := binEditor.Get()
		if err != nil {
			plog.ErrorLn(err)
			return
		}
		if contentOri == contentNew {
			return
		}
		fileName := getOpenFileName()
		if fileName == "" {
			plog.ErrorLn(ErrSaveFileNameEmpty)
			return
		}

		err = saveOpenFile(fileName, contentNew)
		if err != nil {
			plog.ErrorLn(err)
			return
		}
		(*pGuiButSave).Disable()
	})

	guiButSave := *pGuiButSave
	guiButSave.Disable()
}

func initGUIButAdd(pGuiButAdd **widget.Button, pGuiPath **widget.Tree, funRefresh func()) {
	fileType := ""
	win := getWin()
	binName := binding.NewString()
	guiRadio := widget.NewRadioGroup([]string{"file", "folder"}, func(s string) { fileType = s })
	guiRadio.SetSelected("file")
	*pGuiButAdd = widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		guiDiaAdd := dialog.NewForm("New", "OK", "Cancel", []*widget.FormItem{
			widget.NewFormItem("type", guiRadio),
			widget.NewFormItem("name", widget.NewEntryWithData(binName)),
		},
			func(b bool) {
				if !b {
					return
				}
				fileName, err := binName.Get()
				if err != nil {
					plog.Error(err)
					return
				}
				if fileName == phelp.STR_EMPTY {
					return
				}

				files := getFiles()
				basePath := getPathSelect()
				basePath = files.FixPath(basePath)
				if isDir, ok := files.IsDir(basePath); !(isDir && ok) {
					basePath = path.Dir(basePath)
				}
				filePath := path.Join(basePath, fileName)
				if files.Has(filePath) {
					plog.Error(ErrCreateFileFolderAlreadyExist)
					return
				}

				sync()
				if fileType == "file" {
					err = createFile(filePath)
					if err != nil {
						return
					}
				} else {
					createFolder(filePath)
					if err != nil {
						return
					}
				}
				funRefresh()
				(*pGuiPath).Select(filePath)
			}, win)
		guiDiaAdd.Resize(win.Canvas().Size())
		guiDiaAdd.Show()
	})
}

func initGUIButDel(pGuiButDel **widget.Button, pGuiPath **widget.Tree, funRefresh func()) {
	win := getWin()
	*pGuiButDel = widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		fileNameSelect := getPathSelect()
		files := getFiles()
		fileName := files.FixPath(fileNameSelect)
		conf := getCfg()
		if fileName == path.Join(conf.Local, "") {
			dialog.NewError(ErrDeleteRootDataDirInHome, win).Show()
			plog.ErrorLn(ErrDeleteRootDataDirInHome)
			return
		}
		if !files.Has(fileName) {
			dialog.NewError(ErrDeleteFileFolderNotFound, win).Show()
			plog.ErrorLn(ErrDeleteFileFolderNotFound)
			return
		}
		guiDiaDel := dialog.NewConfirm("Delete", fmt.Sprintf("rm -rf %q?", fileName), func(b bool) {
			if !b {
				return
			}
			if err := deleteFileFolder(fileName); err != nil {
				plog.ErrorLn(err)
				return
			}
			(*pGuiPath).UnselectAll()
			funRefresh()
		}, win)
		// guiDiaDel.Resize(win.Canvas().Size())
		guiDiaDel.Show()
	})
}

func initGUIPath(pGuiPath **widget.Tree, binEditor binding.String, pGuiEditorContent **widget.Entry, pGuiEditorPathContent **widget.Entry) {
	files := getFiles()

	(*pGuiPath) = widget.NewTree(
		files.ChildUIDs,
		files.IsBranch,
		files.Create,
		files.Update,
	)

	guiPath := (*pGuiPath)
	guiPath.OnSelected = func(uid widget.TreeNodeID) {
		setPathSelect(uid)
		if isDir, ok := files.IsDir(uid); isDir && ok {
			return
		} else {
			sliByte, err := ioutil.ReadFile(uid)
			if err != nil {
				plog.ErrorLn(err)
				dialog.NewError(err, getWin()).Show()
				return
			}
			content := string(sliByte)

			binEditor.Set(content)
			setOpenFileName(uid)
			setOpenFileContent(content)

			if (*pGuiEditorContent).Disabled() {
				(*pGuiEditorContent).Enable()
			}
			if (*pGuiEditorPathContent).Disabled() {
				(*pGuiEditorPathContent).Enable()
			}
		}

		for p := uid; (p != phelp.STR_EMPTY) && (p != phelp.PATH_CURRENT); p = path.Dir(p) {
			guiPath.OpenBranch(p)
		}
	}
}

func initGuiBodyContent(pGuiEditorContent **widget.Entry, pGuiBodyContent **fyne.Container, binEditor binding.String, funOnEditorChange func(string)) {
	*pGuiEditorContent = widget.NewMultiLineEntry()
	guiEditorContent := *pGuiEditorContent
	guiEditorContent.Bind(binEditor)
	guiEditorContent.OnChanged = funOnEditorChange
	guiEditorContent.Disable()
	*pGuiBodyContent = container.NewMax(guiEditorContent)
	guiBodyContent := *pGuiBodyContent
	guiBodyContent.Hide()
}

func initGuiBodyPathContent(pGuiEditorPathContent **widget.Entry, pGuiBodyPathContent **container.Split, binEditor binding.String, funOnEditorChange func(string), pGuiPath **widget.Tree) {
	*pGuiEditorPathContent = widget.NewMultiLineEntry()
	guiEditorPathContent := *pGuiEditorPathContent
	guiEditorPathContent.Bind(binEditor)
	guiEditorPathContent.OnChanged = funOnEditorChange
	guiEditorPathContent.Disable()
	*pGuiBodyPathContent = container.NewHSplit(*pGuiPath, guiEditorPathContent)
	guiBodyPathContent := *pGuiBodyPathContent
	guiBodyPathContent.SetOffset(GUI_HOME_OFFSET)
}
