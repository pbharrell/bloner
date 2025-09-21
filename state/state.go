package state

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type mapItemType uint8

const (
	MapItemStr mapItemType = iota
	MapItemInt
	MapItemGameState
	MapItemUnsupported
)

type GameState struct {
	Elements map[string]any
}

func GetItemMapType(a any) mapItemType {
	switch a.(type) {
	case string:
		return MapItemStr
	case int:
		return MapItemInt
	case GameState:
		return MapItemGameState
	default:
		return MapItemUnsupported
	}
}

// func WriteKey(m map[string]any, key string, val any) map[string]any {
// 	_, exists := m[key]
// 	if exists {
// 		fmt.Printf("Warning! Overwriting existing value at '%v'\n", key)
// 		fmt.Printf("old type: %T vs new type: %T", m[key], val)
// 	}
// 	m[key] = val
// 	return m
// }

// func ParseVal(val string) any {
// 	if len(val) == 0 {
// 		return nil
// 	}
//
// 	if len(val) >= 2 && '"' == val[0] && '"' == val[len(val)] {
// 		return val[1:]
// 	}
//
// 	intVal, err := strconv.Atoi(val)
// 	if err != nil {
// 		println("Failed to convert string to int! Returning nil")
// 		return nil
// 	}
//
// 	return intVal
// }

// func DeserializeElements(serialized string) map[string]any {
// 	elements := map[string]any{}
// 	element := ""
// 	parsingElement := false
// 	val := ""
// 	var parsedVal any
//
// 	parsingVal := false
// 	for i := 0; i < len(serialized); i++ {
// 		char := serialized[i]
// 		if char == '<' {
// 			if parsingVal {
// 				brackNum := 1
// 				j := i + 1
// 				for j < len(serialized) && brackNum > 0 {
// 					if serialized[j] == '>' {
// 						brackNum--
// 						continue
// 					}
// 					j++
// 				}
//
// 				parsedVal = DeserializeElements(serialized[i : j+1])
// 				i = j
// 			} else {
// 				parsingElement = true
// 			}
// 			continue
// 		}
//
// 		if char == '>' {
// 			if !parsingVal {
// 				println("Error! Should not reach the end of element without first parsing value.")
// 			}
//
// 			isMapType := GetItemMapType(parsedVal) == MapItemAnyMap
// 			if !isMapType {
// 				parsedVal = ParseVal(val)
// 			}
//
// 			elements = WriteKey(elements, element, parsedVal)
// 			return elements
//
// 		}
//
// 		if parsingElement {
// 			if char == ' ' {
// 				parsingElement = false
// 				parsingVal = true
// 				continue
// 			}
//
// 			element += string(char)
//
// 		} else if parsingVal {
// 			val += string(char)
// 		}
// 	}
//
// 	if len(element) > 0 {
// 		fmt.Printf("Should not be here! Somehow wound up with outstanding element: '%v'\n", element)
// 	}
//
// 	if len(val) > 0 {
// 		fmt.Printf("Should not be here! Somehow wound up with outstanding val: '%v'\n", val)
// 	}
//
// 	return elements
// }

func Serialize(s GameState) bytes.Buffer {
	var b bytes.Buffer

	enc := gob.NewEncoder(&b)
	if err := enc.Encode(s); err != nil {
		fmt.Println("Error encoding struct:", err)
	}

	return b
}

func Deserialize(b bytes.Buffer) GameState {
	var state GameState
	dec := gob.NewDecoder(&b)
	if err := dec.Decode(&state); err != nil {
		fmt.Println("Error decoding struct:", err)
	}

	return state
}
