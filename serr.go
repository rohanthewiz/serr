package serr

import (
	"errors"
	"fmt"
	"strings"
)

// SErr is a Structured Error wrapper
// Supports wrapping of errors with a list of key, values to nicely support structured logging
// Works nicely with github.com/rohanthewiz/logger
type SErr struct {
	err error // the usual error
	// support structured logging of the format key1, val1, key2, val2
	// Repeated keys are allowed and will be concatenated in log output
	fields []string
}

// New returns a new SErr as an error type
func New(er string, fields ...string) error {
	se := SErr{err: errors.New(er)}
	return se.newSErr(fields...)
}

// NewSErr returns a new concrete SErr
func NewSErr(er string, fields ...string) SErr {
	ser := SErr{err: errors.New(er)}
	return ser.newSErr(fields...)
}

// AppendKeyValPairs adds pairs of attribute-values to the SErr
// *Note* this method will be used by SErr aware loggers to add extra fields
// at the time of logging, so let's maintain the function signature
func (se *SErr) AppendKeyValPairs(keyValPairs ...string) {
	keyValPairs = fixupFields(keyValPairs) // it doesn't hurt to always fix up fields
	se.fields = append(se.fields, keyValPairs...)
}

// Error satisfies the `error` interface
// The contract here is to return the value of the core error
func (se SErr) Error() string {
	if se.err == nil {
		return "SErr: internal error not set"
	}
	return se.err.Error()
}

// FieldsMap returns all SErr attributes as a map of string keys and values.
// Values of duplicate fields are appended together with ' - '
// such that the innermost attributes are to the right
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

// FieldsAsString builds output for non-structured logging
func (se SErr) FieldsAsString() string {
	mp := se.FieldsMap()
	arr := make([]string, 0, len(mp))
	for key, val := range mp {
		arr = append(arr, key+"->"+val)
	}
	return strings.Join(arr, "; ")
}

// String satisfies the Stringer interface, so this is the default method called by fmt
func (se SErr) String() (out string) {
	return fmt.Sprintf("%s [error_attrs] => %s", se.err, se.FieldsAsString())
}

// Clone returns a new SErr from an existing one
func (se SErr) Clone() SErr {
	return SErr{se.err, se.fields}
}

// GetError returns the wrapped error
func (se SErr) GetError() error {
	return se.err
}

// Unwrap returns the wrapped error
// This is the standard for Go
//
//	see https://blog.golang.org/go1.13-errors#TOC_3.1.
func (se SErr) Unwrap() error {
	return se.err
}

// Fields returns the internal list of keys and values
func (se SErr) Fields() []string {
	return se.fields
}

// appendCallerContext adds Function name and location of the call to SErr new or wrapper functions
func (se *SErr) appendCallerContext(frmLevel int) {
	se.AppendKeyValPairs([]string{"location", FunctionLoc(frmLevel),
		"function", FunctionName(frmLevel)}...)
}

// newSErr is the core method for creating a new SErr from an existing SErr
// This is used in Wrap, New and other methods that add key val pairs and context
func (ser SErr) newSErr(pairs ...string) (out SErr) {
	out = SErr{err: ser.err} // add the internal error

	// Add any existing fields first
	if len(ser.fields) > 0 {
		out.AppendKeyValPairs(ser.fields...) // add existing fields first
	}

	// Add new fields
	out.AppendKeyValPairs(pairs...)

	// Add location info on each wrap
	out.appendCallerContext(FrameLevels.FrameLevel4)
	return
}

// NewSerrNoContext builds an SErr from an err without addition of frame context.
// If err already contains a concrete SErr, it is returned
func NewSerrNoContext(err error) SErr {
	if ser, ok := err.(SErr); !ok {
		return SErr{err: err}
	} else {
		return ser
	}
}

// Wrap wraps an existing error. Attribute keys and values must be strings.
// Returns an SErr (structured err) as an error
// This requires an even number of fields unless a single field is given
// in which case it is added under the key "msg".
func Wrap(err error, fields ...string) error {
	if err == nil {
		fmt.Println("SErr: Not wrapping a nil error", "callerLocation:", FunctionLoc(FrameLevels.FrameLevel2),
			"callerName:", FunctionName(FrameLevels.FrameLevel2))
		return nil
	}

	return NewSerrNoContext(err).newSErr(fields...)
}

// WrapAsSErr wraps an existing error. Attribute keys and values must be strings.
// Returns a concrete SErr (structured err)
// This requires an even number of fields unless a single field is given
// in which case it is added under the key "msg".
func WrapAsSErr(err error, fields ...string) SErr {
	if err == nil {
		fmt.Println("SErr: Not wrapping a nil error", "callerLocation:", FunctionLoc(FrameLevels.FrameLevel2),
			"callerName:", FunctionName(FrameLevels.FrameLevel2))
		return SErr{}
	}

	return NewSerrNoContext(err).newSErr(fields...)
}

//	fixupFields fixes up a  sequence of attribute value pairs
//
// A Single field gets added as {"msg", "field"}
// For an odd number of multiple fields, the first field is considered a message value {"msg", "field"}
// An even number of fields are added without any change in sequence
func fixupFields(fields []string) (flds []string) {
	ln := len(fields)

	if ln == 1 { // Single field becomes a "msg: field" pair
		flds = append(flds, []string{"msg", fields[0]}...)
	} else {
		if ln%2 != 0 { // for odd fields, treat the first as a message
			msg := fields[0]
			fields = fields[1:]                          // drop the first
			flds = append(flds, []string{"msg", msg}...) // add as first pair
		}
		// Add fields
		flds = append(flds, fields...)
	}
	return
}

// SetUserMsg is a convenience method for setting a user message field
// This could be displayed to the user of the app
func (se *SErr) SetUserMsg(msg string, sev string) {
	userInfo := []string{UserMsgKey, msg, UserMsgSeverityKey, sev}
	se.fields = append(se.fields, userInfo...)
}

// UserMsg is a convenience method to return the user message field
// This could be a message displayed to the user of the app
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
