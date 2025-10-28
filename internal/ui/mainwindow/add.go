package mainwindow

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

func (mw *MainWindow2) OpenAddDialog() {
	var addDialog *dialog.CustomDialog
	var addButton *widget.Button

	keyEntry := widget.NewEntry()
	keyEntry.SetPlaceHolder("Enter Key")

	valueEntry := widget.NewMultiLineEntry()
	valueEntry.SetPlaceHolder("Enter Value")

	fileNameButton := widget.NewButton("File Name", nil)

	var fileData []byte
	uploadFileButton := widget.NewButton("Upload File", func() {
		fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				log.Println("Error opening file:", err)
				return
			}
			if reader == nil {
				log.Println("No file selected")
				return
			}

			filename := reader.URI().Name()

			fileData, err = ioutil.ReadAll(reader)
			if err != nil {
				log.Println("Error reading file:", err)
				return
			}

			fileNameButton.SetText(filename)
			addButton.Enable()
			fileNameButton.Refresh()
		}, mw.Window)

		fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{
			".ico", ".svg", ".jpeg", ".jpg", ".png", ".txt", ".json", ".go",
			".md", ".xml", ".csv", ".ini", ".yml", ".yaml", ".log", ".config",
			".properties", ".env", ".sql", ".json5", ".rst", ".tex", ".asm",
			".hbs", ".tpl", ".html", ".conf", ".mdx", ".latex", ".scala", ".swift",
			".lua", ".ts", ".scss", ".less", ".awk", ".bat", ".csh", ".c", ".cpp",
			".h", ".java", ".kt", ".lisp", ".css", ".m", ".pas", ".pl", ".php",
			".ps", ".ps1", ".py", ".r", ".rb", ".sh", ".tcl", ".vbs", ".vhd",
			".vue", ".coffee", ".clj", ".js", ".erl", ".fs", ".dart", ".handlebars",
			".sass", ".mustache", ".jinja", ".asciidoc", ".org", ".sml", ".v",
			".verilog", ".vhdl", ".m4", ".xhtml", ".xml5", ".wsdl", ".xsd", ".dtd",
			".gdsl", ".jsonc", ".hs", ".limbo", ".loco", ".ml", ".nim", ".oz",
			".pddl", ".rexx", ".rmd", ".xsl", ".zsh",
		}))
		fileDialog.Show()
	})

	valueEntry.OnChanged = func(s string) {
		if s != "" {
			addButton.Enable()
		} else {
			addButton.Disable()
		}
	}

	// Default disabled state
	uploadFileButton.Disable()
	valueEntry.Disable()
	fileNameButton.Disable()

	fileTypeLabel := widget.NewLabel("Select file type:")
	fileTypeRadio := widget.NewRadioGroup([]string{"Text", "File"}, func(selected string) {
		addButton.Disable()
		switch selected {
		case "Text":
			fileData = nil
			valueEntry.Enable()
			uploadFileButton.Disable()
			uploadFileButton.SetText("Upload File")
			fileNameButton.SetText("File Name")
			fileNameButton.Disable()
		case "File":
			fileData = nil
			valueEntry.SetText("")
			uploadFileButton.Enable()
			fileNameButton.Enable()
			valueEntry.Disable()
		}
	})
	fileTypeRadio.Horizontal = true

	fileTypeRow := container.NewHBox(fileTypeLabel, fileTypeRadio)
	fileButtonsRow := container.NewHSplit(uploadFileButton, fileNameButton)
	fileButtonsRow.SetOffset(0.8)

	addButton = widget.NewButton("Add", func() {
		if uploadFileButton.Disabled() {
			fileData = []byte(valueEntry.Text)
		}

		if len(keyEntry.Text) == 0 && len(fileData) == 0 {
			dialog.ShowInformation("Error", "Key and Value cannot be empty", mw.Window)
			return
		}

		err := logic.AddKeyLogic(keyEntry.Text, fileData)
		if err != nil {
			dialog.ShowInformation("Error", err.Error(), mw.Window)
			return
		}

		addDialog.Hide()
		mw.RightColumn.Container.Refresh()
	})
	addButton.Importance = widget.HighImportance

	addButton.Disable()

	content := container.NewVBox(
		keyEntry,
		fileTypeRow,
		valueEntry,
		fileButtonsRow,
		layout.NewSpacer(),
		addButton,
		layout.NewSpacer(),
	)

	addDialog = dialog.NewCustom("Add Key and Value", "Close", content, mw.Window)
	addDialog.Resize(fyne.NewSize(600, 400))
	addDialog.Show()
}
