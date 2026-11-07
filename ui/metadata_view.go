package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"idea_app/db"
	"idea_app/models"
)

// BuildMetaDataTree builds an interactive hierarchical list of metadata for an Idea, Memory, or MetaData
func BuildMetaDataTree(w fyne.Window, parentID string, parentType string, onRefresh func()) fyne.CanvasObject {
	metaList, err := db.GetMetaDataForParent(parentID, parentType)
	if err != nil {
		return widget.NewLabel("Failed to load metadata: " + err.Error())
	}

	content := container.NewVBox()

	addTopLevelBtn := widget.NewButtonWithIcon("Add Sub-Thought / Note", theme.ContentAddIcon(), func() {
		ShowAddMetaDataModal(w, parentID, parentType, 1, onRefresh)
	})
	content.Add(container.NewHBox(
		widget.NewLabelWithStyle("Sub-Thoughts & Metadata", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		layoutSpacer(),
		addTopLevelBtn,
	))

	if len(metaList) == 0 {
		content.Add(widget.NewLabelWithStyle("No sub-thoughts or notes added yet.", fyne.TextAlignCenter, fyne.TextStyle{Italic: true}))
		return content
	}

	for _, m := range metaList {
		renderMetaDataNode(content, w, m, onRefresh)
	}

	return content
}

func renderMetaDataNode(parentBox *fyne.Container, w fyne.Window, m *models.MetaData, onRefresh func()) {
	// Indentation prefix based on level (level 1 to 6)
	indentPrefix := strings.Repeat("   ", m.Level-1)
	levelBadge := fmt.Sprintf("[L%d] ", m.Level)

	titleLabel := widget.NewLabelWithStyle(indentPrefix+levelBadge+m.Title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	descLabel := widget.NewLabel(indentPrefix + "    " + m.Description)

	var attachObjects []fyne.CanvasObject
	if len(m.PhotoPaths) > 0 || len(m.FilePaths) > 0 {
		var attachText []string
		for _, p := range m.PhotoPaths {
			attachText = append(attachText, "📷 "+p)
		}
		for _, f := range m.FilePaths {
			attachText = append(attachText, "📎 "+f)
		}
		attachLabel := widget.NewLabelWithStyle(indentPrefix+"    Attachments: "+strings.Join(attachText, ", "), fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
		attachObjects = append(attachObjects, attachLabel)
	}

	btnBox := container.NewHBox(
		widget.NewLabel(indentPrefix),
	)

	// Can add sub-level if less than level 6
	if m.Level < 6 {
		addSubBtn := widget.NewButtonWithIcon("Add Sub", theme.ContentAddIcon(), func() {
			ShowAddMetaDataModal(w, m.ID, "MetaData", m.Level+1, onRefresh)
		})
		btnBox.Add(addSubBtn)
	}

	addAttachBtn := widget.NewButtonWithIcon("Attach", theme.MailAttachmentIcon(), func() {
		ShowAddAttachmentModal(w, m.ID, onRefresh)
	})
	btnBox.Add(addAttachBtn)

	editBtn := widget.NewButtonWithIcon("Edit", theme.DocumentCreateIcon(), func() {
		ShowEditMetaDataModal(w, m, onRefresh)
	})
	btnBox.Add(editBtn)

	delBtn := widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), func() {
		dialog.ShowConfirm("Delete MetaData", "Delete this item and all its nested sub-thoughts?", func(confirmed bool) {
			if confirmed {
				db.DeleteMetaData(m.ID)
				onRefresh()
			}
		}, w)
	})
	btnBox.Add(delBtn)

	nodeBox := container.NewVBox(titleLabel)
	if m.Description != "" {
		nodeBox.Add(descLabel)
	}
	for _, att := range attachObjects {
		nodeBox.Add(att)
	}
	nodeBox.Add(btnBox)
	nodeBox.Add(widget.NewSeparator())

	parentBox.Add(nodeBox)

	// Recursively render child sub-metadata
	for _, sub := range m.SubMetaData {
		renderMetaDataNode(parentBox, w, sub, onRefresh)
	}
}

func ShowAddMetaDataModal(w fyne.Window, parentID string, parentType string, level int, onSaved func()) {
	titleEntry := widget.NewEntry()
	titleEntry.SetPlaceHolder("Sub-thought / Note Title")

	descEntry := widget.NewMultiLineEntry()
	descEntry.SetPlaceHolder("Description / Content")

	form := widget.NewForm(
		widget.NewFormItem("Level", widget.NewLabel(fmt.Sprintf("%d of 6", level))),
		widget.NewFormItem("Title", titleEntry),
		widget.NewFormItem("Description", descEntry),
	)

	d := dialog.NewCustomConfirm(fmt.Sprintf("New Sub-Thought (Level %d)", level), "Add", "Cancel", form, func(save bool) {
		if save && strings.TrimSpace(titleEntry.Text) != "" {
			meta := &models.MetaData{
				ID:          models.NewID(),
				ParentID:    parentID,
				ParentType:  parentType,
				Level:       level,
				Title:       strings.TrimSpace(titleEntry.Text),
				Description: strings.TrimSpace(descEntry.Text),
			}
			if err := db.CreateMetaData(meta); err != nil {
				dialog.ShowError(err, w)
				return
			}
			onSaved()
		}
	}, w)
	d.Resize(fyne.NewSize(380, 280))
	d.Show()
}

func ShowEditMetaDataModal(w fyne.Window, m *models.MetaData, onSaved func()) {
	titleEntry := widget.NewEntry()
	titleEntry.SetText(m.Title)

	descEntry := widget.NewMultiLineEntry()
	descEntry.SetText(m.Description)

	form := widget.NewForm(
		widget.NewFormItem("Title", titleEntry),
		widget.NewFormItem("Description", descEntry),
	)

	d := dialog.NewCustomConfirm("Edit Sub-Thought", "Save", "Cancel", form, func(save bool) {
		if save && strings.TrimSpace(titleEntry.Text) != "" {
			m.Title = strings.TrimSpace(titleEntry.Text)
			m.Description = strings.TrimSpace(descEntry.Text)
			if err := db.UpdateMetaData(m); err != nil {
				dialog.ShowError(err, w)
				return
			}
			onSaved()
		}
	}, w)
	d.Resize(fyne.NewSize(380, 260))
	d.Show()
}

func ShowAddAttachmentModal(w fyne.Window, metadataID string, onSaved func()) {
	pathEntry := widget.NewEntry()
	pathEntry.SetPlaceHolder("/path/to/image.png or file.pdf")

	isPhotoCheck := widget.NewCheck("Is Photo / Image", nil)
	isPhotoCheck.SetChecked(true)

	browseBtn := widget.NewButtonWithIcon("Browse File", theme.FolderOpenIcon(), func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				pathEntry.SetText(reader.URI().Path())
				reader.Close()
			}
		}, w)
		fd.Show()
	})

	form := widget.NewForm(
		widget.NewFormItem("File Path", pathEntry),
		widget.NewFormItem("", browseBtn),
		widget.NewFormItem("Type", isPhotoCheck),
	)

	d := dialog.NewCustomConfirm("Add Attachment", "Attach", "Cancel", form, func(save bool) {
		if save && strings.TrimSpace(pathEntry.Text) != "" {
			err := db.AddAttachment(metadataID, strings.TrimSpace(pathEntry.Text), isPhotoCheck.Checked, models.NewID())
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			onSaved()
		}
	}, w)
	d.Resize(fyne.NewSize(380, 240))
	d.Show()
}

func layoutSpacer() fyne.CanvasObject {
	return widget.NewLabel("")
}
