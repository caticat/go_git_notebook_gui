package main

import (
	"image/color"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func initGUIInfo() fyne.CanvasObject {
	urlRespitory, _ := url.Parse("https://github.com/caticat/go_git_notebook_gui")
	urlLink := widget.NewHyperlink("github.com/caticat/go_git_notebook_gui", urlRespitory)

	return container.NewScroll(
		container.NewVBox(
			widget.NewLabel("Information"),
			widget.NewForm(
				widget.NewFormItem("author", widget.NewLabel(APP_AUTHOR)),
				widget.NewFormItem("repository", urlLink),
				widget.NewFormItem("version", widget.NewLabel(APP_VER)),
			),
			canvas.NewLine(color.Black),
			widget.NewRichTextFromMarkdown(GUI_INFO_README),
		),
	)
}
