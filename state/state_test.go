package state

import (
	"encoding/gob"
	"fmt"
	"testing"
)

func MapItemTypeToString(m mapItemType) string {
	switch m {
	case MapItemStr:
		return "string"
	case MapItemInt:
		return "int"
	case MapItemGameState:
		return "GameState"
	case MapItemUnsupported:
		return "unsupported"
	}

	return "unaccounted for type... shouldn't be here"
}

func MapItemToString(a any) string {
	switch GetItemMapType(a) {
	case MapItemStr:
		return "'" + a.(string) + "'"
	case MapItemInt:
		return fmt.Sprintf("%v", a.(int))
	case MapItemGameState:
		return MapToString(a.(GameState))
	case MapItemUnsupported:
		return "Error! Unsupported type!"
	}

	return "unaccounted for item type... shouldn't be here"
}

func MapToString(m GameState) string {
	var res string
	for key, val := range m.Elements {
		switch GetItemMapType(val) {
		case MapItemStr:
			res += fmt.Sprintf("Key: '%v', Val: '%v'\n", key, val.(string))
		case MapItemInt:
			res += fmt.Sprintf("Key: '%v', Val: %v\n", key, val.(int))
		case MapItemGameState:
			res += fmt.Sprintf("Key: %v, Val: (nested map)\n\t", key)
			res += MapToString(val.(GameState))
		case MapItemUnsupported:
			res += fmt.Sprintf("Key: '%v', Val: Error! Unsupported type!", key)
		}
	}

	return res
}

func MapDiff(m1 GameState, m2 GameState) string {
	diff := ""
	if len(m1.Elements) != len(m2.Elements) {
		diff += fmt.Sprintf("len(m1) == %v vs len(m2) == %v\n", len(m1.Elements), len(m2.Elements))
		diff += "Map 1:\n" + MapToString(m1) + "Map 2:\n" + MapToString(m2)
		return diff
	}

	for key, val1 := range m1.Elements {
		val2, exists := m2.Elements[key]
		if !exists {
			diff += fmt.Sprintf("Key: '%v' not in m2!", key)
			continue
		}

		type1 := GetItemMapType(val1)
		type2 := GetItemMapType(val2)

		if type1 != type2 {
			diff += fmt.Sprintf("Mismatched types! type1: %v vs type2 %v",
				MapItemTypeToString(type1), MapItemTypeToString(type2))
		} else if type1 == MapItemGameState {
			diff += MapDiff(val1.(GameState), val2.(GameState))
		} else if val1 != val2 {
			diff += fmt.Sprintf("m1[%v] != m2[%v]! %v vs %v!", key, key, val1, val2)
		}
	}

	for key := range m2.Elements {
		_, exists := m1.Elements[key]
		if !exists {
			diff += fmt.Sprintf("Key: '%v' not in m1!", key)
		}
	}

	return diff
}

func TestDeserializeEmpty(t *testing.T) {
	gob.Register(GameState{})
	want := GameState{map[string]any{}}
	mapDiff := MapDiff(Deserialize(Serialize(want)), want)
	if len(mapDiff) > 0 {
		t.Errorf(mapDiff)
	}
}

func TestDeserializeSuit(t *testing.T) {
	gob.Register(GameState{})
	want := GameState{map[string]any{"suit": 1}}
	mapDiff := MapDiff(Deserialize(Serialize(want)), want)
	if len(mapDiff) > 0 {
		t.Errorf(mapDiff)
	}
}

func TestDeserializeRecursive(t *testing.T) {
	gob.Register(GameState{})
	want := GameState{map[string]any{"suit": GameState{map[string]any{"bruh": 1}}}}
	mapDiff := MapDiff(Deserialize(Serialize(want)), want)
	if len(mapDiff) > 0 {
		t.Errorf(mapDiff)
	}
}

func TestDeserializeManyNested(t *testing.T) {
	gob.Register(GameState{})
	want := GameState{map[string]any{
		"many":  GameState{map[string]any{"suit": GameState{map[string]any{"bruh": 1}}}},
		"many1": GameState{map[string]any{"suit": GameState{map[string]any{"bruh": 2}}}},
	}}

	mapDiff := MapDiff(Deserialize(Serialize(want)), want)
	if len(mapDiff) > 0 {
		t.Errorf(mapDiff)
	}
}
