package main

import (
	"fmt"
	"os"
	"time"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/caticat/go_game_server/pgit"
	"github.com/caticat/go_game_server/phelp"
	"github.com/caticat/go_game_server/plog"
)

func sync() error {
	g := getPGit()
	if g == nil {
		c := getCfg()
		gn, err := pgit.NewPGit(c.Repository, c.Local, c.Username, c.Password, c.AuthorName, c.AuthorEMail)
		if err != nil {
			return err
		}
		setPGit(gn)
		g = gn
	}

	return g.Sync()
}

// 编辑器内容变更 更新保存按钮状态
func onEditorChange(binEntry binding.String, guiButSave *widget.Button) {
	contentNew, err := binEntry.Get()
	if err != nil {
		plog.ErrorLn(err)
		return
	}

	contentOri := getOpenFileContent()
	// plog.InfoLn("change ori:", contentOri)
	// plog.InfoLn("change new:", contentNew)
	if contentOri == contentNew {
		if !guiButSave.Disabled() {
			guiButSave.Disable()
		}
		return
	}

	if guiButSave.Disabled() {
		guiButSave.Enable()
	}
}

func saveOpenFile(fileName, fileContent string) error {
	// 本地文件保存
	if err := os.WriteFile(fileName, []byte(fileContent), NOTEBOOK_PERM_FILE); err != nil {
		plog.ErrorLn(err)
		return err
	}
	c := getCfg()
	if c == nil {
		plog.ErrorLn(ErrConfigNotFound)
		return ErrConfigNotFound
	}
	g := getPGit()
	if g == nil {
		plog.ErrorLn(ErrGitNotSync)
		return ErrGitNotSync
	}

	// 同步流程
	return fileOperationProgress("save", fileName, true, func() error {
		return g.Commit(fmt.Sprintf("save %q by %s", fileName, c.AuthorName))
	})
}

func createFile(fileName string) error {
	// 本地文件创建
	if err := os.WriteFile(fileName, []byte(""), NOTEBOOK_PERM_FILE); err != nil {
		plog.Error(err)
		return err
	}
	c := getCfg()
	if c == nil {
		plog.ErrorLn(ErrConfigNotFound)
		return ErrConfigNotFound
	}
	g := getPGit()
	if g == nil {
		plog.ErrorLn(ErrGitNotSync)
		return ErrGitNotSync
	}

	// 同步流程
	return fileOperationProgress("add file", fileName, true, func() error {
		return g.Commit(fmt.Sprintf("add file %q by %s", fileName, c.AuthorName))
	})
}

func createFolder(fileName string) error {
	// 本地文件夹创建
	if err := os.MkdirAll(fileName, NOTEBOOK_PERM_FOLDER); err != nil {
		plog.Error(err)
		return err
	}
	c := getCfg()
	if c == nil {
		plog.ErrorLn(ErrConfigNotFound)
		return ErrConfigNotFound
	}
	g := getPGit()
	if g == nil {
		plog.ErrorLn(ErrGitNotSync)
		return ErrGitNotSync
	}

	// 同步流程
	return fileOperationProgress("add folder", fileName, true, func() error {
		return g.Commit(fmt.Sprintf("add folder %q by %s", fileName, c.AuthorName))
	})
}

func deleteFileFolder(fileName string) error {
	return deleteFiles(map[string]bool{
		fileName: true, // 这里不考虑文件或者文件夹,操作都一样
	})
}

func deleteFiles(files map[string]bool) error {
	// 本地删除文件/文件夹
	fileName := ""
	for f := range files {
		if err := phelp.Rm(f); err != nil {
			plog.Error(err)
			return err
		}
		fileName += f + "\n"
	}
	c := getCfg()
	if c == nil {
		plog.ErrorLn(ErrConfigNotFound)
		return ErrConfigNotFound
	}
	g := getPGit()
	if g == nil {
		plog.ErrorLn(ErrGitNotSync)
		return ErrGitNotSync
	}

	// 同步流程
	return fileOperationProgress("delete file/folder", fileName, true, func() error {
		return g.Commit(fmt.Sprintf("delete file/folder %q by %s", fileName, c.AuthorName))
	})
}

