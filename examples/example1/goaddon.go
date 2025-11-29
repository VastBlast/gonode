package main

import "C"
import (
	"encoding/json"
	"fmt"
	"time"
)

//export IntSum32
func IntSum32(x, y int32) C.int {
	// Accept int32 and return int32
	return C.int(x + y)
}

//export IntSum64
func IntSum64(x, y int64) C.longlong {
	// Accept int64 and return int64
	return C.longlong(x + y)
}

//export UintSum32
func UintSum32(x, y uint32) C.uint {
	// Accept uint32 and return uint32
	return C.uint(x + y)
}

//export CompareInt
func CompareInt(x, y int32) bool {
	// Accept int32 and return boolean
	return x > y
}

//export FloatSum
func FloatSum(x, y float32) C.float {
	// Accept float and return float
	return C.float(x + y)
}

//export DoubleSum
func DoubleSum(x, y float64) C.double {
	// Accept double and return double
	return C.double(x + y)
}

//export FormatStr
func FormatStr(s *C.char) *C.char {
	// Accept string and return string
	ss := C.GoString(s)
	return C.CString("golang out >>> " + ss)
}

//export EmptyString
func EmptyString(s *C.char) bool {
	// Accept string and return boolean
	ss := C.GoString(s)
	return len(ss) <= 0
}

//export FilterMap
func FilterMap(s *C.char) *C.char {
	// Accept object and return object
	ss := C.GoString(s)
	fmt.Println(fmt.Sprintf("golang out >>> %s", string(ss)))

	var m2 = make(map[string]interface{}, 0)
	err2 := json.Unmarshal([]byte(string(ss)), &m2)
	if err2 != nil {
		fmt.Println("golang out >>> err: ", err2.Error())
	}
	fmt.Println("golang out >>> map len: ", m2, len(m2))

	var m = make(map[string]string, 2)
	m["a"] = "aaaaa"
	m["b"] = "bbbbb"
	jsonStr, err := json.Marshal(m)
	if err != nil {
		fmt.Println("golang out >>> err: ", err.Error())
	}

	result := string(jsonStr)
	return C.CString(result)
}

//export CountMap
func CountMap(s *C.char) C.int {
	// Accept object and return int32
	ss := C.GoString(s)
	fmt.Println(fmt.Sprintf("golang out >>> %s", string(ss)))

	var m2 = make(map[string]interface{}, 0)
	err2 := json.Unmarshal([]byte(string(ss)), &m2)
	if err2 != nil {
		fmt.Println("golang out >>>err: ", err2.Error())
		return 0
	}

	return C.int(len(m2))
}

//export IsMapType
func IsMapType(s *C.char) bool {
	// Accept object and return boolean
	ss := C.GoString(s)
	fmt.Println(fmt.Sprintf("golang out >>> %s", string(ss)))

	var m2 = make(map[string]interface{}, 0)
	err2 := json.Unmarshal([]byte(string(ss)), &m2)
	if err2 != nil {
		fmt.Println("golang out >>>err: ", err2.Error())
		return false
	}
	return true
}

//export FilterSlice
func FilterSlice(s *C.char) *C.char {
	// Accept array and return array
	ss := C.GoString(s)
	fmt.Println("golang out >>> slice len: ", ss, len(ss))

	var m2 = make([]interface{}, 0)
	err2 := json.Unmarshal([]byte(string(ss)), &m2)
	if err2 != nil {
		fmt.Println("golang out >>>err: ", err2.Error())
	}
	fmt.Println("golang out >>> slice len: ", m2, len(m2))

	var m = make([]interface{}, 2)
	m[0] = "hello"
	m[1] = "world"

	jsonStr, err := json.Marshal(m)
	if err != nil {
		fmt.Println("golang out >>>err: ", err.Error())
	}

	result := string(jsonStr)

	return C.CString(result)
}

