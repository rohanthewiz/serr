package serr

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestWrapF(t *testing.T) {
	const baseErr = "base error"

	// nil error should return nil
	if got := WrapF(nil, "ignored %d", 1); got != nil {
		t.Error("WrapF(nil, ...) should return nil")
	}

	// non-nil error
	er := errors.New(baseErr)
	format := "number %d, word %s"
	out := WrapF(er, format, 7, "cats")
	if out == nil {
		t.Fatal("WrapF returned nil")
	}

	se, ok := out.(SErr)
	if !ok {
		t.Fatal("WrapF should return an error containing a concrete SErr type")
	}

	// original error preserved
	if se.GetError().Error() != baseErr {
		t.Errorf("Expected wrapped error '%s', got '%s'", baseErr, se.GetError().Error())
	}

	// fields length and ordering
	flds := se.Fields()
	if len(flds) != 6 {
		t.Errorf("Expected 6 fields, got %d", len(flds))
	}

	expectedMsg := fmt.Sprintf(format, 7, "cats")
	if len(flds) > 0 && flds[0] != "msg" {
		t.Error("First field key should be 'msg'")
	}
	if len(flds) > 1 && flds[1] != expectedMsg {
		t.Errorf("Expected message '%s', got '%s'", expectedMsg, flds[1])
	}

	// map contains msg, location, function
	mp := se.FieldsMap()
	if v, ok := mp["msg"]; !ok {
		t.Error("Expected 'msg' field to be present")
	} else if v != expectedMsg {
		t.Errorf("Expected 'msg' value '%s', got '%s'", expectedMsg, v)
	}

	if _, ok := mp["location"]; !ok {
		t.Error("Expected 'location' field to be present")
	}

	if _, ok := mp["function"]; !ok {
		t.Error("Expected 'function' field to be present")
	}

	// string representation contains base error and formatted message
	str := se.String()
	if !strings.Contains(str, baseErr) {
		t.Errorf("String() should contain '%s', got '%s'", baseErr, str)
	}
	if !strings.Contains(str, expectedMsg) {
		t.Errorf("String() should contain '%s', got '%s'", expectedMsg, str)
	}
}
