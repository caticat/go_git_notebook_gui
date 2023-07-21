package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"github.com/caticat/go_game_server/pfyne_theme_cn"
)

func initGUI() {
	// 窗口初始化
	app := getApp()
	app.Settings().SetTheme(pfyne_theme_cn.NewThemeCN())
	win := app.NewWindow(WINDOW_TITLE)
	win.SetMaster()
	setWin(win)

	// 界面元素声明
	var (
		guiTabMain *container.AppTabs
	)

	// 界面页签
	guiHomeTab := initGUIHome("Home", theme.HomeIcon())
	guiLog := initGUILog()
	guiSetting := initGUISetting()
	guiInfo := initGUIInfo()
	guiTabMain = container.NewAppTabs(
		guiHomeTab,
		container.NewTabItemWithIcon("Log", theme.DocumentIcon(), guiLog),
		container.NewTabItemWithIcon("Setting", theme.SettingsIcon(), guiSetting),
		container.NewTabItemWithIcon("Info", theme.InfoIcon(), guiInfo),
	)
	// guiTabMain.OnSelected = func(ti *container.TabItem) {
	// 	plog.InfoLn(ti.Text)
	// 	// iconName := ti.Icon.Name()
	// 	// if iconName == theme.InfoIcon().Name() { // 更新Info界面的数据
	// 	// 	infoData.Set(phelp.ToJsonIndent(getConf()))
	// 	// }
	// 	// if iconName == theme.DocumentIcon().Name() { // 日志界面隐藏逻辑,减少界面刷新
	// 	// 	guiLog.Show()
	// 	// } else {
	// 	// 	guiLog.Hide()
	// 	// }
	// }
	guiTabMain.SetTabLocation(container.TabLocationLeading)

	// 窗口尺寸
	win.SetContent(guiTabMain)
	win.Resize(fyne.NewSize(GUI_WINDOW_INIT_SIZE_W, GUI_WINDOW_INIT_SIZE_H))
}

func runGUI() {
	win := getWin()
	win.ShowAndRun()
}
