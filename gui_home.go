package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/caticat/go_game_server/phelp"
	"github.com/caticat/go_game_server/phelp/ppath"
	"github.com/caticat/go_game_server/plog"
)

func initGUIHome(text string, icon fyne.Resource) *container.TabItem {
	var (
		funBodyShow        func()
		guiSearch          *widget.SelectEntry
		guiButToggleLeft   *widget.Button
		guiButToggleMiddle *widget.Button
		guiButToggleRight  *widget.Button
		guiButRefresh      *widget.Button
		guiButSave         *widget.Button
		guiButAdd          *widget.Button
		guiButDel          *widget.Button
		guiButMove         *widget.Button
		guiPath            *widget.Tree
		guiEditor          *fyne.Container
		guiEditorContent   *PCSEntry
		guiImage           *canvas.Image
		guiPreview         *widget.RichText
		guiBody            *fyne.Container
		binEditor          = binding.NewString()
		guiLabLogLast      *widget.Label // 身体 日志最后一行
	)

	funOnEditorChange := func(s string) {
		if err := binEditor.Set(s); err != nil {
			plog.ErrorLn(err)
		}
		onEditorChange(binEditor, guiButSave, guiPreview)
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
		plog.InfoLn("refresh done")
	}
	setFunRefresh(funOnRefresh)
	initGUISearch(&guiSearch, &guiPath)
	initGUIButToggle(&funBodyShow, &guiButToggleLeft, &guiButToggleMiddle, &guiButToggleRight, &guiBody, &guiPath, &guiEditor, &guiPreview)
	initGUIButRefresh(&guiButRefresh, funOnRefresh)
	initGUIButSave(&guiButSave, binEditor, &guiPath)
	initGUIButAdd(&guiButAdd, &guiPath, funOnRefresh)
	initGUIButDel(&guiButDel, &guiPath, funOnRefresh)
	initGUIButMove(&guiButMove, &guiPath, funOnRefresh)
	guiHead := container.NewGridWithColumns(2, guiSearch, container.NewHBox(guiButToggleLeft, guiButToggleMiddle, guiButToggleRight, guiButRefresh, guiButSave, guiButAdd, guiButDel, guiButMove))

	initGUIPath(&guiPath, binEditor, &guiEditorContent, &guiImage)
	initGuiBodyContent(&guiEditor, &guiEditorContent, &guiImage, &guiPreview, binEditor, funOnEditorChange, &guiButSave)
	guiBody = container.NewMax()
	funBodyShow()

	guiLabLogLast = widget.NewLabelWithData(getLogLast())
	guiMain := container.NewBorder(guiHead, container.NewHScroll(guiLabLogLast), nil, nil, guiBody)

	guiHome := container.NewTabItemWithIcon(text, icon, guiMain)

	return guiHome
}