func move(pathFrom, pathTo string) error {
	return fileOperationProgress(fmt.Sprintf("move %q %q", pathFrom, pathTo), pathFrom, true, func() error {
		phelp.Mv(pathFrom, pathTo, phelp.PBinFlag_Recursive|phelp.PBinFlag_Force)
		return nil
	})
}

func forceUpdate() error {
	c := getCfg()
	if c == nil {
		plog.ErrorLn(ErrConfigNotFound)
		return ErrConfigNotFound
	}
	g := getPGit()
	if g == nil {
		plog.ErrorLn(ErrGitNotSync)
		return ErrGitNotSync
	}

	// 同步流程
	return fileOperationProgress("force update", c.Local, false, func() error {
		return g.ResetToRemote()
	})
}

func forcePush() error {
	c := getCfg()
	if c == nil {
		plog.ErrorLn(ErrConfigNotFound)
		return ErrConfigNotFound
	}
	g := getPGit()
	if g == nil {
		plog.ErrorLn(ErrGitNotSync)
		return ErrGitNotSync
	}

	// 同步流程
	return fileOperationProgress("force push", c.Repository, false, func() error {
		return g.PushForce()
	})
}

func export(fileTo string) error {
	c := getCfg()
	if c == nil {
		plog.ErrorLn(ErrConfigNotFound)
		return ErrConfigNotFound
	}

	// 同步流程
	return fileOperationProgress("export notebook", fileTo, false, func() error {
		_, err := phelp.Cp(c.Local, fileTo, phelp.PBinFlag_Recursive|phelp.PBinFlag_Force)
		return err
	})
}

func importFrom(fileFrom string) error {
	c := getCfg()
	if c == nil {
		plog.ErrorLn(ErrConfigNotFound)
		return ErrConfigNotFound
	}
	g := getPGit()
	if g == nil {
		plog.ErrorLn(ErrGitNotSync)
		return ErrGitNotSync
	}

	// 同步流程
	return fileOperationProgress("import notebook", fileFrom, false, func() error {
		if _, err := phelp.Cp(fileFrom, c.Local, phelp.PBinFlag_Recursive|phelp.PBinFlag_Force); err != nil {
			return err
		}
		return g.Commit(fmt.Sprintf("import notebook from %q by %s", fileFrom, c.AuthorName))
	})
}

func fileOperationProgress(strOperation, fileName string, autoClose bool, funOperation func() error) error {
	// 弹窗
	guiProgress := widget.NewProgressBarInfinite()
	guiDiaProgress := dialog.NewCustom(fmt.Sprintf("%s %q", strOperation, fileName), strOperation, guiProgress, getWin())
	go guiDiaProgress.Show()

	// 本地提交
	guiDiaProgress.SetDismissText("commit local")
	c := getCfg()
	if c == nil {
		plog.ErrorLn(ErrConfigNotFound)
		return ErrConfigNotFound
	}
	g := getPGit()
	if g == nil {
		plog.ErrorLn(ErrGitNotSync)
		return ErrGitNotSync
	}
	if err := funOperation(); err != nil {
		plog.ErrorLn(err)
		return err
	}

	// 推送
	guiDiaProgress.SetDismissText("syncing")

	// 同步数据
	if err := sync(); err != nil {
		plog.ErrorLn(err)
		return err
	}

	// 标记完成
	guiProgress.Stop()
	guiDiaProgress.SetDismissText("done")
	if autoClose {
		time.Sleep(time.Millisecond * GUI_DIALOG_AUTOCLOSE_WAIT_TIME)
		guiDiaProgress.Hide()
	}

	return nil
}
