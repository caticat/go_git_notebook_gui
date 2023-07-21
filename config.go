package main

import (
	"encoding/json"

	"github.com/caticat/go_game_server/phelp"
	"github.com/caticat/go_game_server/plog"
)

type PNoteBookCfg struct {
	Repository  string `json:"repository"`  // 仓库地址
	Local       string `json:"local"`       // 本地地址
	Username    string `json:"username"`    // 仓库账号,可以为空
	Password    string `json:"password"`    // 仓库密码/仓库Token
	AuthorName  string `json:"authorname"`  // 提交用户名
	AuthorEMail string `json:"authoremail"` // 提交用户邮箱
	LogLevel    int    `json:"loglevel"`    // 日志等级
}

func newPNoteBookCfg() *PNoteBookCfg {
	return &PNoteBookCfg{}
}

func (p *PNoteBookCfg) init() error {
	app := getApp()
	if app == nil {
		return ErrAppNotInit
	}

	// 配置读取
	strCfg := app.Preferences().StringWithFallback(APP_CFG_KEY, APP_CFG_VALUE_DEFAULT)
	if err := json.Unmarshal([]byte(strCfg), p); err != nil {
		return err
	}

	// 设置文件数据
	files := NewPFyneTreePath(p.Local)
	if err := files.Refresh(); err != nil {
		plog.ErrorLn(err)
		// return err
	}
	setFiles(files)

	// git同步
	if err := sync(); err != nil {
		plog.ErrorLn(err)
	}

	return nil
}

func (p *PNoteBookCfg) save() error {
	app := getApp()
	if app == nil {
		return ErrAppNotInit
	}

	strCfg := phelp.ToJsonIndent(p)
	app.Preferences().SetString(APP_CFG_KEY, strCfg)

	return nil
}

func (p *PNoteBookCfg) getLogLevel() int { return p.LogLevel }
func (p *PNoteBookCfg) setLogLevel(logLevel int) {
	p.LogLevel = logLevel
}
