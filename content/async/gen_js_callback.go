package async

import (
	"github.com/VastBlast/gonode/config"
	"github.com/VastBlast/gonode/content/returns/reasync"
	"github.com/VastBlast/gonode/tools"
)

// JS callback arguments
func genJsCallbackArgs(export config.Export) string {
	returnType := export.ReturnType

	code := ""
	if returnType == "string" {
		code += reasync.GenAsyncCallbackArgStringTypeCode()
	} else if returnType == "boolean" {
		code += reasync.GenAsyncCallbackArgBooleanTypeCode()
	} else if returnType == "int" {
		code += reasync.GenAsyncCallbackArgIntTypeCode()
	} else if returnType == "int32" {
		code += reasync.GenAsyncCallbackArgIntTypeCode()
	} else if returnType == "int64" {
		code += reasync.GenAsyncCallbackArgIntTypeCode()
	} else if returnType == "uint32" {
		code += reasync.GenAsyncCallbackArgIntTypeCode()
	} else if returnType == "float" {
		code += reasync.GenAsyncCallbackArgFloatTypeCode()
	} else if returnType == "double" {
		code += reasync.GenAsyncCallbackArgDoubleTypeCode()
	} else if returnType == "array" {
		code += reasync.GenAsyncCallbackArgArrayTypeCode()
	} else if returnType == "object" {
		code += reasync.GenAsyncCallbackArgObjectTypeCode()
	} else if returnType == "arraybuffer" {
		code += reasync.GenAsyncCallbackArgArrayBufferTypeCode()
	} else {
		code += tools.FormatCodeIndentLn(`napi_value wg_result = wg_env.Null();`, 4)
	}

	return code
}

// JS callback
func genJsCallbackCode(export config.Export, jsCallbackName string, resultStructName string) string {
	parseCode, parsePreCode := genResultParseCode(export.ReturnType, resultStructName)
	code := `
// ------------ genJsCallbackCode
static void ` + jsCallbackName + `(napi_env wg_env, napi_value wg_js_cb, void* wg_context, void* wg_data) {
  (void)wg_context;
#ifdef NAPI_CPP_EXCEPTIONS
  try {
#endif` + parseCode + `
  if (wg_env != NULL) {` + parsePreCode + `
    Env wg_env_scope = Env(wg_env);
    napi_value wg_err = wg_env_scope.Null();
    if (wg_is_error) {
      wg_err = Error::New(wg_env_scope, !wg_err_str.empty() ? wg_err_str.c_str() : "async error").Value();
    }
    ` + genJsCallbackArgs(export) + `
    napi_value wg_argv[] = { wg_err, wg_result };
    napi_value wg_global;
    napi_status wg_sts = napi_get_global(wg_env, &wg_global);
    if (wg_sts != napi_ok) {
      wg_catch_err(wg_env, wg_sts);
      return;
    }
    wg_sts = napi_call_function(wg_env, wg_global, wg_js_cb, 2, wg_argv, NULL);
    if (wg_sts != napi_ok) {
      wg_catch_err(wg_env, wg_sts);
      return;
    }
  }
#ifdef NAPI_CPP_EXCEPTIONS
  } catch (const Error& wg_ex) {
    wg_ex.ThrowAsJavaScriptException();
  } catch (const std::exception& wg_ex) {
    napi_throw_error(wg_env, NULL, wg_ex.what());
  } catch (...) {
    napi_throw_error(wg_env, NULL, "native exception");
  }
#endif
}`
	return code
}
