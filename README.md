# Do not use this in production. I recommend you use [Koffi](https://github.com/Koromix/koffi) to develop a node bindings for a Go project.

<div align="center">
<img width="120" style="padding-top: 50px" src="http://47.104.180.148/gonode/gonode_logo.svg"/>
<h1 style="margin: 0; padding: 0">Gonode</h1>
<p>This is a development tool that can quickly use Golang to develop NodeJS Addon.</p>
<a href="https://goreportcard.com/report/github.com/VastBlast/gonode"><img src="https://goreportcard.com/badge/github.com/VastBlast/gonode"/></a>
<a href="https://pkg.go.dev/github.com/VastBlast/gonode"><img src="https://pkg.go.dev/badge/github.com/VastBlast/gonode"/></a>
<a href="https://github.com/VastBlast/gonode/releases"><img src="https://img.shields.io/github/v/release/VastBlast/gonode.svg"/></a>
<a href="https://github.com/VastBlast/gonode/blob/master/LICENSE"><img src="https://img.shields.io/github/license/VastBlast/gonode.svg"/></a>
<a href="https://github.com/VastBlast/gonode"><img src="https://img.shields.io/github/stars/VastBlast/gonode.svg"/></a>
<a href="https://github.com/VastBlast/gonode"><img src="https://img.shields.io/github/last-commit/VastBlast/gonode.svg"/></a>
</div>

<br/>


<a href="https://github.com/VastBlast/gonode">GONODE</a> is a development tool that quickly uses Golang to develop <b>NodeJS Addon</b>. You only need to concentrate on the development of Golang, and you don't need to care about the implementation of the bridge layer. It supports JavaScript sync calls and async callbacks.

> Originally created as <a href="https://github.com/wenlng/gonacli">gonacli</a> by <a href="https://github.com/wenlng">wenlng</a> — many thanks for the foundation.

<br/>

 ⭐️ If it helps you, please give a star.

