package serr

import (
	"errors"
	"fmt"
	"testing"
)

func TestSErr(t *testing.T) {
	const strErr1 = "This is a test err"
	const thisIsMyMessage = "This is my message"

	// todo: constantize most test literals

	fmt.Println("Testing SErr")

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
	} /*else {
		t.Log("at this point ser is:", s.String())
	}*/

	// Add some fields to an existing sErr
	err := Wrap(ser, "thing2", "thing2NewVal")
	se, ok := err.(SErr)
	if !ok {
		t.Error("Wrap should return an error containing a concrete SErr type")
	} else {
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
