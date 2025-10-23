package serr

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestSErrFromatting(t *testing.T) {
	ser := NewSErr("my error", "att1", "val1", "att2", "val2")
	ser2 := WrapAsSErr(ser, "att2", "valNew")
	const expected = "att2[val2 -> valNew]"

	result := ser2.FieldsAsCustomString(", ", " -> ")
	if !strings.Contains(result, expected) {
		t.Errorf("Expected result to contain '%s', got '%s'", expected, result)
	}
}

func TestSErr(t *testing.T) {
	const strErr1 = "This is a test err"
	const thisIsMyMessage = "This is my message"

	// We should safely ignore a nil err
	ret := Wrap(nil, "We should be able to handle a nil error without crashing")
	if ret != nil {
		t.Error("We should return a nil when a nil error is wrapped", "got:", ret)
	}

	ser := New(strErr1, "thing1", "thing1val", "thing2", "thing2val")
	if ser.Error() != strErr1 {
		t.Errorf("Expected custom error to contain '%s', got '%s'", strErr1, ser.Error())
		t.FailNow()
	}
	if _, ok := ser.(SErr); !ok {
		t.Error("ser should be a SErr")
		t.FailNow()
	}

	// Add some fields to an existing sErr
	err := Wrap(ser, "thing2", "thing2NewVal")
	se, ok := err.(SErr)
	if !ok {
		t.Error("Wrap should return an error containing a concrete SErr type")
	} else {
		t.Log("se =>", se.String())

		// Test SErr#GetError
		strErr := se.GetError().Error()
		if strErr != strErr1 {
			t.Errorf(`Expected wrapped error string to be "%s", got "%s"`, strErr1, strErr)
		}

		// Test SErr#Fields
		strFlds := se.Fields()
		// fmt.Printf("[Debug] strFlds: %#v; Immediate location: %s\n", strFlds, FunctionLoc(FuncLevel1)) // debug
		if len(strFlds) != 14 {
			t.Error("Expected length of SErr.Fields() to be 14, got", len(strFlds))
		}

		// Test SErr#FieldsMap
		mapFlds := se.FieldsMap()
		if len(mapFlds) != 4 {
			t.Error("Expected length of SErr.MapFlds() to be 4, got", len(mapFlds))
			t.FailNow()
		}

		if val, ok := mapFlds["thing1"]; ok {
			if val != "thing1val" {
				t.Errorf("Expected thing1 to be 'thing1Val', got '%s'", val)
			}
		} else {
			t.Error("mapFlds should contain the key: thing1")
		}

		if val, ok := mapFlds["thing2"]; ok {
			if val != "thing2NewVal - thing2val" {
				t.Errorf("Expected thing2 to be 'thing2NewVal - thing2val', got '%s'", val)
			}
		} else {
			t.Error("mapFlds should contain the key: thing2")
		}
	}

	// We should be able to wrap with a single field which becomes `"msg": field`

	er := Wrap(errors.New(strErr1), thisIsMyMessage)
	se, ok = er.(SErr)
	if !ok {
		t.Error("er should be a SErr")
		t.FailNow()
	}

	sl := se.Fields()
	if len(sl) != 6 {
		t.Error("Structured error from an error wrapped with a single field should contain 6 fields, got", len(sl))
	}
	if len(sl) > 0 && sl[0] != "msg" {
		t.Error("Structured error from an error wrapped with a single field should have 'map' as the first field")
	}
	if len(sl) > 1 && sl[1] != thisIsMyMessage {
		t.Errorf(`The structured error should have "%s" as the second field, got "%s"`, thisIsMyMessage, sl[1])
		fmt.Println()
	}

	// Test WrapAsSErr
	sr := WrapAsSErr(errors.New(strErr1), thisIsMyMessage, UserMsgKey, "Your account balance is very low", UserMsgSeverityKey, Severity.Warn)
	sl = sr.Fields()
	if len(sl) != 10 {
		t.Error("Structured error from an error wrapped with a single field should contain 10 fields, got", len(sl))
	}

	newSr := NewSErr(strErr1, thisIsMyMessage)
	nsf := newSr.Fields()
	if len(nsf) != 6 {
		t.Error(fmt.Sprintf("Structured error from an error wrapped with a single field should contain %d fields, got %d", 6, len(nsf)))
	}
}

func TestF(t *testing.T) {
	const expectedErr = "test error: 42"

	err := F("test error: %d", 42)

	// Test that it returns a proper error
	if err == nil {
		t.Error("F() returned nil instead of an error")
	}

	// Test that the error message is correctly formatted
	if err.Error() != expectedErr {
		t.Errorf("Expected error message '%s', got '%s'", expectedErr, err.Error())
	}

	// Test that it returns a SErr type
	se, ok := err.(SErr)
	if !ok {
		t.Error("F() should return a SErr type")
		return
	}

	// Test that location and function fields were added
	fields := se.FieldsMap()
	if _, ok := fields["location"]; !ok {
		t.Error("Expected 'location' field to be present")
	}

	if _, ok := fields["function"]; !ok {
		t.Error("Expected 'function' field to be present")
	}

	// Test the string representation
	strErr := se.String()
	if !strings.Contains(strErr, expectedErr) {
		t.Errorf("String() should contain '%s', got '%s'", expectedErr, strErr)
	}
}

func TestGetAttribute(t *testing.T) {
	// Build a SErr without auto-added context so we can control attributes precisely
	se := NewSerrNoContext(errors.New("base error"))
	se.AppendAttributes("k1", 123, "k2", "v2")

	// Existing key with non-string value
	if val, ok := se.GetAttribute("k1"); !ok {
		t.Fatalf("Expected to find attribute 'k1'")
	} else {
		iv, ok := val.(int)
		if !ok || iv != 123 {
			t.Fatalf("Expected 'k1' to be int 123, got %#v", val)
		}
	}

	// Missing key should return (nil, false)
	if val, ok := se.GetAttribute("missing"); ok || val != nil {
		t.Fatalf("Expected missing attribute to return (nil,false), got (%#v,%v)", val, ok)
	}

	// Duplicate key should concatenate with newest value first as a string
	se.AppendAttributes("k1", "next")
	if val, ok := se.GetAttribute("k1"); !ok {
		t.Fatalf("Expected to find attribute 'k1' after duplicate add")
	} else {
		str, ok := val.(string)
		if !ok {
			t.Fatalf("Expected 'k1' to be a string after duplicate add, got %#v", val)
		}
		expected := "next - 123"
		if str != expected {
			t.Fatalf("Expected concatenated value %q, got %q", expected, str)
		}
	}
}
