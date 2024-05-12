package serr

const UserMsgKey = "userMsgKey"                 // user message key
const UserMsgSeverityKey = "userMsgSeverityKey" // user message severity key

type UserMsgOptions struct {
	Severity string
}

// poor man's enum
type severity struct {
	Success, Error, Warn, Info string
}

var Severity = severity{"success", "error", "warn", "info"}

// UserMsg is a convenience function for getting the user message,
// and severity fields from a standard error
// This could be a message displayed to the user of the app
func UserMsg(err error) (msg, severity string) {
	if err == nil {
		return
	}
	if ser, ok := err.(SErr); ok {
		msg, severity = ser.UserMsg()
	}
	return
}