func initGUISearch(pGuiSearch **widget.SelectEntry, pGuiPath **widget.Tree) {
	files := getFiles()

	(*pGuiSearch) = widget.NewSelectEntry(files.KeysAll(""))

	guiSearch := (*pGuiSearch)
	guiSearch.SetPlaceHolder(GUI_HOME_SEARCH_PLACEHOLDER)
	guiSearch.OnChanged = func(s string) {
		if (s == "") || (!files.Has(s)) {
			guiSearch.SetOptions(files.KeysAll(s))
			// guiSearch.ActionItem.(*widget.Button).OnTapped() // 打开选择框,会丢失输入焦点
			// getWin().Canvas().Focus(&guiSearch.Entry) // 强制选择输入框做焦点,无效...没有找对canvas
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

func initGUIButToggle(pFunBodyShow *func(),
	pGuiButToggleLeft **widget.Button,
	pGuiButToggleMiddle **widget.Button,
	pGuiButToggleRight **widget.Button,
	pGuiBody **fyne.Container,
	pGuiPath **widget.Tree,
	pGuiEditor **fyne.Container,
	pGuiPreview **widget.RichText) {
	conf := getCfg()
	bodyshow := conf.getHomeLayout()
	*pGuiButToggleLeft = widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
		if (bodyshow & GUI_HOME_TOGGLE_LEFT) > 0 {
			bodyshow = bodyshow & ^GUI_HOME_TOGGLE_LEFT
		} else {
			bodyshow = bodyshow | GUI_HOME_TOGGLE_LEFT
		}
		(*pFunBodyShow)()
	})
	*pGuiButToggleMiddle = widget.NewButtonWithIcon("", theme.MoveUpIcon(), func() {
		if (bodyshow & GUI_HOME_TOGGLE_MIDDLE) > 0 {
			bodyshow = bodyshow & ^GUI_HOME_TOGGLE_MIDDLE
		} else {
			bodyshow = bodyshow | GUI_HOME_TOGGLE_MIDDLE
		}
		(*pFunBodyShow)()
	})
	*pGuiButToggleRight = widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		if (bodyshow & GUI_HOME_TOGGLE_RIGHT) > 0 {
			bodyshow = bodyshow & ^GUI_HOME_TOGGLE_RIGHT
		} else {
			bodyshow = bodyshow | GUI_HOME_TOGGLE_RIGHT
		}
		(*pFunBodyShow)()
	})

	(*pFunBodyShow) = func() {
		if bodyshow == 0 {
			bodyshow = GUI_HOME_TOGGLE_MIDDLE
		}
		body := *pGuiBody
		body.RemoveAll()
		sliWid := make([]fyne.CanvasObject, 0)
		if bodyshow&GUI_HOME_TOGGLE_MIDDLE > 0 {
			sliWid = append(sliWid, container.NewScroll(*pGuiEditor))
			(*pGuiButToggleMiddle).SetIcon(theme.MoveUpIcon())
		} else {
			(*pGuiButToggleMiddle).SetIcon(theme.MoveDownIcon())
		}
		if bodyshow&GUI_HOME_TOGGLE_RIGHT > 0 {
			sliWid = append(sliWid, container.NewScroll(*pGuiPreview))
			(*pGuiButToggleRight).SetIcon(theme.NavigateNextIcon())
		} else {
			(*pGuiButToggleRight).SetIcon(theme.NavigateBackIcon())
		}
		lenWid := len(sliWid)
		if bodyshow&GUI_HOME_TOGGLE_LEFT > 0 {
			if lenWid > 0 {
				guiSplit := container.NewHSplit(*pGuiPath, container.NewAdaptiveGrid(lenWid, sliWid...))
				guiSplit.SetOffset(GUI_HOME_OFFSET)
				body.Add(guiSplit)
			} else {
				body.Add(*pGuiPath)
			}
			(*pGuiButToggleLeft).SetIcon(theme.NavigateBackIcon())
		} else {
			body.Add(container.NewAdaptiveGrid(lenWid, sliWid...))
			(*pGuiButToggleLeft).SetIcon(theme.NavigateNextIcon())
		}

		conf.setHomeLayout(bodyshow)
		conf.save()
	}
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

		plog.InfoF("save %q done\n", fileName)
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

				plog.InfoF("add %s %s done\n", fileType, filePath)
			}, win)
		canSize := win.Canvas().Size()
		guiDiaAdd.Resize(fyne.NewSize(canSize.Width/2, 0)) // 总宽度的一般,最小高度
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
			plog.InfoLn("del " + fileName + " done")
		}, win)
		// guiDiaDel.Resize(win.Canvas().Size())
		guiDiaDel.Show()
	})
}

