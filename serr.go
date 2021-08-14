package serr

import (
	"errors"
	"fmt"
	"strings"
)

// Backend Structured Error wrapper
// Supports wrapping of errors with a list of key, values to nicely support structured logging
// Works nicely with logger.LogErr
type SErr struct {
	err error // the usual error
	// support structured logging of the format key1, val1, key2, val2
	// Repeated keys are allowed and will be concatenated in log output
	fields []string
}

// New returns a new SErr as an error type
func New(er string, fields ...string) error {
	se := SErr{err: errors.New(er)}
	se.fields = fixupFields(fields)
	// Add location info on each wrap
	se.Append([]string{"location", FunctionLoc()}...)
	return se
}

// NewSerr returns a new concrete SErr
func NewSErr(er string, fields ...string) SErr {
	se := SErr{err: errors.New(er)}
	se.fields = fixupFields(fields)
	// Add location info on each wrap
	se.Append([]string{"location", FunctionLoc()}...)
	return se
}

// Append adds variable number of strings to the SErr
// These should be key value pairs
func (se *SErr) Append(fields ...string) {
	se.fields = append(se.fields, fields...)
}

// AppendIfHasErr adds variable number of strings to the SErr
// These should be key value pairs on condition
// that the wrapped error is not nil
func (se *SErr) AppendIfHasErr(fields ...string) {
	if se.err != nil {
		se.fields = append(se.fields, fields...)
	}
}

// Yes we are also an error type -- sweet!
// Satisfy the `error` interface
// The contract here is to return the value of the core error
func (se SErr) Error() string {
	if se.err == nil {
		return ""
	}
	return se.err.Error()
}

// Return all SErr attributes as a map of string keys and values
func (se SErr) FieldsMap() map[string]string {
	flds := map[string]string{}
	key := ""
	for i, str := range se.fields {
		if i%2 == 0 { // even indices are presumed to be keys
			key = str
		} else {
			if orig, ok := flds[key]; ok { // we've seen this key before
				flds[key] = str + " - " + orig
			} else {
				flds[key] = str
			}
		}
	}
	return flds
}

// Build output for non-structured logging
func (se SErr) FieldsString() string {
	mp := se.FieldsMap()
	arr := make([]string, 0, len(mp))
	for key, val := range mp {
		arr = append(arr, key+"->"+val)
	}
	return strings.Join(arr, "; ")
}

// Satisfies the Stringer interface
func (se SErr) String() (out string) {
	return fmt.Sprintf("%s --> %s", se.err, se.FieldsString())
}

// Use case: we want to use the convenience functions of SErr
// to build an error then assign it to an existing SErr
func (se SErr) Clone() SErr {
	return SErr{se.err, se.fields}
}

// Return the wrapped error
func (se SErr) GetError() error {
	return se.err
}

// Return the wrapped error
// I believe this is the standard for Go
// 	  see https://blog.golang.org/go1.13-errors#TOC_3.1.
func (se SErr) Unwrap() error {
	return se.err
}

// Return the internal list of keys and values
func (se SErr) Fields() []string {
	return se.fields
}

// Convenience method for setting a user message field
// This is a message displayable to the user of the app
func (se *SErr) SetUserMsg(msg string, sev string) {
	userInfo := []string{UserMsgKey, msg, UserMsgSeverityKey, sev}
	se.fields = append(se.fields, userInfo...)
}

// Convenience method to return the user message field
// This is a message displayable to the user of the app
func (se SErr) UserMsg() (userMsg, severity string) {
	mp := se.FieldsMap()
	if str, ok := mp[UserMsgKey]; ok {
		userMsg = str
	}
	if str, ok := mp[UserMsgSeverityKey]; ok {
		severity = str
	}
	return
}

// Convenience function for getting the user message, and severity fields
// from a standard error
// This is a message displayable to the user of the app
func UserMsg(err error) (msg, severity string) {
	if ser, ok := err.(SErr); ok {
		msg, severity = ser.UserMsg()
	}
	return
}

// Wrap an existing error. Attribute keys and values must be strings.
// Returns an SErr (structured err) as an error
// This requires an even number of fields unless a single field is given
// in which case it is added under the key "msg".
func Wrap(err error, fields ...string) error {
	if err == nil {
		println("SErr: Not wrapping a nil error", "called from", FunctionLoc())
		return nil
	}

	newSErr := SErr{err, []string{}}

	// Add any existing fields first
	if se, ok := err.(SErr); ok && len(se.fields) > 0 {
		newSErr.Append(se.fields...) // add existing fields first
	}

	// Add new fields
	newSErr.Append(fixupFields(fields)...)

	// Add location info on each wrap
	newSErr.Append([]string{"location", FunctionLoc()}...)

	return newSErr // return
}

// Wrap an existing error. Attribute keys and values must be strings.
// Returns a concrete SErr (structured err)
// This requires an even number of fields unless a single field is given
// in which case it is added under the key "msg".
func WrapAsSErr(err error, fields ...string) SErr {
	if err == nil {
		println("SErr: Not wrapping a nil error", "called from", FunctionLoc())
		return SErr{}
	}

	newSErr := SErr{err, []string{}}

	// Add any existing fields first
	if se, ok := err.(SErr); ok && len(se.fields) > 0 {
		newSErr.Append(se.fields...) // add existing fields first
	}

	// Add new fields
	newSErr.Append(fixupFields(fields)...)

	// Add location info on each wrap
	newSErr.Append([]string{"location", FunctionLoc()}...)

	return newSErr // return
}

// Fix up sequence of attribute value pairs
// A Single field gets added as {"msg", "field"}
// For an odd number of multiple fields, the first field is considered a message value {"msg", "field"}
// An even number of fields are added without any change in sequence
func fixupFields(fields []string) (flds []string) {
	ln := len(fields)
	if ln == 1 { // Single field becomes a "msg: field" pair
		flds = append(flds, []string{"msg", fields[0]}...)
	} else {
		if ln%2 != 0 { // for odd fields, first is a message
			msg := fields[0]
			fields = fields[1:] // drop the first
			flds = append(flds, []string{"msg", msg}...)
		}
		flds = append(flds, fields...)
	}
	return
}
