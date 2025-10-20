package addkeyui

import (
	"DatabaseDB/internal/logic"
	"io/ioutil"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func OpenWindowAddButton(rightColumnContent *fyne.Container, mainWindow fyne.Window) {
	var ded *dialog.CustomDialog

	iputKey := widget.NewEntry()
	iputKey.SetPlaceHolder("Key")

	iputvalue := widget.NewMultiLineEntry()
	iputvalue.SetPlaceHolder("value")

	nameFile := widget.NewButton("Name File", nil)

	var valueFinish []byte
	uploadFile := widget.NewButton("UploadFile", func() {
		folderPath := dialog.NewFileOpen(func(dir fyne.URIReadCloser, err error) {
			if err != nil {
				log.Println("Error opening folder:", err)
				return
			}
			if dir == nil {
				log.Println("No folder selected")
				return
			}

			filename := dir.URI().Name()

			valueFinish, err = ioutil.ReadAll(dir)
			if err != nil {
				log.Println(err.Error())
				return
			}

			nameFile.SetText(filename)
			nameFile.Refresh()
		}, mainWindow)
		folderPath.SetFilter(storage.NewExtensionFileFilter([]string{
			".ico", ".svg",
			".jpeg", ".jpg", ".png", ".txt", ".json", ".go",
			".md", ".xml", ".csv", ".ini", ".yml", ".yaml", ".log", ".config", ".properties",
			".env", ".sql", ".xml", ".json5", ".rst", ".tex", ".asm", ".hbs", ".tpl", ".html",
			".conf", ".mdx", ".latex", ".scala", ".swift", ".lua", ".ts", ".scss", ".less",
			".asm", ".awk", ".bat", ".csh", ".c", ".cpp", ".h", ".java", ".kt", ".lisp", ".css",
			".m", ".pas", ".pl", ".php", ".ps", ".ps1", ".py", ".r", ".rb", ".sh", ".sql",
			".tcl", ".vbs", ".vhd", ".vue", ".yaml", ".yml", ".zsh", ".coffee", ".clj", ".js",
			".erl", ".fs", ".dart", ".handlebars", ".scss", ".sass", ".mustache", ".jinja",
			".asciidoc", ".org", ".tex", ".rst", ".sml", ".v", ".verilog", ".vhdl", ".scala",
			".swift", ".m4", ".xhtml", ".xml5", ".wsdl", ".xsd", ".dtd", ".gdsl", ".jsonc",
			".hbs", ".hs", ".limbo", ".loco", ".ml", ".nim", ".oz", ".pddl", ".rexx", ".rmd",
			".sh", ".tcl", ".xsl", ".yml",
		}))
		folderPath.Show()
	})

	uploadFile.Disable()
	iputvalue.Disable()
	nameFile.Disable()

	typeValue := widget.NewLabel("Select the type of file you want")
	redioType := widget.NewRadioGroup([]string{"Text", "File"}, func(typeRedio string) {
		switch typeRedio {
		case "Text":
			iputvalue.Enable()
			uploadFile.Disable()
			nameFile.Disable()
		case "File":
			uploadFile.Enable()
			nameFile.Enable()
			iputvalue.Disable()

		}
	})

	redioType.Horizontal = true
	rowRedio := container.NewHBox(typeValue, redioType)

	columns := container.NewHSplit(uploadFile, nameFile)
	columns.SetOffset(0.80)

	ButtonAddAdd := widget.NewButton("Add", func() {
		if uploadFile.Disabled() {
			valueFinish = []byte(iputvalue.Text)
		}
		err := logic.AddKeyLogic(iputKey.Text, valueFinish)
		if err != nil {
			dialog.ShowInformation("Error", err.Error(), mainWindow)
		} else {
			ded.Hide()
		}
	})
	ButtonAddAdd.Importance = widget.HighImportance
	cont := container.NewVBox(
		iputKey,
		rowRedio,
		iputvalue,
		columns,
		layout.NewSpacer(),
		ButtonAddAdd,
		layout.NewSpacer(),
	)

	ded = dialog.NewCustom("Add Key and Value", "Close", cont, mainWindow)
	ded.Resize(fyne.NewSize(600, 400))
	ded.Show()

}
