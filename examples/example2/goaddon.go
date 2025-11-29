package main

import "C"
import (
	"encoding/json"
	"time"
)

func durationFromMS(ms int32) time.Duration {
	if ms < 0 {
		return 0
	}
	return time.Duration(ms) * time.Millisecond
}

//export AddNumbers
func AddNumbers(a, b int32) C.int {
	return C.int(a + b)
}

//export SyncDelay
func SyncDelay(ms int32) bool {
	time.Sleep(durationFromMS(ms))
	return true
}

//export AsyncDelay
func AsyncDelay(ms int32, cb *C.char) bool {
	_ = cb // callback name used by gonode for async exports
	time.Sleep(durationFromMS(ms))
	return true
}

//export RemoveItemFromArray
func RemoveItemFromArray(arrJSON *C.char, idx int32) *C.char {
	var items []interface{}
	if err := json.Unmarshal([]byte(C.GoString(arrJSON)), &items); err != nil {
		return C.CString("[]")
	}

	i := int(idx)
	if i >= 0 && i < len(items) {
		items = append(items[:i], items[i+1:]...)
	}

	result, err := json.Marshal(items)
	if err != nil {
		return C.CString("[]")
	}

	return C.CString(string(result))
}

//export RemoveKeyFromObject
func RemoveKeyFromObject(objJSON, key *C.char) *C.char {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(C.GoString(objJSON)), &obj); err != nil {
		obj = make(map[string]interface{})
	}

	delete(obj, C.GoString(key))

	result, err := json.Marshal(obj)
	if err != nil {
		return C.CString("{}")
	}

	return C.CString(string(result))
}

func returnMessageAfterTwoSeconds(message string) string {
	time.Sleep(2 * time.Second)
	return message
}

//export SyncReturnMessageAfterTwoSeconds
func SyncReturnMessageAfterTwoSeconds(message *C.char) *C.char {
	output := returnMessageAfterTwoSeconds(C.GoString(message))
	return C.CString(output)
}

//export AsyncReturnMessageAfterTwoSeconds
func AsyncReturnMessageAfterTwoSeconds(message, cb *C.char) *C.char {
	_ = cb // callback name used by gonode for async exports
	output := returnMessageAfterTwoSeconds(C.GoString(message))
	return C.CString(output)
}

func main() {}
