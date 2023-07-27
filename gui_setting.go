package main

import (
	"fmt"
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

func initGUISetting() fyne.CanvasObject {
	return container.NewAppTabs(
		container.NewTabItem("basic", initGUISettingBasic()),
		container.NewTabItem("advance", initGUISettingAdvance()),
	)
}

func initGUISettingBasic() *widget.Form {
	conf := getCfg()
	win := getWin()
	sliFormItemData := []*struct {
		K, V       string
		IsPassword bool
		BinData    binding.String
		GUI        *widget.FormItem
	}{
		{K: "repository", V: conf.Repository},
		{K: "username", V: conf.Username},
		{K: "password", V: conf.Password, IsPassword: true},
		{K: "authorname", V: conf.AuthorName},
		{K: "authoremail", V: conf.AuthorEMail},
	}

	guiForm := widget.NewForm()
	for _, d := range sliFormItemData {
		d.BinData, d.GUI = initGUISettingFormItem(d.K, d.V, d.IsPassword)
		guiForm.AppendItem(d.GUI)
	}

	guiForm.OnCancel = func() {
		index := 0
		if err := sliFormItemData[index].BinData.Set(conf.Repository); err != nil {
			plog.ErrorLn(err)
		}

		index++
		if err := sliFormItemData[index].BinData.Set(conf.Username); err != nil {
			plog.ErrorLn(err)
		}

		index++
		if err := sliFormItemData[index].BinData.Set(conf.Password); err != nil {
			plog.ErrorLn(err)
		}

		index++
		if err := sliFormItemData[index].BinData.Set(conf.AuthorName); err != nil {
			plog.ErrorLn(err)
		}

		index++
		if err := sliFormItemData[index].BinData.Set(conf.AuthorEMail); err != nil {
			plog.ErrorLn(err)
		}
	}
	guiForm.OnSubmit = func() {
		index := 0
		repository, err := sliFormItemData[index].BinData.Get()
		if err != nil {
			dialog.NewError(err, win).Show()
			plog.ErrorLn(err)
			return
		}
		needRmLocal := repository != conf.Repository

		funChangeSetting := func() {
			index++
			conf.Username, err = sliFormItemData[index].BinData.Get()
			if err != nil {
				dialog.NewError(err, win).Show()
				plog.ErrorLn(err)
				return
			}

			index++
			conf.Password, err = sliFormItemData[index].BinData.Get()
			if err != nil {
				dialog.NewError(err, win).Show()
				plog.ErrorLn(err)
				return
			}

			index++
			conf.AuthorName, err = sliFormItemData[index].BinData.Get()
			if err != nil {
				dialog.NewError(err, win).Show()
				plog.ErrorLn(err)
				return
			}

			index++
			conf.AuthorEMail, err = sliFormItemData[index].BinData.Get()
			if err != nil {
				dialog.NewError(err, win).Show()
				plog.ErrorLn(err)
				return
			}
			conf.Repository = repository // 最后才更新仓库地址

			err = conf.save()
			if err != nil {
				dialog.NewError(err, win).Show()
				plog.ErrorLn(err)
				return
			}

			err = changeConfig(needRmLocal)
			if err != nil {
				dialog.NewError(err, win).Show()
				plog.ErrorLn(err)
				return
			}
			getFunRefresh()()

			// dialog.NewInformation("Setting Config", "Update Config Done!", win).Show()
			plog.InfoLn("update config done")
		}

		if needRmLocal {
			guiCon := dialog.NewConfirm("Change Repo", fmt.Sprintf("Change Repo From:%q to %q while delete all local data\nSure?", conf.Repository, repository),
				func(b bool) {
					if b {
						funChangeSetting()
					} else {
						err = sliFormItemData[index].BinData.Set(conf.Repository)
						if err != nil {
							plog.ErrorLn(err)
						}
					}
				}, win)
			guiCon.Resize(fyne.NewSize(win.Canvas().Size().Width/2, 0))
			guiCon.Show()
		} else {
			funChangeSetting()
		}
	}
	guiForm.Refresh()

	return guiForm
}

func initGUISettingFormItem(text string, value string, isPassword bool) (binding.String, *widget.FormItem) {
	binData := binding.BindString(&value)
	guiEnt := widget.NewEntryWithData(binData)
	guiEnt.Password = isPassword
	return binData, widget.NewFormItem(text, guiEnt)
}

func initGUISettingAdvance() fyne.CanvasObject {
	win := getWin()

	return container.NewVBox(
		widget.NewButtonWithIcon("Force Update(drop local modify data)", theme.DownloadIcon(), func() {
			dialog.NewConfirm("Force Update", "drop all local change?", func(b bool) {
				if !b {
					return
				}
				if err := forceUpdate(); err != nil {
					plog.ErrorLn(err)
					return
				}
				getFunRefresh()()
				plog.InfoLn("Force Update Done")
			}, win).Show()
		}),
		widget.NewButtonWithIcon("Force Push(use local data)", theme.UploadIcon(), func() {
			dialog.NewConfirm("Force Push", "force push local data to remote?", func(b bool) {
				if !b {
					return
				}
				if err := forcePush(); err != nil {
					plog.ErrorLn(err)
					return
				}
				getFunRefresh()()
				plog.InfoLn("Force Push Done")
			}, win).Show()
		}),
		widget.NewButtonWithIcon("Delete Data All", theme.DeleteIcon(), func() {
			guiDiaDel := dialog.NewConfirm("Delete", "delete all data(sync to repo)?", func(b bool) {
				if !b {
					return
				}
				c := getCfg()
				files := make(map[string]bool)
				if err := phelp.Ls(c.Local, files, phelp.PBinFlag_None); err != nil {
					plog.ErrorLn(err)
					return
				}
				if err := deleteFiles(files); err != nil {
					plog.ErrorLn(err)
					return
				}
				if err := saveOpenFile(path.Join(c.Local, NOTEBOOK_TMP_FILE_NAME), NOTEBOOK_TMP_FILE_CONTENT); err != nil {
					plog.ErrorLn(err)
					return
				}
				getFunRefresh()()
				plog.InfoLn("Delete Data All Done")
			}, win)
			guiDiaDel.Show()
		}),
		widget.NewButtonWithIcon("Export Notebook", theme.LogoutIcon(), func() {
			diaFolder := dialog.NewFolderOpen(func(lu fyne.ListableURI, err error) {
				if lu == nil {
					return
				}
				fileTo := lu.Path()
				dialog.NewConfirm("Export Notebook", fmt.Sprintf("export notebook to:%q", fileTo), func(b bool) {
					if !b {
						return
					}
					if err := export(fileTo); err != nil {
						plog.ErrorLn(err)
						return
					}
					plog.InfoF("Export all notebook to %q done\n", fileTo)
				}, win).Show()
			}, win)
			diaFolder.Resize(win.Canvas().Size())
			diaFolder.Show()
		}),
		widget.NewButtonWithIcon("Import Notebook", theme.LoginIcon(), func() {
			binFileType := binding.NewString()
			guiRadio := widget.NewRadioGroup([]string{"file", "folder"}, func(s string) { binFileType.Set(s) })
			guiRadio.SetSelected("file")
			dialog.NewCustomConfirm("Import Notebook", "OK", "Cancel", guiRadio, func(b bool) {
				if !b {
					return
				}
				fileType, err := binFileType.Get()
				if err != nil {
					plog.ErrorLn(err)
					return
				}
				if fileType == "file" {
					diaFile := dialog.NewFileOpen(func(uc fyne.URIReadCloser, err error) {
						if uc == nil {
							return
						}
						defer uc.Close()
						fileFrom := uc.URI().Path()
						if err := importFrom(fileFrom); err != nil {
							plog.ErrorLn(err)
							return
						}
						sync()
						getFunRefresh()()
						plog.InfoF("Import notebook from %q done\n", fileFrom)
					}, win)
					diaFile.Resize(win.Canvas().Size())
					diaFile.Show()
				} else {
					diaFolder := dialog.NewFolderOpen(func(lu fyne.ListableURI, err error) {
						if lu == nil {
							return
						}
						fileFrom := lu.Path()
						dialog.NewConfirm("Import Notebook", fmt.Sprintf("import notebook from:%q", fileFrom), func(b bool) {
							if !b {
								return
							}
							if err := importFrom(fileFrom); err != nil {
								plog.ErrorLn(err)
								return
							}
							sync()
							getFunRefresh()()
							plog.InfoF("Import notebook from %q done\n", fileFrom)
						}, win).Show()
					}, win)
					diaFolder.Resize(win.Canvas().Size())
					diaFolder.Show()
				}
			}, win).Show()
		}),
	)
}
