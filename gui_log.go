package main

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/caticat/go_game_server/phelp"
	"github.com/caticat/go_game_server/plog"
)

func initGUILog() fyne.CanvasObject {
	guiTab := container.NewAppTabs(
		initGUILogNormal(),
		initGUILogGit(),
	)

	return guiTab
}

func initGUILogNormal() *container.TabItem {
	// 上侧日志等级选择与功能按钮
	guiSelLogLevel := widget.NewSelect([]string{
		string(plog.SLogLevel_Debug),
		string(plog.SLogLevel_Info),
		string(plog.SLogLevel_Warn)}, func(s string) {
		logLevel := plog.ToLogLevel(s)
		getCfg().setLogLevel(int(logLevel))
		plog.InfoLn("change log level to:", s)
	})
	guiSelLogLevel.SetSelected(plog.ToLogLevelName(plog.ELogLevel(getCfg().LogLevel)))
	guiButClearLog := widget.NewButtonWithIcon(phelp.STR_EMPTY, theme.DeleteIcon(), func() {
		getLogData().Set(phelp.STR_EMPTY)
	})
	guiConLogLevel := container.NewAdaptiveGrid(2, guiSelLogLevel, container.NewHBox(guiButClearLog))
	guiForLogLevel := widget.NewForm(widget.NewFormItem("LogLevel", guiConLogLevel))

	// 日志内容
	guiScrLog := container.NewScroll(widget.NewLabelWithData(getLogData()))

	// 界面组合
	return container.NewTabItemWithIcon("App", theme.DocumentIcon(), container.NewBorder(guiForLogLevel, nil, nil, nil, container.NewMax(guiScrLog)))
}

func initGUILogGit() *container.TabItem {
	// 上侧 日志数量
	binLogNum := binding.NewString()
	guiLogNum := widget.NewEntryWithData(binLogNum)
	guiLogNum.SetPlaceHolder(GUI_LOG_GIT_PLACEHOLDER)
	guiLogNum.SetText(GUI_LOG_GIT_DEFAULT)
	guiLogNum.Validator = func(s string) error {
		if getRegNum().MatchString(s) {
			return nil
		} else {
			return ErrGitLogNumNeedToBeInt
		}
	}

	// 上侧 查询按钮
	binLog := binding.NewString()
	guiButLog := widget.NewButtonWithIcon("", theme.SearchIcon(), func() {
		strLogNum, err := binLogNum.Get()
		if err != nil {
			dialog.NewError(err, getWin()).Show()
			return
		}
		logNum, err := strconv.Atoi(strLogNum)
		if err != nil {
			dialog.NewError(err, getWin()).Show()
			return
		}
		binLog.Set(gitLogs(logNum))
	})

	// 上侧 总
	guiHead := container.NewAdaptiveGrid(2, guiLogNum, container.NewHBox(guiButLog))

	// 下侧 日志内容
	guiScrLog := container.NewScroll(widget.NewLabelWithData(binLog))

	// 界面组合
	return container.NewTabItemWithIcon("Git", theme.ListIcon(), container.NewBorder(guiHead, nil, nil, nil, guiScrLog))
}
