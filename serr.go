package serr

import (
	"fmt"
	"errors"
)

// Structured Error wrapper
// Supports wrapping of errors with a list of key, values to nicely support structured logging
// Works nicely with logger.LogErr
type SErr struct {
	err error  // the usual error
	// support structured logging of the format key1, val1, key2, val2
	// Repeated keys are allowed and should be concatenated in log output
	fields []string
}

func NewSErr(err error, fields []string) error {
	se := SErr{}
	se.err = err
	se.fields = fields
	return se
}

// Yes we are also an error type -- sweet!
// Satisfy the `error` interface
func (se SErr) Error() string {
	if se.err == nil { return "" }
	return se.err.Error()
}

// Return all SErr attributes as a map of string keys and values
func (se SErr) FieldsMap() map[string]string {
	flds := map[string]string{}
	key := ""
	lenFields := len(se.fields)
	for i, str := range se.fields {
		if i % 2 == 0 {
			if i == lenFields - 1 {
				fmt.Println("[SErr] Key:", str, "has no matching value", "location:", FunctionLoc(), "fields:", se.fields)
				break  // this the last item - we don't have a matching value, so drop
			}
			key = str
		} else {
			if orig, ok := flds[key]; ok {  // we've seen this key before - accumulate
				flds[key] = str + " - " + orig
			} else {
				flds[key] = str
			}
		}
	}
	return flds
}

// Wrap an existing error. Attribute keys and values must be strings.
// Returns an SErr (structured err)
// This requires an even number of fields unless a single field is given
// in which case it is added under the key "msg".
func Wrap(err error, fields ... string) error {
	if err == nil {
		err = errors.New("[SErr] Not wrapping a nil error from " + FunctionLoc())
		fmt.Println(err.Error(), "location:", FunctionLoc(), "fields:", fields)
		return err
	}
	if ln := len(fields); ln > 1 && ln % 2 != 0 {
		err = errors.New("[SErr] Odd number of fields provided from " + FunctionLoc())
		fmt.Println(err.Error(), "location:", FunctionLoc(), "fields:", fields)
		return err
	}

	var flds []string

	// Add Existing fields
	if se, ok := err.(SErr); ok && len(se.fields) > 0 {
		flds = append(flds, se.fields...)
	}

	// Add new fields
	if len(fields) == 1 {
		flds = append(flds, []string{"msg", fields[0]}...)
	} else {
		flds = append(flds, fields...)
	}

	// Add location info on each wrap
	flds = append(flds, "location")
	flds = append(flds, FunctionLoc())

	return SErr{err, flds} // return
}

// Return the wrapped error
func (se SErr) OriginalErr() error {
	return se.err
}

// Return the internal list of keys and values
func (se SErr) FieldsSlice() []string {
	return se.fields
}
