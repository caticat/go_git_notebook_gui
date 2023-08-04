package main

import (
	"fmt"
	"os"
	"path"
	"strings"
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
func onEditorChange(binEntry binding.String, guiButSave *widget.Button, guiPreview *widget.RichText) {
	contentNew, err := binEntry.Get()
	if err != nil {
		plog.ErrorLn(err)
		return
	}

	if guiPreview.Visible() {
		if phelp.IsImage([]byte(contentNew)) {
			plog.DebugLn("![image](" + getOpenFileName() + ")")
			guiPreview.ParseMarkdown("![image](" + getOpenFileName() + ")")
		} else {
			// 相对路径转化为绝对路径
			pathReplace := getCfg().Local
			regPathRelative := getRegPathRelative()
			sliMatch := regPathRelative.FindAllStringSubmatch(contentNew, -1)
			for _, ls := range sliMatch {
				if len(ls) != GUI_HOME_REG_PATHRELATIVE_NUM {
					continue
				}
				if ls[GUI_HOME_REG_PATHRELATIVE_BASE] == GUI_HOME_REG_PATHRELATIVE_EXCEPTION {
					continue
				} else {
					pathNew := ls[GUI_HOME_REG_PATHRELATIVE_PREFIX] + path.Join(pathReplace, ls[GUI_HOME_REG_PATHRELATIVE_BASE])
					contentNew = strings.ReplaceAll(contentNew, ls[GUI_HOME_REG_PATHRELATIVE_ALL], pathNew)
				}
			}

			// 解析markdown
			guiPreview.ParseMarkdown(contentNew)
		}
	}

	contentOri := getOpenFileContent()
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
		plog.InfoF("add folder %q by %s\n", fileName, c.AuthorName)
		// 创建文件夹无需提交仓库,空文件夹也提交不上去
		// return g.Commit(fmt.Sprintf("add folder %q by %s", fileName, c.AuthorName))
		return nil
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
		if err := g.Commit(fmt.Sprintf("add all to notebook by %s", c.AuthorName)); err != nil {
			return err
		}
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

	// 出错函数处理
	funOnError := func(err error) {
		plog.ErrorLn(err)
		guiDiaProgress.Hide()
		dialog.NewError(err, getWin()).Show()
	}

	// 本地提交
	guiDiaProgress.SetDismissText("commit local")
	c := getCfg()
	if c == nil {
		funOnError(ErrConfigNotFound)
		return ErrConfigNotFound
	}
	g := getPGit()
	if g == nil {
		funOnError(ErrGitNotSync)
		return ErrGitNotSync
	}
	if err := funOperation(); err != nil {
		funOnError(err)
		return err
	}

	// 推送
	guiDiaProgress.SetDismissText("syncing")

	// 同步数据
	if err := sync(); err != nil {
		funOnError(err)
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

func gitLogs(logNum int) string {
	g := getPGit()
	if g == nil {
		plog.ErrorLn(ErrGitNotSync)
		return ""
	}

	sliLog, err := g.Log(logNum)
	if err != nil {
		plog.ErrorLn(err)
		return ""
	}

	ret := ""
	for _, log := range sliLog {
		ret += fmt.Sprintf("[%s][%s]%s\n", log.Author.When.Format(GUI_LOG_GIT_TIME_FORMAT), log.Author.String(), strings.TrimSpace(log.Message))
	}

	return ret
}

func changeConfig(needRmLocal bool) error {
	// 同步流程
	conf := getCfg()
	return fileOperationProgress("change config", conf.Repository, false, func() error {
		setPGit(nil) // 清空git信息

		if needRmLocal {
			if err := phelp.Rm(conf.Local); err != nil {
				plog.ErrorLn(err) // 这里不停止后续处理,因为`getCfg().Local`这个文件夹一定被占用无法删除,但是内部的文件会被删除干净
			}
		}

		if err := sync(); err != nil {
			return err
		}
		return nil
	})
}
