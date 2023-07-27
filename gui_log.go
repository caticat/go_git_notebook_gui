package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/caticat/go_game_server/phelp"
	"github.com/caticat/go_game_server/plog"
)

func initGUILog() fyne.CanvasObject {
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
	return container.NewBorder(guiForLogLevel, nil, nil, nil, container.NewMax(guiScrLog))
}
