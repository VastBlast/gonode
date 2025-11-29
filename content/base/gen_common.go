package base

import (
	"github.com/VastBlast/gonode/config"
	"github.com/VastBlast/gonode/tools"
)

func genWgAddonDataCode() string {
	return `
//---------- genWgAddonArg ----------
typedef struct {
  int type; // [1]char [2]int [3]float [4]double [5]bool
  int len;
  void* value;
} WgAddonArgInfo;`
}

func genBuildGoStringCode() string {
	return `
//---------- genBuildGoString ----------
GoString wg_build_go_string(const char* p, size_t n){
  return {p, static_cast<ptrdiff_t>(n)};
}`
}

func genBuildGoSliceCode() string {
	return `
//---------- genBuildGoString ----------
GoSlice wg_build_go_slice(void *data, int len, int cap){
  return { data, len, cap };
}`
}

func genCatchErrCode() string {
	return `
// ------------- genCatchErr -----------
static void wg_catch_err(napi_env env, napi_status status) {
  if (status != napi_ok) {
    const napi_extended_error_info* error_info = NULL;
    napi_get_last_error_info(env, &error_info);
    printf("addon >>>>> %s\n", error_info->error_message);
    exit(0);
  }
}`
}

func genArrayToStringCode() string {
	return `
// ------------- genStringToArray2 -----------
string wg_array_to_string(Array arr, Env env) {
  Value wg_json_val = env.Global().Get("JSON");
  if (!wg_json_val.IsObject()) return "[]";
  Object wg_json = wg_json_val.As<Object>();
  Value wg_stringify_val = wg_json.Get("stringify");
  if (!wg_stringify_val.IsFunction()) return "[]";
  Function wg_stringify = wg_stringify_val.As<Function>();
  Value wg_res = wg_stringify.Call(wg_json, { arr });
  if (wg_res.IsString()) {
    return wg_res.As<String>().Utf8Value();
  }
  return "[]";
}`
}

func genStringToArrayCode() string {
	return `
// ------------- genStringToArray -----------
Array wg_string_to_array(string str, Env env) {
  if (str == "") return Array::New(env);
  Value wg_json_val = env.Global().Get("JSON");
  if (!wg_json_val.IsObject()) return Array::New(env);
  Object wg_json = wg_json_val.As<Object>();
  Value wg_parse_val = wg_json.Get("parse");
  if (!wg_parse_val.IsFunction()) return Array::New(env);
  Function wg_parse = wg_parse_val.As<Function>();
  Value wg_res = wg_parse.Call(wg_json, { String::New(env, str) });
  if (wg_res.IsArray()) {
    return wg_res.As<Array>();
  }
  return Array::New(env);
}`
}

func genObjectToStringCode() string {
	return `
// ------------- genObjectToString -----------
string wg_object_to_string(Object objs, Env env) {
  Value wg_json_val = env.Global().Get("JSON");
  if (!wg_json_val.IsObject()) return "{}";
  Object wg_json = wg_json_val.As<Object>();
  Value wg_stringify_val = wg_json.Get("stringify");
  if (!wg_stringify_val.IsFunction()) return "{}";
  Function wg_stringify = wg_stringify_val.As<Function>();
  Value wg_res = wg_stringify.Call(wg_json, { objs });
  if (wg_res.IsString()) {
    return wg_res.As<String>().Utf8Value();
  }
  return "{}";
}`
}

func genStringToObject() string {
	code := `
// ------------- genStringToObject -----------
 Object wg_string_to_object(string str, Env env) {
  if (str == "") return Object::New(env);
  Value wg_json_val = env.Global().Get("JSON");
  if (!wg_json_val.IsObject()) return Object::New(env);
  Object wg_json = wg_json_val.As<Object>();
  Value wg_parse_val = wg_json.Get("parse");
  if (!wg_parse_val.IsFunction()) return Object::New(env);
  Function wg_parse = wg_parse_val.As<Function>();
  Value wg_res = wg_parse.Call(wg_json, { String::New(env, str) });
  if (wg_res.IsObject()) {
    return wg_res.As<Object>();
  }
  return Object::New(env);
}`

	return code
}

func GenBeforeCode(hasAsync bool) string {
	code := `// [common]++++++++++++++++++++++++++++++++++++++ start`
	code += genWgAddonDataCode()
	code += genBuildGoStringCode()
	code += genBuildGoSliceCode()
	code += genArrayToStringCode()
	code += genStringToArrayCode()

	code += genObjectToStringCode()
	code += genStringToObject()

	if hasAsync {
		code += genCatchErrCode()
	}

	code += `
// [common]++++++++++++++++++++++++++++++++++++++ end`
	return code
}

func genExportsJsCallApi(exports []config.Export) string {
	code := ""
	for _, export := range exports {
		jsApiName := export.JsCallName
		goApiName := export.Name
		if export.JsCallMode == "sync" {
			code += tools.FormatCodeIndentLn(`exports.Set(String::New(env, "`+jsApiName+`"), Function::New(env, _`+goApiName+`));`, 2)
		}
	}

	return code
}

func GenAfterCode(cfg config.Config, asyncCode string) string {
	name := cfg.Name

	exportsCode := genExportsJsCallApi(cfg.Exports)
	code := `
Object Init(Env env, Object exports) {` + exportsCode + asyncCode + `
  return exports;
}

NODE_API_MODULE(` + name + `, Init)`
	return code
}
