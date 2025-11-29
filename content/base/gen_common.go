package base

import (
	"github.com/wenlng/gonacli/config"
	"github.com/wenlng/gonacli/tools"
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

func genStringSplitCode() string {
	return `
// ------------- genStringSplit -----------
void wg_string_split(const string& str, const char split, vector<string>& res){
  if (str.empty()) return;
  size_t start = 0;
  size_t pos = str.find(split, start);
  while (pos != string::npos){
	res.emplace_back(str.substr(start, pos - start));
	start = pos + 1;
	pos = str.find(split, start);
  }
  res.emplace_back(str.substr(start));
}`
}

func genArrayToStringCode() string {
	return `
// ------------- genStringToArray2 -----------
string wg_array_to_string(Array arr) {
  string res = "[";
  for(uint32_t i = 0; i < arr.Length(); i++){
    if (i > 0) {
      res += ",";
    }
    Value v = arr[i];
    if (v.IsArray()){
      Array arr2 = v.As<Array>();
      res += wg_array_to_string(arr2); 
    } else { 
      string ss = v.ToString();
      res += "\"" + ss + "\""; 
    }
  }
  res += "]";
  return res;
}`
}

func genStringToArrayCode() string {
	return `
// ------------- genStringToArray -----------
Array wg_string_to_array(string str, Env env) {
  Array arr = Array::New(env);
  vector<string> strList;
  if (str == "") return arr;
  for (char &c : str) {
    if (c == '[' || c == ']') {
      c = ',';
    }
  }  
  wg_string_split(str, ',', strList);
  int index = 0;
  for (auto s : strList) {
    if (s.size() > 0) {
      int _spos = s.find("\"");
      s = s.substr(_spos + 1);
      int _epos = s.find("\"");
      s = s.substr(0, _epos);
      arr.Set(Number::New(env, index), String::New(env, s));
      index++;
    }
  }
  return arr;
}`
}

func genObjectArrToStringCode() string {
	return `
// ------------- genObjectArrToString -----------
string wg_object_to_string(Object objs);
string wg_object_array_to_string(Array arr) {
  string res = "[";
  for(uint32_t i = 0; i < arr.Length(); i++){
    if (i > 0) {
      res += ",";
    }
    Value v = arr[i];
    if (v.IsArray()){
      Array arr2 = v.As<Array>();
      res += wg_object_array_to_string(arr2);
    } else if (v.IsObject()){
      Object obj2 = v.As<Object>();
      res += wg_object_to_string(obj2); 
    } else {
      string ss = v.ToString();
      res += "\"" + ss + "\""; 
    }
  }
  res += "]";
  return res;
}`
}

func genObjectToStringCode() string {
	return `
// ------------- genObjectToString -----------
string wg_object_to_string(Object objs) {
  string res = "{";
  Array keyArr = objs.GetPropertyNames();
  for(uint32_t i = 0; i < keyArr.Length(); i++){
    if (i > 0) {
      res += ",";
    }
    Value key = keyArr[i];
    Value v = objs.Get(key);
    string name = key.As<String>().Utf8Value();
    res += "\"" + name + "\":";
    if (v.IsArray()) {
      Array arr = v.As<Array>();
      res += wg_object_array_to_string(arr);
    } else if (v.IsObject()){
      Object obj2 = v.As<Object>();
      res += wg_object_to_string(obj2); 
    } else {
      string ss = v.ToString();
      res += "\"" + ss + "\""; 
    }
  }
  res += "}";
  return res;
}`
}

func genStringToObject() string {
	code := `
// ------------- genStringToObject -----------
Object wg_string_to_object(string str, Env env) {
  Object obj = Object::New(env);
  vector<string> strList;
  if (str == "") return obj;
  for (char &c : str) {
    if (c == '{' || c == '}') {
      c = ',';
    }
  }  
  wg_string_split(str, ',', strList);
  for (auto s : strList) {
    size_t pos;
    if (s.size() > 0) {
      while ((pos = s.find("\"")) != string::npos) {
        s.replace(pos, 1, "");
      }
      vector<string> keyValue;
      wg_string_split(s, ':', keyValue);
      string key = keyValue.size() >= 0 ? keyValue[0] : "";
      string value = keyValue.size() >= 1 ? keyValue[1] : "";
      obj.Set(String::New(env, key), String::New(env, value));
    }
  }
  return obj;
}`

	return code
}

func GenBeforeCode(hasAsync bool) string {
	code := `// [common]++++++++++++++++++++++++++++++++++++++ start`
	code += genWgAddonDataCode()
	//code += genBuildGoStringCode()
	//code += genBuildGoSliceCode()
	code += genStringSplitCode()
	code += genArrayToStringCode()
	code += genStringToArrayCode()

	code += genObjectArrToStringCode()
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
