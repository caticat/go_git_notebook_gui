package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"github.com/caticat/go_game_server/pgit"
)

var (
	g_app             fyne.App
	g_win             fyne.Window
	g_cfg             = newPNoteBookCfg()
	g_pgit            *pgit.PGit
	g_logData         = binding.NewString()
	g_logLast         = binding.NewString()
	g_files           *PFyneTreePath
	g_pathSelect      string
	g_openFileName    string
	g_openFileContent string
	g_funRefresh      func()
)

func getApp() fyne.App                  { return g_app }
func setApp(a fyne.App)                 { g_app = a }
func getWin() fyne.Window               { return g_win }
func setWin(w fyne.Window)              { g_win = w }
func getCfg() *PNoteBookCfg             { return g_cfg }
func getPGit() *pgit.PGit               { return g_pgit }
func setPGit(p *pgit.PGit)              { g_pgit = p }
func getLogData() binding.String        { return g_logData }
func getLogLast() binding.String        { return g_logLast }
func getFiles() *PFyneTreePath          { return g_files }
func setFiles(files *PFyneTreePath)     { g_files = files }
func getPathSelect() string             { return g_pathSelect }
func setPathSelect(pathSelect string)   { g_pathSelect = pathSelect }
func getOpenFileName() string           { return g_openFileName }
func setOpenFileName(name string)       { g_openFileName = name }
func getOpenFileContent() string        { return g_openFileContent }
func setOpenFileContent(content string) { g_openFileContent = content }
func getFunRefresh() func()             { return g_funRefresh }
func setFunRefresh(fun func())          { g_funRefresh = fun }
