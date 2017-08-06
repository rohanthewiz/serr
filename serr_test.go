package serr

import (
	"testing"
	"fmt"
	"errors"
)

func TestSErr(t *testing.T) {
	const strErr1 = "Ok. This is a test err"
	// todo: constantize most test literals

	fmt.Println("Testing SErr")
	ser := NewSErr(
		errors.New(strErr1),
		[]string{"thing1", "thing1val", "thing2", "thing2val"},
	)
	if ser.Error() != strErr1 {
		t.Errorf("Expected custom error to contain '%s', got '%s'", strErr1, ser.Error())
		t.FailNow()
	}
	if _, ok := ser.(SErr); !ok {
		t.Error("ser should be a SErr")
		t.FailNow()
	}

	// Add some fields to an existing customErr
	err := Wrap(ser, "thing2", "thing2NewVal")
	se, ok := err.(SErr)
	if !ok {
		t.Error("Wrap should return an error containing a concrete SErr type")
	} else {
		// Test SErr#OriginalErr
		strErr := se.OriginalErr().Error()
		if strErr != strErr1 {
			t.Errorf(`Expected wrapped error string to be "%s", got "%s"`, strErr1, strErr)
		}
		// Test SErr#FieldsSlice
		strFlds := se.FieldsSlice()
		fmt.Println("[Debug] strFlds:", strFlds)  // debug
		if len(strFlds) != 6 {
			t.Error("Expected length of customErr.Fields() to be 6, got", len(strFlds))
		}
		// Test SErr#FieldsMap
		mapFlds := se.FieldsMap()
		if len(mapFlds) != 2 {
			t.Error("Expected length of customErr.MapFlds() to be 2, got", len(mapFlds))
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
	const thisIsMyMessage = "This is my message"
	er := Wrap(errors.New(strErr1), thisIsMyMessage)
	se, ok = er.(SErr)
	if !ok {
		t.Error("er should be a SErr"); t.FailNow()
	}
	sl := se.FieldsSlice()
	if len(sl) != 2 {
		t.Error("Structured error from an error wrapped with a single field should contain two fields")
	}
	if len(sl) > 0 && sl[0] != "msg" {
		t.Error("Structured error from an error wrapped with a single field should have 'map' as the first field")
	}
	if len(sl) > 1 && sl[1] != thisIsMyMessage {
		t.Errorf(`The structured error should have "%s" as the second field, got "%s"`, thisIsMyMessage, sl[1])
		fmt.Println()
	}
}
