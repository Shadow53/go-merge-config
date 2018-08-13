package merge

import (
	"reflect"
	"testing"
)

func TestCopyMap(t *testing.T) {
	m := map[string]interface{}{
		"slice": []interface{}{1, 2, 3, 4},
		"map": map[string]interface{}{
			"test": "testing"},
		"number": 5}

	m2 := CopyMap(m)
	if m2 == nil {
		t.Error("Copied map should not return nil")
	}

	if !reflect.DeepEqual(m, m2) {
		t.Errorf("Map and copy are not deeply equal:\n%#v\n%#v\n", m, m2)
	}

	m["number"] = 6
	num := m2["number"].(int)
	if num != 5 {
		t.Errorf("Expected number to be 5, was %d", num)
	}

	s := m["slice"].([]interface{})
	s = append(s, 9)
	m["slice"] = s
	s2 := m2["slice"].([]interface{})
	if len(s2) != 4 {
		t.Errorf(
			"slice expected to have four values, has %d: %#v", len(s2), s2)
	}

	mInner := m["map"].(map[string]interface{})
	mInner["test"] = "foobar"
	m2Inner := m2["map"].(map[string]interface{})
	if m2Inner["test"] == mInner["test"] {
		t.Error("Nested maps should be independent copies")
	} else if m2Inner["test"] != "testing" {
		t.Error("Copied inner map was not initialized correctly")
	}
}

func TestSliceCopy(t *testing.T) {
	s := []interface{}{1, "foo", []interface{}{5, 6, 7},
		map[string]interface{}{"test": "bar"}}
	s2 := CopySlice(s)
	if !reflect.DeepEqual(s, s2) {
		t.Errorf("Slice and copy are not deeply equal:\n%#v\n%#v", s, s2)
	}

	sInner := (s[2]).([]interface{})
	s2Inner := s2[2].([]interface{})
	if !reflect.DeepEqual(sInner, s2Inner) {
		t.Errorf(
			"Nested slices are not deeply equal:\n%#v\n%#v", sInner, s2Inner)
	}
	sInner = append(sInner, 5)
	// If appending to one changes the length of the other, they are the
	// same slice
	if len(sInner) == len(s2Inner) {
		t.Error("Nested slices are the same slice")
	}

	sMap := s[3].(map[string]interface{})
	s2Map := s2[3].(map[string]interface{})
	if !reflect.DeepEqual(sMap, s2Map) {
		t.Errorf("Nested maps are not deeply equal: \n%#v\n%#v", sMap, s2Map)
	}
	sMap["undefined"] = "not nil"
	// Item should be nil unless both are the same map
	if s2Map["undefined"] != nil {
		t.Errorf("Nested maps are the same map")
	}
}

func TestMerge(t *testing.T) {
	m := map[string]interface{}{
		"foo": "bar",
		"bar": 4,
		"baz": []interface{}{1, 2, 3},
		"nil": []interface{}{4, 5, 6},
		"map": map[string]interface{}{
			"key": "value",
			"baz": "foo"},
		"map2": map[string]interface{}{
			"key": "value"},
		"not-map": 4}
	m2 := map[string]interface{}{
		"foo": "baz",
		"bar": "baz",
		"baz": []interface{}{7, 8, 9},
		"map": map[string]interface{}{
			"key": "other",
			"foo": "baz"},
		"nil-map": map[string]interface{}{
			"key": "other",
			"foo": "baz"},
		"not-map": map[string]interface{}{
			"key": "other",
			"foo": "baz"}}
	merge := Merge(m, m2)
	if merge["foo"] != m2["foo"] {
		t.Errorf("Two like-type keys were not merged: %v, %v",
			merge["foo"], m2["foo"])
	}
	if merge["bar"] != m2["bar"] {
		t.Errorf("Two non-like-type keys were not merged: %v, %v",
			merge["bar"], m2["bar"])
	}
	if !reflect.DeepEqual(merge["baz"], m2["baz"]) {
		t.Errorf("Slice key should have been overwritten: %#v, %#v",
			merge["baz"], m2["baz"])
	}
	s := merge["baz"].([]interface{})
	merge["baz"] = append(s, 22)
	if reflect.DeepEqual(merge["baz"], m2["baz"]) {
		t.Errorf("Slice key should not be connected: %#v, %#v",
			merge["baz"], m2["baz"])
	}
	if !reflect.DeepEqual(merge["nil"], m["nil"]) {
		t.Errorf("Key not in second argument should not change: %#v, %#v",
			merge["nil"], m["nil"])
	}
	s = m["nil"].([]interface{})
	m["nil"] = append(s, 22)
	if reflect.DeepEqual(merge["nil"], m["nil"]) {
		t.Errorf("Unmerged argument should not be reference: %#v, %#v",
			merge["nil"], m["nil"])
	}
	if !reflect.DeepEqual(merge["map2"], m["map2"]) {
		t.Error("Map value not in second argument should not change")
	}
	mInner := m["map2"].(map[string]interface{})
	mInner["other"] = "something"
	if reflect.DeepEqual(merge["map2"], m["map2"]) {
		t.Errorf("Unmerged map value should not be reference: %#v, %#v",
			merge["nil"], m["nil"])
	}
	merged := map[string]interface{}{
		"key": "other",
		"foo": "baz",
		"baz": "foo"}
	if !reflect.DeepEqual(merged, merge["map"]) {
		t.Errorf("Map was not merged correctly: Wanted %#v, got %#v",
			merged, merge["map"])
	}

	// Map overwriting non-map non-nil value
	if !reflect.DeepEqual(merge["not-map"], m2["not-map"]) {
		t.Error("Map value in second argument should overwrite non-map")
	}
	mInner = m2["not-map"].(map[string]interface{})
	mInner["other"] = "something"
	if reflect.DeepEqual(merge["not-map"], m["not-map"]) {
		t.Errorf("Overwriting map value should not be reference: %#v, %#v",
			merge["not-map"], m["not-map"])
	}

	// Map overwriting nil value
	if !reflect.DeepEqual(merge["nil-map"], m2["nil-map"]) {
		t.Error("Map value in second argument should overwrite nil interface{}")
	}
	mInner = m2["nil-map"].(map[string]interface{})
	mInner["other"] = "something"
	if reflect.DeepEqual(merge["nil-map"], m["nil-map"]) {
		t.Errorf("Nil-overwriting map value should not be reference: %#v, %#v",
			merge["nil-map"], m["nil-map"])
	}
}