- [https://github.com/VastBlast/gonode](https://github.com/VastBlast/gonode)


<br/>


## Compatible Support Of Gonode
- Linux
- Mac OS
- Windows

## Compatible Support Of NodeJS Addon
- Linux / Mac OS / Windows
- NodeJS(12.0+)
- Npm(6.0+)
- Node-gyp(9.0+)
- Go(1.14+)

## Use Golang Install
<p>Ensure that the system is configured with GOPATH environment variables before installation</p>

Linux or Mac OS
``` shell
# .bash_profile
export GOPATH="/Users/awen/go"
# set bin dir
export PATH="$PATH:$GOPATH:$GOPATH/bin"
```

Windows
``` shell
# set system path
GOPATH: C:\awen\go
# set bin dir
PATH: %GOPATH%\bin
``` 

Install
``` shell
$ go install github.com/VastBlast/gonode@latest

$ gonode version
```
<br/>


## Compilation Of Windows OS Environment
In the Windows OS environment, you need to install the `gcc/g++` compiler support required by Go CGO, download the `MinGW` installation, configure the `PATH` environment variable of `MinGW/bin`, and execute `gcc` normally on the command line.

``` shell
$ gcc -v
```

When compiling Node Addon in the Windows OS environment, you also need to install the build tool that `node-gyp` and depends on.

``` shell
$ npm install --global node-gyp

$ npm install --global --production windows-build-tools
```
<br/>

## Gonode Command

### 1. generate

Generate bridge code related to NodeJS Addon according to the configuration of goaddon

``` shell
# By default, it reads the goaddon in the current directory Json configuration file
$ gonode generate

# --config: Specify Profile
$ gonode generate --config demoaddon.json
```

### 2. build

Same as the `go build - buildmode=c-archive` command, compile the library

``` shell
# Compile to generate library
$ gonode build

# --args: Specify the args of go build
# --config: Specify Profile
$ gonode build --args '-ldflags "-s -w"'
```

### 3. make

Same as the `node-gyp configure && node-gyp build` command，Compile NodeJS Addon

``` text
Please ensure that the node gyp compiler has been installed on the system before using the "make" command

Before using the "--npm-i" arg, ensure that the system has installed the npm package dependency management tool
```

``` shell
# --args: Specify the parameters of node-gyp build，for example "--debug"
$ gonode make --args '--debug'
```

### 4. clean

Clean the configured output directory when it exists.

``` shell
$ gonode clean

# --config: Specify Profile
$ gonode clean --config demoaddon.json
```

### 5. all

Run the full pipeline: clean -> generate -> build.

``` shell
$ gonode all

# --config: Specify Profile
# --build-args: Args forwarded to the Go build step
$ gonode all --config demoaddon.json --build-args '-ldflags "-s -w"'
```

<br/>

## Use Golang to develop a Demo of NodeJS Addon

Tip: Ensure that relevant commands can be used normally. This is a demo under Linux/OS environment.

``` shell
# go
$ go version

# node
$ node -v

# npm
$ npm -v

# node-gyp
$ node-gyp -v
```


#### 1. Create Goaddon Configure File

`/goaddon.json`

``` json
{
  "name": "demoaddon",
  "sources": [
    "demoaddon.go"
  ],
  "output": "./demoaddon/",
  "exports": [
    {
      "name": "Hello",
      "args": [
        {
          "name": "name",
          "type": "string"
        }
      ],
      "returntype": "string",
      "jscallname": "hello",
      "jscallmode": "sync"
    }
  ]
}
```

#### 2. Write Golang Code

`/demoaddon.go`

gonode will auto-generate a temporary `temp_gonode_helpers.go` (with `FreeCString`) if you don't define one. If you prefer to own the implementation, add the export below (no goaddon.json entry needed).

``` go
package main

// #include <stdlib.h>
import "C"
import "unsafe"

// notice：//export xxxx is necessary

//export FreeCString
func FreeCString(str *C.char) {
	C.free(unsafe.Pointer(str))
}

//export Hello
func Hello(_name *C.char) *C.char {
	// args string type，return string type
	name := C.GoString(_name)
	
	res := "hello"
	if len(name) > 0 {
	    res += "," + name
	}
	
	return C.CString(res)
}
```

#### 3. Generate Bridging Napi C/C++ Code
``` shell
# Save to the "./demoaddon/" directory
$ gonode generate --config ./goaddon.json
```

#### 4.Compile Libraries
``` shell
# Save to the "./demoaddon/" directory
$ gonode build
```


#### 5. Compile Nodejs Addon
``` shell
# Save to the "./demoaddon/build" directory
$ gonode make
```

#### 6. Create JS Test File

`/demoaddon/test.js`

``` javascript
const demoaddon = require('.')

const name = "awen"
const res = demoaddon.hello(name)
console.log('>>> ', res)

```

``` shell
$ node ./test.js
# >>> hello, awen
```

<br/>

## Configure File Description
``` text
{
  "name": "demoaddon",      // Name of Nodejs Addon
  "sources": [              // File list of go build，Cannot have path
    "demoaddon.go"
  ],
  "output": "./demoaddon/", // Output directory path
  "exports": [              // Exported interface, generating the Napi and C/C++ code of Addon
    {
      "name": "Hello",      // The name of the "//export Hello" interface corresponding to Golang must be consistent
      "args": [             // The parameter type of the passed parameter list must be consistent with the type table
        {                  
          "name": "name",
          "type": "string"
        }
      ],
      "returntype": "string",   // The type returned to JavaScript，has no callback type
      "jscallname": "hello",    // JavaScript call name
      "jscallmode": "sync"      // Sync is synchronous execution, and Async is asynchronous execution
    },
    {
        name: "xxx",
        ....
    }
  ]
}
```

## Type Table

|    Type     | Golang Args | Golang Return  |   JS / TS   |
|:-----------:|:-----------:|:--------------:|:-----------:|
|     int     |    int32    |     C.int      |   number    |
|    int32    |    int32    |     C.int      |   number    |
|    int64    |    int64    |   C.longlong   |   number    |
|   uint32    |   uint32    |     C.uint     |   number    |
|    float    |   float32   |    C.float     |   number    |
|   double    |   float64   |    C.double    |   number    |
|   boolean   |    bool     |      bool      |   boolean   |
|   string    |   *C.char   |    *C.char     |   string    |
|    array    |   *C.char   |    *C.char     |    Array    |
|   object    |   *C.char   |    *C.char     |   Object    |
|  callback   |   *C.char   |       -        |  Function   |

### The returntype field type of the configuration file
The `returntype` field has no callback type

### array type
When there are multiple levels when returning, it is not recommended to use in the `returntype`


1. The `array` type received in Golang is a string `*C.Char` type, which needs to be use `[]interface{}` and `json.Unmarshal`


2. The `array` type is when Golang returns `*C.Char` type, use `json.Marshal`


3. The `array` type is an Array type when JavaScript is passed, but currently only supports one layer when receiving. Please use string method to return multiple layers in Golang, and then use JavaScript's `JSON.parse`

### object type
When there are multiple levels when returning, it is not recommended to use in the `returntype`

1. The `object` type received in Golang is a string type. You need to use `[string]interface{}` and `json.Unmarshal`


2. The `object` type is when Golang returns `*C.Char` type, use `json.Marshal`


3. The `object` type is an Object type when JavaScript is passed, but currently only supports one layer when receiving. Please use string method to return multiple layers in Golang, and then use JavaScript's `JSON.parse`

<br/>

## JavaScript Sync Call

`/goaddon.json`

``` json
{
  "name": "demoaddon",
  "sources": [
    "demoaddon.go"
  ],
  "output": "./demoaddon/",
  "exports": [
    {
      "name": "Hello",
      "args": [
        {
          "name": "name",
          "type": "string"
        }
      ],
      "returntype": "string",
      "jscallname": "hello",
      "jscallmode": "sync"
    }
  ]
}
```

#### 2. Golang Code

`/demoaddon.go`

If you skip defining `FreeCString`, gonode will inject a temporary helper for you; add this export only if you want to supply your own implementation.

``` go
package main

// #include <stdlib.h>
import "C"
import (
	"time"
	"unsafe"
)

//export FreeCString
func FreeCString(str *C.char) {
	C.free(unsafe.Pointer(str))
}

//export Hello
func Hello(_name *C.char) *C.char {
	// args is string type，return string type
	name := C.GoString(_name)

	res := "hello"
	ch := make(chan bool)

	go func() {
		// Time consuming task processing
		time.Sleep(time.Duration(2) * time.Second)
		if len(name) > 0 {
			res += "," + name
		}
		ch <- true
	}()

	<-ch

	return C.CString(res)
}
```

#### 3. Test

`/test.js`

``` javascript
const demoaddon = require('./demoaddon')

const name = "awen"
const res = demoaddon.hello(name)
console.log('>>> ', res)
```

<br/>

## JavaScript Async Call

`/goaddon.json`

``` json
{
  "name": "demoaddon",
  "sources": [
    "demoaddon.go"
  ],
  "output": "./demoaddon/",
  "exports": [
    {
      "name": "Hello",
      "args": [
        {
          "name": "name",
          "type": "string"
        },
        {
          "name": "cbs",
          "type": "callback"
        }
      ],
      "returntype": "string",
      "jscallname": "hello",
      "jscallmode": "async"
    }
  ]
}
```

#### 2. Golang Code

`/demoaddon.go`

You can omit `FreeCString` and let gonode inject the helper automatically, or add the export yourself if you prefer explicit control.

``` go
package main

// #include <stdlib.h>
import "C"
import (
	"time"
	"unsafe"
)

//export FreeCString
func FreeCString(str *C.char) {
	C.free(unsafe.Pointer(str))
}

//export Hello
func Hello(_name *C.char, cbsFnName *C.char) *C.char {
	// args is string type，return string type
	name := C.GoString(_name)

	res := "hello"
	ch := make(chan bool)

	go func() {
		// Time consuming task processing
		time.Sleep(time.Duration(2) * time.Second)
		if len(name) > 0 {
			res += "," + name
		}
		ch <- true
	}()

	<-ch

	return C.CString(res)
}
```

#### 3. Test

`/test.js`

``` javascript
const demoaddon = require('./demoaddon')

const name = "awen"
demoaddon.hello(name, function(res){
    console.log('>>> ', res)
})
```

<br/>

## LICENSE
    MIT