//export CountSlice
func CountSlice(s *C.char) C.int {
	// Accept array and return int32
	ss := C.GoString(s)
	fmt.Println("golang out >>> slice len: ", ss, len(ss))

	var m2 = make([]interface{}, 0)
	err2 := json.Unmarshal([]byte(string(ss)), &m2)
	if err2 != nil {
		fmt.Println("golang out >>>err: ", err2.Error())
	}
	fmt.Println("golang out >>> slice len: ", m2, len(m2))

	return C.int(len(m2))
}

//export IsSliceType
func IsSliceType(s *C.char) bool {
	// Accept object and return boolean
	ss := C.GoString(s)
	fmt.Println(fmt.Sprintf("golang out >>> %s", string(ss)))

	var m2 = make([]interface{}, 0)
	err2 := json.Unmarshal([]byte(string(ss)), &m2)
	if err2 != nil {
		fmt.Println("golang out >>>err: ", err2.Error())
		return false
	}
	return true
}

type CallbackOutput struct {
	Data   string `json:"data"`
	Output string `json:"output"`
}

var callbackCount = 0

// ===========================
// Synchronous execution

//export  SyncCallbackReStr
func SyncCallbackReStr(arg *C.char) *C.char {
	// Runs synchronously and blocks the main thread
	// Accept string and return string
	result := ""
	ch := make(chan bool)

	go func() {
		callbackCount++
		curCount := callbackCount

		fmt.Println("golang out >>> run", curCount, C.GoString(arg))
		time.Sleep(time.Duration(2) * time.Second)
		var co CallbackOutput
		co.Data = "hello wait return hello"

		if curCount%2 == 1 {
			co.Data = "hello wait return world"
		}

		co.Output = fmt.Sprintf("%d", curCount)
		jsonStr, err := json.Marshal(co)
		if err != nil {
			fmt.Println("golang out >>>err: ", err.Error())
		}
		result = string(jsonStr)
		ch <- true
	}()

	<-ch
	return C.CString(result)
}

//export SyncCallbackReArr
func SyncCallbackReArr(arg *C.char) *C.char {
	// Runs synchronously and blocks the main thread
	// Accept array and return array
	result := ""
	ch := make(chan bool)
	ss := C.GoString(arg)

	go func() {
		callbackCount++
		curCount := callbackCount

		fmt.Println("golang out >>>> run", curCount, ss)

		var m2 = make([]interface{}, 0)
		err2 := json.Unmarshal([]byte(string(ss)), &m2)
		if err2 != nil {
			fmt.Println("err: ", err2.Error())
		}
		fmt.Println("golang out >>> slice len: ", m2, len(m2))

		var m = make([]interface{}, 3)
		m[0] = "hello"
		m[1] = "world"
		m[2] = curCount

		jsonStr, err := json.Marshal(m)
		if err != nil {
			fmt.Println("golang out >>>err: ", err.Error())
		}

		result = string(jsonStr)
		ch <- true
	}()

	<-ch
	return C.CString(result)
}

//export SyncCallbackReObject
func SyncCallbackReObject(arg *C.char) *C.char {
	// Runs synchronously and blocks the main thread
	// Accept object and return object
	result := ""
	ch := make(chan bool)
	ss := C.GoString(arg)

	go func() {
		callbackCount++
		curCount := callbackCount

		fmt.Println("golang out >>> run", curCount, C.GoString(arg))

		var m2 = make(map[string]interface{}, 0)
		err2 := json.Unmarshal([]byte(string(ss)), &m2)
		if err2 != nil {
			fmt.Println("err: ", err2.Error())
		}
		fmt.Println("golang out >>> map len: ", m2, len(m2))

		var m = make(map[string]interface{}, 3)
		m["k1"] = "hello"
		m["k2"] = "world"
		m["k3"] = curCount

		jsonStr, err := json.Marshal(m)
		if err != nil {
			fmt.Println("golang out >>>err: ", err.Error())
		}

		result = string(jsonStr)
		ch <- true
	}()

	<-ch
	return C.CString(result)
}

