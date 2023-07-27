package main

const (
	// app
	APP_ID       = "github.com.caticat.go_git_notebook_gui"
	WINDOW_TITLE = "PNoteBook"
	APP_AUTHOR   = "Pan J"
	APP_VER      = "v0.0.4"
	APP_CFG_KEY  = "config"

	// config
	APP_CFG_VALUE_DEFAULT = `
{
	"repository" : "",
	"username" : "",
	"password" : "",
	"authorname" : "",
	"authoremail" : "",
	"loglevel" : 1,
	"homelayout" : 7
}
`
	APP_CFG_PATH_LOCAL = "local.repo"

	// gui
	GUI_WINDOW_INIT_SIZE_W         = 1000                  // 窗口初始大小W
	GUI_WINDOW_INIT_SIZE_H         = 700                   // 窗口初始大小H
	GUI_HOME_OFFSET                = 0.3                   // 初始主界面分隔比率
	GUI_DIALOG_AUTOCLOSE_WAIT_TIME = 300                   // 弹窗自动关闭事件 毫秒
	GUI_HOME_SEARCH_PLACEHOLDER    = "Search Pattern ..."  // 主界面搜索提示文字
	GUI_HOME_TOGGLE_LEFT           = 1 << 0                // 主界面显示标记 目录树
	GUI_HOME_TOGGLE_MIDDLE         = 1 << 1                // 主界面显示标记 编辑
	GUI_HOME_TOGGLE_RIGHT          = 1 << 2                // 主界面显示标记 预览
	GUI_HOME_PATH_LOOP_MAX         = 32                    // 递归打开界面次数上限
	GUI_LOG_GIT_PLACEHOLDER        = "git log -..."        // 日志显示条目数提示文字
	GUI_LOG_GIT_DEFAULT            = "10"                  // 日志显示条目数 默认
	GUI_LOG_GIT_TIME_FORMAT        = "2006-01-02 15:04:05" // 日志时间格式
	GUI_INFO_README                = `# 说明

基础云笔记功能

## 功能

- 基础编辑
	- "Ctrl + S",保存
- Markdown
- github/gitee同步
- 强制推送
- 强制拉取
- 导入
- 导出

## 使用说明

1. 配置仓库地址(数据存储)
1. 配置认证信息(提交认证)
	- 账号密码方式
		- 账号
		- 密码
	- token方式
		- 账号
		- 密码:填写token的值
1. 配置作者信息(git日志用)
	- 作者名
	- 作者邮箱

## 问题解决

- 如果出现同步失败的情况
	- 可能的原因
		- 本地更新后其他人提交到仓库了文件
	- 解决方法
		- 强制拉取(会丢弃本地变更)
		- 强制推送(会丢失远程变更)

` // 信息界面说明文字

	// plog
	PLOG_MAX_SIZE = 1 << 20 // 日志数据量
	PLOG_CHAN_LEN = 100     // 数据管道长度

	// notebook
	NOTEBOOK_PERM_FILE        = 0666      // 文件权限
	NOTEBOOK_PERM_FOLDER      = 0755      // 文件夹权限
	NOTEBOOK_TMP_FILE_NAME    = "tmp.txt" // 空项目占位文件,因为go-git无法清空项目提交的问题,需要至少一个文件保证可以提交清空操作
	NOTEBOOK_TMP_FILE_CONTENT = "this file is generated by " + WINDOW_TITLE + " for a go-git bug to commit empty workspace\nyou can treat this file as an regular file\nand do whatever you want."
)
