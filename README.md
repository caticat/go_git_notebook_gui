# go_git_notebook_gui

a git based notebook gui implement by go

## 待制作

## 已实现

- 新建
	- 文件
	- 文件夹
- 删除
	- 文件
	- 文件夹
- 移动
	- 文件
	- 文件夹
- 导入
- 导出
- 预览界面
	- markdown支持
- 搜索
- 全部删除并推送
- 强制推送
- 强制拉取
- 启动时无法连接git时跳转到配置页签
- 首次打开编辑窗口处理

## 更新内容

### v0.0.2

- 手机目录树展开死循环问题修正

## 打包

- `fyne package -os windows -appID github.com.caticat.go_git_notebook_gui -icon assets/myapp.png`
- `fyne package -os android -appID github.com.caticat.go_git_notebook_gui -icon assets/myapp.png`
