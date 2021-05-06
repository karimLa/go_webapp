package views

import "webapp/models"

const (
	AlertLvlError    = "danger"
	AlertLvlWarning  = "warning"
	AlertLevlInfo    = "info"
	AlertLevelSucess = "success"

	AlertMsgGeneric = "Something went wrong. Please try again, and contact us if the problem presists."
)

// Alert is used to render Bootstrap Alert messages in templates
type Alert struct {
	Level   string
	Message string
}

// Data is the top level structure that views expect data
// to come in.
type Data struct {
	Alert *Alert
	Yield interface{}
	User  *models.User
}

func (d *Data) SetAlert(err error) {
	if pErr, ok := err.(PublicError); ok {
		d.AlertError(pErr.Public())
	} else {
		d.AlertError(AlertMsgGeneric)
	}
}

func (d *Data) AlertError(msg string) {
	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: msg,
	}
}

type PublicError interface {
	error
	Public() string
}
