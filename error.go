package main

import "errors"

var (
	ErrAppNotInit                   = errors.New("app not init")
	ErrSaveFileNameEmpty            = errors.New("save file name is empty")
	ErrGitNotSync                   = errors.New("git is not sync, change config or clear local data")
	ErrConfigNotFound               = errors.New("config not found")
	ErrCreateFileFolderAlreadyExist = errors.New("can not create file/folder already exist")
	ErrDeleteRootDataDirInHome      = errors.New("can not delete root dir by home delete button,\nuse setting->advance->delete-all-data instead")
	ErrDeleteFileFolderNotFound     = errors.New("file/folder not found while delete")
	ErrMoveRootDataDirInHome        = errors.New("can not move root dir by home move button")
	ErrMoveFileFolderNotFound       = errors.New("file/folder not found while move")
	ErrMoveFilePathOutOfData        = errors.New("can not move file out of data dir")
)
