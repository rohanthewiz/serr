package serr

const UserMsgKey = "userMsgKey"                 // user message key
const UserMsgSeverityKey = "userMsgSeverityKey" // user message severity key

var Severity severity

type severity struct {
	Success, Error, Warn, Info string
}

func init() {
	Severity = severity{
		Success: "success",
		Error:   "error",
		Warn:    "warn",
		Info:    "info",
	}
}

type UserMsgOptions struct {
	Severity string
}