//export SyncCallbackReCount
func SyncCallbackReCount(arg *C.char) C.int {
	// Runs synchronously and blocks the main thread
	// Accept string and return int32
	result := 0
	ch := make(chan bool)

	go func() {
		callbackCount++
		curCount := callbackCount

		fmt.Println("golang out >>> run", curCount, C.GoString(arg))

		result = curCount
		ch <- true
	}()

	<-ch
	return C.int(result)
}

//export SyncCallbackReBool
func SyncCallbackReBool(arg *C.char) bool {
	// Runs synchronously and blocks the main thread
	// Accept string and return boolean
	result := false
	ch := make(chan bool)

	go func() {
		callbackCount++
		curCount := callbackCount

		fmt.Println("golang out >>> run", curCount, C.GoString(arg))

		result = true
		ch <- true
	}()

	<-ch
	return result
}

//export SyncCallbackSleep
func SyncCallbackSleep(t int32) bool {
	// Runs synchronously and blocks the main thread
	// Accept int32 and return boolean
	ch := make(chan bool)

	go func() {
		callbackCount++
		curCount := callbackCount
		d := t
		time.Sleep(time.Duration(d) * time.Second)
		fmt.Println("golang out >>> run", curCount, d)
		ch <- true
	}()

	<-ch
	return true
}

// =========== Async

//export  ASyncCallbackReStr
func ASyncCallbackReStr(arg *C.char, cbFuncStr *C.char) *C.char {
	// Runs asynchronously without blocking the main thread
	// Accepts string and string, returns string
	return SyncCallbackReStr(arg)
}

//export ASyncCallbackReIntSum64
func ASyncCallbackReIntSum64(x, y int64, cbFuncStr *C.char) C.longlong {
	// Accept int64 and return int64
	SyncCallbackSleep(1)
	return C.longlong(x + y)
}

//export ASyncCallbackReUintSum32
func ASyncCallbackReUintSum32(x, y uint32, cbFuncStr *C.char) C.uint {
	// Accept uint32 and return uint32
	SyncCallbackSleep(1)
	return C.uint(x + y)
}

//export ASyncCallbackReArr
func ASyncCallbackReArr(arg *C.char, cbFuncStr *C.char) *C.char {
	// Runs asynchronously without blocking the main thread
	// Accepts string and string, returns array
	fmt.Println("golang out >>> cbFuncStr = ", C.GoString(cbFuncStr))
	return SyncCallbackReArr(arg)
}

//export ASyncCallbackReObject
func ASyncCallbackReObject(arg *C.char, cbFuncStr *C.char) *C.char {
	// Runs asynchronously without blocking the main thread
	// Accepts object and string, returns object
	fmt.Println("golang out >>> cbFuncStr = ", C.GoString(cbFuncStr))
	return SyncCallbackReObject(arg)
}

//export ASyncCallbackReCount
func ASyncCallbackReCount(arg *C.char, cbFuncStr *C.char) C.int {
	// Runs asynchronously without blocking the main thread
	// Accepts string and string, returns int32
	fmt.Println("golang out >>> cbFuncStr = ", C.GoString(cbFuncStr))
	return SyncCallbackReCount(arg)
}

//export ASyncCallbackReBool
func ASyncCallbackReBool(arg *C.char, cbFuncStr *C.char) bool {
	// Runs asynchronously without blocking the main thread
	// Accepts string and string, returns boolean
	fmt.Println("golang out >>> cbFuncStr = ", C.GoString(cbFuncStr))
	return SyncCallbackReBool(arg)
}

//export ASyncCallbackMArg
func ASyncCallbackMArg(arg *C.char, cbFuncStr *C.char, ext *C.char) bool {
	// Runs asynchronously without blocking the main thread
	// Accepts string and string, returns boolean
	fmt.Println("golang out >>> cbFuncStr = ", C.GoString(cbFuncStr))
	return SyncCallbackReBool(arg)
}

func main() {
	// ...
}
