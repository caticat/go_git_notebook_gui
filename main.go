package main

import (
	"fyne.io/fyne/v2/app"
	"github.com/caticat/go_game_server/plog"
)

func main() {
	// 初始化日志
	plog.Init(plog.ELogLevel_Debug, "")

	// 初始化App
	setApp(app.NewWithID(APP_ID))

	// 初始化配置
	conf := getCfg()
	if err := conf.init(); err != nil {
		plog.FatalLn(err)
	}

	// 界面初始化
	initGUI()

	// 日志额外设置
	plog.SetOutput(NewPLogWriter()) // 配置读取后才将日志输出只向GUI界面
	plog.SetShortFile()             // GUI中采用短文件名记录日志

	// 运行
	runGUI()

	// 关闭
	close()
}

func close() {
	app := getApp()
	if !app.Driver().Device().IsMobile() {
		app.Quit()
	}
}