func initGUIButMove(pGuiButMove **widget.Button, pGuiPath **widget.Tree, funRefresh func()) {
	win := getWin()
	*pGuiButMove = widget.NewButtonWithIcon("", theme.ContentCutIcon(), func() {
		fileNameSelect := getPathSelect()
		files := getFiles()
		fileName := files.FixPath(fileNameSelect)
		conf := getCfg()
		if fileName == path.Join(conf.Local, "") {
			dialog.NewError(ErrMoveRootDataDirInHome, win).Show()
			plog.ErrorLn(ErrMoveRootDataDirInHome)
			return
		}
		if !files.Has(fileName) {
			dialog.NewError(ErrMoveFileFolderNotFound, win).Show()
			plog.ErrorLn(ErrMoveFileFolderNotFound)
			return
		}
		fileNameMove, err := filepath.Rel(conf.Local, fileName)
		if err != nil {
			plog.ErrorLn(err)
			return
		}

		binMoveTo := binding.NewString()
		guiDiaMov := dialog.NewCustomConfirm(fmt.Sprintf("Move %q to", fileNameMove), "OK", "Cancel", widget.NewEntryWithData(binMoveTo), func(b bool) {
			if !b {
				return
			}
			pathTo, err := binMoveTo.Get()
			if err != nil {
				plog.ErrorLn(err)
				return
			}
			conf := getCfg()
			pathTo = path.Join(conf.Local, pathTo)
			// if files.Has(pathTo) { // 重名文件校验处理
			// 	plog.ErrorLn()
			// 	return
			// }
			pathTo, err = filepath.Abs(pathTo)
			if err != nil {
				plog.ErrorLn(err)
				return
			}
			pathFrom, err := filepath.Abs(getPathSelect())
			if err != nil {
				plog.ErrorLn(err)
				return
			}
			pathBase, err := filepath.Abs(conf.Local)
			if err != nil {
				plog.ErrorLn(err)
				return
			}
			if !ppath.IsSubDir(pathBase, pathTo) {
				plog.ErrorLn(ErrMoveFilePathOutOfData)
				return
			}
			if err := move(pathFrom, pathTo); err != nil {
				plog.ErrorLn(err)
				return
			}
			(*pGuiPath).UnselectAll()
			funRefresh()

			plog.InfoF("move %q %q done\n", pathFrom, pathTo)
		}, win)
		guiDiaMov.Resize(fyne.NewSize(win.Canvas().Size().Width/2, 0))
		guiDiaMov.Show()
	})
}

func initGUIPath(pGuiPath **widget.Tree, binEditor binding.String, pGuiEditorContent **PCSEntry, pGuiImage **canvas.Image) {
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
			sliByte, err := os.ReadFile(uid)
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

			// 编辑界面展示内容变更 文字/图片
			guiImage := *pGuiImage
			if phelp.IsImage(sliByte) {
				guiImage := *pGuiImage
				guiImage.File = uid
				if !guiImage.Visible() {
					guiImage.Show()
				}
				if (*pGuiEditorContent).Visible() {
					(*pGuiEditorContent).Hide()
				}
			} else {
				if guiImage.Visible() {
					guiImage.Hide()
				}
				if !(*pGuiEditorContent).Visible() {
					(*pGuiEditorContent).Show()
				}
			}
		}
		loopLimit := GUI_HOME_PATH_LOOP_MAX
		for p := uid; strings.Contains(p, APP_CFG_PATH_LOCAL) && loopLimit > 0; p = path.Dir(p) {
			guiPath.OpenBranch(p)
			loopLimit--
		}
		plog.InfoF("open file %q done\n", uid)
	}
}

func initGuiBodyContent(pGuiEditor **fyne.Container, pGuiEditorContent **PCSEntry, pGuiImage **canvas.Image, pGuiPreview **widget.RichText, binEditor binding.String, funOnEditorChange func(string), pGuiButSave **widget.Button) {

	*pGuiEditorContent = newPCSEntry(func() {
		(*pGuiButSave).OnTapped()
	})
	guiEditorContent := *pGuiEditorContent
	guiEditorContent.Bind(binEditor)
	guiEditorContent.OnChanged = funOnEditorChange
	guiEditorContent.Disable()
	// guiEditorContent.TypedShortcut(&fyne)

	*pGuiImage = &canvas.Image{}
	guiImage := *pGuiImage
	guiImage.FillMode = canvas.ImageFillOriginal
	guiImage.Hide()

	*pGuiPreview = widget.NewRichText()

	*pGuiEditor = container.NewMax(guiEditorContent, guiImage)
}

// func initGuiBodyPathContent(pGuiEditorPathContent **widget.Entry, pGuiBodyPathContent **container.Split, binEditor binding.String, funOnEditorChange func(string), pGuiPath **widget.Tree) {
// 	*pGuiEditorPathContent = widget.NewMultiLineEntry()
// 	guiEditorPathContent := *pGuiEditorPathContent
// 	guiEditorPathContent.Bind(binEditor)
// 	guiEditorPathContent.OnChanged = funOnEditorChange
// 	guiEditorPathContent.Disable()
// 	*pGuiBodyPathContent = container.NewHSplit(*pGuiPath, guiEditorPathContent)
// 	guiBodyPathContent := *pGuiBodyPathContent
// 	guiBodyPathContent.SetOffset(GUI_HOME_OFFSET)
// }
