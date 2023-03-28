package main

import (
	"io"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type config struct {
	EditWidget    *widget.Entry
	PreviewWidget *widget.RichText
	CurrentFile   fyne.URI
	SaveMenuItem  *fyne.MenuItem
}

var cfg config

func main() {
	// create a fyne app
	a := app.New()
	// create a window for app
	win := a.NewWindow("Markdown Editor")
	// get user interface
	edit, preview := cfg.makeUI()
	cfg.createMenuItems(win)
	// set the content of the window
	win.SetContent(container.NewHSplit(edit, preview))
	// show the window and run
	win.Resize(fyne.NewSize(800, 600))
	win.CenterOnScreen()
	win.ShowAndRun()
}

func (app *config) makeUI() (*widget.Entry, *widget.RichText) {
	edit := widget.NewMultiLineEntry()
	preview := widget.NewRichTextFromMarkdown("")
	app.EditWidget = edit
	app.PreviewWidget = preview
	edit.OnChanged = preview.ParseMarkdown
	return edit, preview
}

func (app *config) createMenuItems(win fyne.Window) {
	openMenuItem := fyne.NewMenuItem("Open...", app.openFunc(win))
	saveMenuItem := fyne.NewMenuItem("Save", app.saveFunc(win))
	app.SaveMenuItem = saveMenuItem
	app.SaveMenuItem.Disabled = true
	saveAsMenuItem := fyne.NewMenuItem("Save As...", app.saveAsFunc(win))
	fileMenu := fyne.NewMenu("File", openMenuItem, saveMenuItem, saveAsMenuItem)
	menu := fyne.NewMainMenu(fileMenu)
	win.SetMainMenu(menu)
}

var filter = storage.NewExtensionFileFilter([]string{".md", ".MD"})

func (app *config) saveFunc(win fyne.Window) func() {
	return func() {
		if app.CurrentFile == nil {
			app.saveAsFunc(win)()
			return
		}

		write, err := storage.Writer(app.CurrentFile)
		if err != nil {
			dialog.ShowError(err, win)
			return
		}
		defer write.Close()

		// save file
		write.Write([]byte(app.EditWidget.Text))
	}
}

func (app *config) saveAsFunc(win fyne.Window) func() {
	return func() {
		saveDialog := dialog.NewFileSave(func(write fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if write == nil {
				// user canceled
				return
			}

			if !strings.HasSuffix(strings.ToLower(write.URI().String()), ".md") {
				dialog.ShowInformation("Invalid file extension", "Please use .md extension", win)
				return
			}

			// save file
			write.Write([]byte(app.EditWidget.Text))
			app.CurrentFile = write.URI()
			defer write.Close()

			win.SetTitle(win.Title() + " - " + app.CurrentFile.Name())

			// enable save menu item
			app.SaveMenuItem.Disabled = false
		}, win)

		saveDialog.SetFilter(filter)
		saveDialog.Show()
	}
}

func (app *config) openFunc(win fyne.Window) func() {
	return func() {
		openDialog := dialog.NewFileOpen(func(read fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if read == nil {
				// user canceled
				return
			}

			// open file
			defer read.Close()
			data, err := io.ReadAll(read)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			app.EditWidget.SetText(string(data))
			app.CurrentFile = read.URI()
			win.SetTitle(win.Title() + " - " + app.CurrentFile.Name())
			app.SaveMenuItem.Disabled = false
		}, win)

		openDialog.SetFilter(filter)
		openDialog.Show()
	}
}
