package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type LineCounterApp struct {
	window     fyne.Window
	pathEntry  *widget.Entry
	runButton  *widget.Button
	linesCount map[string]int
}

func NewLineCounterApp() *LineCounterApp {
	myApp := app.New()
	window := myApp.NewWindow("Line Counter")

	pathEntry := widget.NewEntry()
	pathEntry.SetPlaceHolder("Input the folder path...")
	pathEntry.Resize(fyne.NewSize(300, 300))

	runButton := widget.NewButton("\n\n\n\n\n\nRUN\n\n\n\n\n\n", func() {
		folderPath := pathEntry.Text
		if folderPath == "" {
			dialog.ShowError(errors.New("you haven't entered a correct path to the folder"), window)
			return
		}

		app := &LineCounterApp{window: window, pathEntry: pathEntry, linesCount: make(map[string]int)}
		if err := app.CountLines(folderPath); err != nil {
			dialog.ShowError(err, window)
			return
		}

		dialog.ShowInformation("Result", app.FormatLinesCount(), window)
	})

	content := container.NewVBox(
		pathEntry,
		runButton,
	)

	window.SetContent(content)
	window.Resize(fyne.NewSize(400, 300))

	return &LineCounterApp{
		window:    window,
		pathEntry: pathEntry,
		runButton: runButton,
		linesCount: make(map[string]int),
	}
}

func (app *LineCounterApp) CountLines(folderPath string) error {
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == "node_modules" {
			return filepath.SkipDir
		}
		if !info.IsDir() {
			ext := filepath.Ext(path)
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			lines := strings.Split(string(content), "\n")
			app.linesCount[ext] += len(lines)
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (app *LineCounterApp) FormatLinesCount() string {
	var result string
	for ext, count := range app.linesCount {
		result += fmt.Sprintf("Extension: %s, Number of lines: %d\n", ext, count)
	}
	return result
}

func (app *LineCounterApp) Run() {
	app.window.ShowAndRun()
}

func main() {
	app := NewLineCounterApp()
	app.Run()
}
