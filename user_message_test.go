package serr

import (
	"errors"
	"testing"
)

func TestUserMsg(t *testing.T) {
	const strErr1 = "This is a test err"
	const thisIsMyMessage = "This is my message"

	er := Wrap(errors.New(strErr1), thisIsMyMessage)
	se, ok := er.(SErr)
	if !ok {
		t.Error("er should be a SErr")
		t.FailNow()
	}
	sl := se.Fields()
	if len(sl) != 6 {
		t.Error("Structured error from an error wrapped with a single field should contain 6 fields, got", len(sl))
	}

	const umsg = "Your app needs to be updated"
	se.SetUserMsg(umsg, Severity.Warn)
	if msg, sev := UserMsg(se); msg != umsg || sev != Severity.Warn {
		t.Errorf(`User message or severity is not as expected.
			Expected message %s, Got %s; Expected Severity %s, Got %s`, umsg, msg, Severity.Warn, sev)
	}
}

func TestUserMsgFromErr(t *testing.T) {
	tests := []struct {
		name     string
		errMsg   string
		usrMsg   string
		altMsg   string
		expected string
	}{
		{
			name:     "Empty error",
			expected: "",
		},
		{
			name:     "Error is blank, usrMsg provided",
			usrMsg:   "User message",
			expected: "User message",
		},
		{
			name:     "Only errMsg provided",
			errMsg:   "Error message",
			usrMsg:   "",
			expected: "",
		},
		{
			name:     "Error is SErr, no alt provided",
			usrMsg:   "Some user message",
			expected: "Some user message",
		},
		{
			name:     "User msg is blank, alt provided",
			usrMsg:   "",
			altMsg:   "Alt message",
			expected: "Alt message",
		},
	}

	for _, test := range tests {
		var testErr error

		// We must have an error to work with
		err := errors.New(test.errMsg)

		// Make it an SErr
		ser := WrapAsSErr(err)
		ser.SetUserMsg(test.usrMsg, Severity.Info)
		testErr = ser

		t.Run(test.name, func(t *testing.T) {
			result := UserMsgFromErr(testErr, test.altMsg)
			if result != test.expected {
				t.Errorf("Expected '%s', got '%s'", test.expected, result)
			}
		})
	}
}
