package dialog

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Custom struct {
	*dialog
}

// NewCustom creates and returns a dialog over the specified application using custom
// content. The button will have the dismiss text set.
// The MinSize() of the CanvasObject passed will be used to set the size of the window.
func NewCustom(title, dismiss string, content fyne.CanvasObject, parent fyne.Window) *Custom {
	d := &dialog{content: content, title: title, parent: parent}
	d.layout = &dialogLayout{d: d}

	d.dismiss = &widget.Button{Text: dismiss,
		OnTapped: d.Hide,
	}
	d.create(container.NewHBox(layout.NewSpacer(), d.dismiss, layout.NewSpacer()))

	return &Custom{dialog: d}
}

// ShowCustom shows a dialog over the specified application using custom
// content. The button will have the dismiss text set.
// The MinSize() of the CanvasObject passed will be used to set the size of the window.
func ShowCustom(title, dismiss string, content fyne.CanvasObject, parent fyne.Window) {
	NewCustom(title, dismiss, content, parent).Show()
}

// CustomConfirm is a custom dialog with dismiss and confirm buttons.
//
// Since 2.4
type CustomConfirm struct {
	*dialog

	confirm *widget.Button
}

// SetConfirmImportance sets the importance level of the confirm button.
//
// Since 2.4
func (c *CustomConfirm) SetConfirmImportance(importance widget.ButtonImportance) {
	c.confirm.Importance = importance
}

// NewCustomConfirm creates and returns a dialog over the specified application using
// custom content. The cancel button will have the dismiss text set and the "OK" will
// use the confirm text. The response callback is called on user action.
// The MinSize() of the CanvasObject passed will be used to set the size of the window.
func NewCustomConfirm(title, confirm, dismiss string, content fyne.CanvasObject,
	callback func(bool), parent fyne.Window) *CustomConfirm {
	d := &dialog{content: content, title: title, parent: parent}
	d.layout = &dialogLayout{d: d}
	d.callback = callback

	d.dismiss = &widget.Button{Text: dismiss, Icon: theme.CancelIcon(),
		OnTapped: d.Hide,
	}
	ok := &widget.Button{Text: confirm, Icon: theme.ConfirmIcon(), Importance: widget.HighImportance,
		OnTapped: func() {
			d.hideWithResponse(true)
		},
	}
	d.create(container.NewHBox(layout.NewSpacer(), d.dismiss, ok, layout.NewSpacer()))

	return &CustomConfirm{dialog: d, confirm: ok}
}

// ShowCustomConfirm shows a dialog over the specified application using custom
// content. The cancel button will have the dismiss text set and the "OK" will use
// the confirm text. The response callback is called on user action.
// The MinSize() of the CanvasObject passed will be used to set the size of the window.
func ShowCustomConfirm(title, confirm, dismiss string, content fyne.CanvasObject,
	callback func(bool), parent fyne.Window) {
	NewCustomConfirm(title, confirm, dismiss, content, callback, parent).Show()
}
