package main

import (
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
	sliFormItemData := []*struct {
		K, V       string
		IsPassword bool
		BinData    binding.String
		GUI        *widget.FormItem
	}{
		{K: "repository", V: conf.Repository},
		{K: "username", V: conf.Username},
		{K: "password", V: conf.Password},
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
		var err error

		index := 0
		conf.Repository, err = sliFormItemData[index].BinData.Get()
		if err != nil {
			plog.ErrorLn(err)
		}

		index++
		conf.Username, err = sliFormItemData[index].BinData.Get()
		if err != nil {
			plog.ErrorLn(err)
		}

		index++
		conf.Password, err = sliFormItemData[index].BinData.Get()
		if err != nil {
			plog.ErrorLn(err)
		}

		index++
		conf.AuthorName, err = sliFormItemData[index].BinData.Get()
		if err != nil {
			plog.ErrorLn(err)
		}

		index++
		conf.AuthorEMail, err = sliFormItemData[index].BinData.Get()
		if err != nil {
			plog.ErrorLn(err)
		}

		err = conf.save()
		if err != nil {
			plog.ErrorLn(err)
		}

		err = sync()
		if err != nil {
			plog.ErrorLn(err)
		}
	}

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
			// TODO 功能待制作
		}),
		widget.NewButtonWithIcon("Force Push(use local data)", theme.UploadIcon(), func() {
			// TODO 功能待制作
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
			}, win)
			guiDiaDel.Show()
		}),
		widget.NewButtonWithIcon("Export Notebook", theme.LogoutIcon(), func() {
			// TODO 功能待制作
		}),
		widget.NewButtonWithIcon("Import Notebook", theme.LoginIcon(), func() {
			// TODO 功能待制作
		}),
	)
}
