package async

import (
	"fmt"
	"github.com/VastBlast/gonode/config"
	"github.com/VastBlast/gonode/content/returns/reasync"
	"github.com/VastBlast/gonode/tools"
	"strings"
)

// Generate registration code
func genRegisterCode(name string, jsCallName string, workName string) (string, string) {
	funName := "register_" + strings.ToLower(name)
	code := `
// ---------- GenRegisterAsyncCode ---------- 
void ` + funName + `(Env env, Object exports){
  napi_property_descriptor desc = {"` + jsCallName + `",NULL,` + workName + `,NULL,NULL,NULL,napi_default,NULL};
  napi_define_properties(env, exports, 1, &desc);
}`
	return code, funName + "(env, exports);"
}

// Parse result arguments
func genResultParseCode(returnType string, resultStructName string) (string, string) {
	code := ""
	preCode := ""
	code += reasync.GenAsyncReturnCommonCode(resultStructName)
	if returnType == "string" {
		ccode, pcode := reasync.GenAsyncReturnStringTypeCode(resultStructName)
		code += ccode
		preCode += pcode
	} else if returnType == "boolean" {
		code += reasync.GenAsyncReturnBooleanTypeCode(resultStructName)
	} else if returnType == "int" {
		code += reasync.GenAsyncReturnIntTypeCode("int32", resultStructName)
	} else if returnType == "int32" {
		code += reasync.GenAsyncReturnIntTypeCode("int32", resultStructName)
	} else if returnType == "int64" {
		code += reasync.GenAsyncReturnIntTypeCode("int64", resultStructName)
	} else if returnType == "uint32" {
		code += reasync.GenAsyncReturnIntTypeCode("unit32", resultStructName)
	} else if returnType == "float" {
		code += reasync.GenAsyncReturnFloatTypeCode(resultStructName)
	} else if returnType == "double" {
		code += reasync.GenAsyncReturnDoubleTypeCode(resultStructName)
	} else if returnType == "array" {
		code += reasync.GenAsyncReturnArrayTypeCode(resultStructName)
	} else if returnType == "object" {
		code += reasync.GenAsyncReturnObjectTypeCode(resultStructName)
	} else if returnType == "arraybuffer" {
		code += reasync.GenAsyncReturnArrayBufferTypeCode(resultStructName)
	} else {
		code += tools.FormatCodeIndentLn(`if (wg_async_res != NULL && wg_async_res->err != NULL) {
    free((void*)wg_async_res->err);
  }
  if (wg_async_res != NULL) {
    free(wg_async_res);
  }`, 2)
	}
	/*
		 else if returnType == "arraybuffer" {
				code += reasync.GenAsyncReturnArrayBufferTypeCode()
			}
	*/

	return code, preCode
}

func genStructCallbackCode(export config.Export, structName string, resultStructName string) string {
	argLen := len(export.Args)
	code := `
typedef struct{
  napi_async_work work;
  napi_threadsafe_function tsfn;
  int argc;
  WgAddonArgInfo *args[` + fmt.Sprintf("%d", argLen) + `];
} ` + structName + `;

typedef struct{
  bool is_error;
  void* data;
  char* err;
} ` + resultStructName + `;`
	return code
}

// Generate async callback code
func GenAsyncCallbackCode(export config.Export) (string, string) {
	methodName := export.Name
	//args := export.Args
	workName := "wg_work_" + strings.ToLower(methodName)
	workCompleteName := "wg_work_complete_" + strings.ToLower(methodName)
	executeworkName := "wg_execute_work" + strings.ToLower(methodName)
	jsCallbackName := "wg_js_callback_" + strings.ToLower(methodName)
	structDataName := "WgAddonData" + methodName
	resultStructName := "WgAsyncResult" + methodName

	code := `
// [` + methodName + `] +++++++++++++++++++++++++++++++++ start`
	code += genStructCallbackCode(export, structDataName, resultStructName)
	code += genJsCallbackCode(export, jsCallbackName, resultStructName)
	code += genExecuteWorkCode(export, executeworkName, structDataName, resultStructName)
	code += genWorkCompleteCode(workCompleteName, structDataName)
	code += genWorkThreadCode(
		export,
		workName,
		workCompleteName,
		executeworkName,
		jsCallbackName,
		structDataName,
	)

	// Generate registration code
	rCode, rFunc := genRegisterCode(methodName, export.JsCallName, workName)
	code += rCode

	code += `
// [` + methodName + `]+++++++++++++++++++++++++++++++++ end`

	registerCode := tools.FormatCodeIndentLn(rFunc, 2)
	return code, registerCode
}
