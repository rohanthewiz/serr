package serr

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
	return se.err.Error()
}

// Return all SErr attributes as a map of string keys and values
func (se SErr) FieldsMap() map[string]string {
	flds := map[string]string{}
	key := ""
	for i, str := range se.fields {
		if i % 2 == 0 {  // even indices are presumed to be keys
			key = str
		} else {
			if orig, ok := flds[key]; ok {  // we've seen this key before
				flds[key] = str + " - " + orig
			} else {
				flds[key] = str
			}
		}
	}
	return flds
}

// Wrap an existing error. Attribute keys and values must be strings.
// This requires an even number of fields unless a single field is given in which case it is added under the key "msg".
// Returns an SErr (structured err)
func Wrap(err error, fields ...string) error {
	var flds []string

	if se, ok := err.(SErr); ok && len(se.fields) > 0 {
		flds = append(flds, se.fields...)  // add existing fields first
	}

	if len(fields) == 1 {
		flds = append(flds, []string{"msg", fields[0]}...)
	} else {
		flds = append(flds, fields...)
	}
	return SErr{err, flds}  // return
}

// Return the wrapped error
func (se SErr) OriginalErr() error {
	return se.err
}

// Return the internal list of keys and values
func (se SErr) FieldsSlice() []string {
	return se.fields
}
